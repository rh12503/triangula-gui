package main

import (
	"encoding/base64"
	"errors"
	"github.com/RH12503/Triangula-GUI/export"
	"github.com/RH12503/Triangula/algorithm"
	"github.com/RH12503/Triangula/algorithm/evaluator"
	"github.com/RH12503/Triangula/generator"
	image2 "github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/mutation"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/render"
	"github.com/RH12503/Triangula/triangulation"
	"github.com/wailsapp/wails"
	"image"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	none     int = iota
	gradient int = iota
	split    int = iota
)

type Runner struct {
	algorithm    algorithm.Algorithm // The algorithm being used
	runtime      *wails.Runtime

	running      bool // Indicates if the algorithm is running or not
	runningMutex sync.Mutex

	stopped      bool // Indicates if the algorithm is stopped or not. stopped doesn't mean !running as the algorithm could be paused
	stoppedMutex sync.Mutex

	frameTime int // The increment between each frame rendered

	image     image.Image // The image selected
	normImage image2.Data // The image data used by the algorithm
}

func (r *Runner) WailsInit(runtime *wails.Runtime) error {
	r.runtime = runtime

	return nil
}

// Run runs the algorithm if it is not already running
func (r *Runner) Run(mutations int, mutationAmount float64, numPoints, population, cutoff, blockSize, cacheSize, threads, frameTime int) {

	if threads != 0 {
		runtime.GOMAXPROCS(threads)
	}

	if r.image == nil {
		return
	}

	if !r.Running() {
		// Setup algorithm
		img := image2.ToData(r.image)
		r.normImage = img

		evaluatorFactory := func(n int) evaluator.Evaluator {
			return evaluator.NewParallel(img, cacheSize, blockSize, n)
		}
		var mutator mutation.Method
		mutator = mutation.NewGaussianMethod(float64(mutations)/float64(numPoints), mutationAmount)

		pointFactory := func() normgeom.NormPointGroup {
			return generator.RandomGenerator{}.Generate(numPoints)
		}

		algo := algorithm.NewModifiedGenetic(pointFactory, population, cutoff, evaluatorFactory, mutator)
		r.algorithm = algo

		r.stopped = false
		r.frameTime = frameTime


		r.Start()

		r.runtime.Events.Emit("running") // Notify frontend
	}
}

// Stop stops the algorithm
func (r *Runner) Stop() {
	r.runningMutex.Lock()
	r.running = false
	r.runningMutex.Unlock()
	r.stopped = true
}

// Running returns if the algorithm is running or not
func (r *Runner) Running() bool {
	r.runningMutex.Lock()
	defer r.runningMutex.Unlock()
	return r.running
}

// SelectImage opens a dialog to select an image file
func (r *Runner) SelectImage() {
	path := r.runtime.Dialog.SelectFile("Select an image", "*.jpg,*.png,*.jpeg")

	if path == "" {
		return
	}
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	imageFile, _, err := image.Decode(file)

	// Obtain the mime type of the image
	bytes, _ := ioutil.ReadFile(path)
	format := mimeType(bytes)

	file.Close()

	if err != nil {
		log.Fatal(err)
	}

	r.image = imageFile

	imageBase64 := "data:" + format + ";base64," + base64.StdEncoding.EncodeToString(bytes)

	r.runtime.Events.Emit("image", path, imageBase64)
}

// LoadImage processes an image to be used
func (r *Runner) LoadImage(name, data, base string) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))

	imageFile, _, err := image.Decode(reader)

	if err != nil {
		log.Fatal(err)
	}

	r.image = imageFile
	r.runtime.Events.Emit("image", name, base+data)
}

// SaveSVG opens a dialog to save the result of the algorithm to a SVG file
func (r *Runner) SaveSVG() error {
	if r.algorithm == nil {
		return errors.New("algorithm not initialized")
	}

	name := r.runtime.Dialog.SelectSaveFile("Export to SVG", "*.svg")
	err := export.WriteSVG(name, r.algorithm.Best(), r.normImage)

	return err
}

// SavePNG opens a dialog to save the result of the algorithm to a PNG file
func (r *Runner) SavePNG(scale float64, effect int) error {
	if r.algorithm == nil {
		return errors.New("algorithm not initialized")
	}

	name := r.runtime.Dialog.SelectSaveFile("Export to PNG", "*.png")

	var err error

	if effect == none {
		err = export.WritePNG(name, r.algorithm.Best(), r.normImage, scale)
	} else if effect == gradient {
		err = export.WriteEffectPNG(name, r.algorithm.Best(), r.normImage, scale, true)
	} else if effect == split {
		err = export.WriteEffectPNG(name, r.algorithm.Best(), r.normImage, scale, false)
	}

	return err
}

// Start starts the algorithm if it isn't already started
func (r *Runner) Start() {
	r.runningMutex.Lock()
	r.running = true
	r.runningMutex.Unlock()
	go func() {
		for r.Running() {
			w, h := r.normImage.Size()
			triangles, _ := triangulation.Triangulate(r.algorithm.Best(), w, h)
			triangleData := render.TrianglesOnImage(triangles, r.normImage)

			// Send rendering data to the frontend
			r.runtime.Events.Emit("renderData", RenderData{
				Width:  w,
				Height: h,
				Data:   triangleData,
			})
			r.runtime.Events.Emit("stats", r.algorithm.Stats())
			ti := time.Now()
			for time.Since(ti).Milliseconds() < int64(r.frameTime) {
				r.algorithm.Step()
			}
		}
		if r.stopped {
			r.runtime.Events.Emit("stopped")
		}
	}()
}

// TogglePause pauses or resumes the algorithm
func (r *Runner) TogglePause() {
	if r.algorithm == nil || r.stopped {
		return
	}

	if r.Running() {
		r.runningMutex.Lock()
		r.running = false
		r.runningMutex.Unlock()
		r.runtime.Events.Emit("paused")
	} else {
		r.Start()
		r.runtime.Events.Emit("resumed")
	}
}

// RenderData is sent to the frontend for rendering
type RenderData struct {
	Width, Height int
	Data          []render.TriangleData
}

//https://stackoverflow.com/a/25959527/15283541
var types = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
}

func mimeType(incipit []byte) string {
	incipitStr := string(incipit)
	for magic, mime := range types {
		if strings.HasPrefix(incipitStr, magic) {
			return mime
		}
	}

	return ""
}
