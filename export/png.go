// Package export provides utility functions for saving a result from the algorithm to a file,
// as well as different filters that can be applied to the results.

package export

import (
	"github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/render"
	"github.com/RH12503/Triangula/triangulation"
	"github.com/fogleman/gg"
	"math"
)

// WritePNG saves a PNG of a result
func WritePNG(filename string, points normgeom.NormPointGroup, img image.Data, pixelScale float64) error {
	imageW, imageH := img.Size()

	w := int(math.Round(float64(imageW) * pixelScale))
	h := int(math.Round(float64(imageH) * pixelScale))

	dc := gg.NewContext(w, h)
	triangles, _ := triangulation.Triangulate(points, imageW, imageH)
	triangleData := render.TrianglesOnImage(triangles, img)


	dc.Fill()
	for _, d := range triangleData {
		tri := d.Triangle.Points
		col := d.Color

		dc.SetRGB(col.R, col.G, col.B)

		dc.NewSubPath()
		dc.LineTo(tri[0].X*float64(w), tri[0].Y*float64(h))
		dc.LineTo(tri[1].X*float64(w), tri[1].Y*float64(h))
		dc.LineTo(tri[2].X*float64(w), tri[2].Y*float64(h))
		dc.ClosePath()

		dc.SetLineWidth(1)
		dc.Stroke()

		dc.NewSubPath()
		dc.LineTo(tri[0].X*float64(w), tri[0].Y*float64(h))
		dc.LineTo(tri[1].X*float64(w), tri[1].Y*float64(h))
		dc.LineTo(tri[2].X*float64(w), tri[2].Y*float64(h))
		dc.ClosePath()
		dc.Fill()
	}
	err := dc.SavePNG(filename)


	return err
}
