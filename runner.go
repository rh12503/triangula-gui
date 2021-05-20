package main

import (
	"encoding/base64"
	"errors"
	"github.com/RH12503/Triangula-GUI/util"
	"github.com/RH12503/Triangula/algorithm"
	image2 "github.com/RH12503/Triangula/image"
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
	triangulated = iota
	polygonal    = iota
)

type Runner struct {
	algorithm algorithm.Algorithm // The algorithm being used
	runtime   *wails.Runtime

	running      bool // Indicates if the algorithm is running or not
	runningMutex sync.Mutex

	stopped      bool // Indicates if the algorithm is stopped or not. !stopped doesn't mean running as the algorithm could be paused
	stoppedMutex sync.Mutex

	tempPauseMutex sync.Mutex // Used to indicate for the algorithm to temporarily pause.

	frameTime int // The increment between each frame rendered

	image     image.Image // The image selected
	pixelData image2.Data // The image data used by the algorithm

	lightMode bool // If light mode is set

	logic Logic
}

func (r *Runner) WailsInit(runtime *wails.Runtime) error {
	r.runtime = runtime
	r.lightMode = true
	return nil
}

// Run runs the algorithm if it is not already running
func (r *Runner) Run(shape int, mutations int, mutationAmount float64, numPoints, population, cutoff, blockSize, cacheSize, threads, frameTime int) {

	switch shape {
	case triangulated:
		r.logic = TriangleLogic{}
		break
	case polygonal:
		r.logic = PolygonLogic{}
		break
	}

	if threads != 0 {
		runtime.GOMAXPROCS(threads)
	}

	if r.image == nil {
		return
	}

	if !r.Running() {
		// Setup algorithm
		img := image2.ToData(r.image)
		r.pixelData = img

		r.algorithm = r.logic.NewAlgorithm(img, mutations, mutationAmount, numPoints, population, cutoff, blockSize, cacheSize)

		r.stopped = false
		r.frameTime = frameTime

		r.StartAlgorithm()

		r.runtime.Events.Emit("running") // Notify frontend
	}
}

// Stop stops the algorithm
func (r *Runner) Stop() {
	r.runningMutex.Lock()
	r.stopped = true
	if !r.running {
		r.runtime.Events.Emit("resumed")
		r.runtime.Events.Emit("stopped")
	} else {
		r.running = false
	}
	r.runningMutex.Unlock()
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
	format := util.MimeType(bytes)

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
	r.tempPauseMutex.Lock()
	best := r.algorithm.Best().Copy()
	r.tempPauseMutex.Unlock()
	return r.logic.SaveSVG(name, best, r.pixelData)
}

// SavePNG opens a dialog to save the result of the algorithm to a PNG file
func (r *Runner) SavePNG(scale float64, effect int) error {
	if r.algorithm == nil {
		return errors.New("algorithm not initialized")
	}

	name := r.runtime.Dialog.SelectSaveFile("Export to PNG", "*.png")

	r.tempPauseMutex.Lock()
	best := r.algorithm.Best().Copy()
	r.tempPauseMutex.Unlock()

	return r.logic.SavePNG(name, best, r.pixelData, scale, effect)
}

// StartAlgorithm starts the algorithm if it isn't already started
func (r *Runner) StartAlgorithm() {
	r.runningMutex.Lock()
	r.running = true
	r.runningMutex.Unlock()
	go func() {
	out:
		for {
			r.tempPauseMutex.Lock()
			best := r.algorithm.Best().Copy()
			r.tempPauseMutex.Unlock()
			// Send rendering data to the frontend
			r.runtime.Events.Emit("renderData", r.logic.RenderData(best, r.pixelData))
			r.runtime.Events.Emit("stats", r.algorithm.Stats())

			ti := time.Now()
			statsTime := time.Now()

			for time.Since(ti).Milliseconds() < int64(r.frameTime) {
				r.tempPauseMutex.Lock()
				r.algorithm.Step()
				r.tempPauseMutex.Unlock()

				if !r.Running() {
					break out
				}
				if time.Since(statsTime).Milliseconds() < 200 { // Update stats at least 5 times per second
					r.runtime.Events.Emit("stats", r.algorithm.Stats())
					statsTime = time.Now()
				}
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
		r.StartAlgorithm()
		r.runtime.Events.Emit("resumed")
	}
}

// ToggleMode toggles light/dark mode
func (r *Runner) ToggleMode() {
	r.lightMode = !r.lightMode

	if r.lightMode {
		r.runtime.Events.Emit("lightmode")
	} else {
		r.runtime.Events.Emit("darkmode")
	}
}
