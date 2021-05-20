package polygons

import (
	"github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/polygonation"
	"github.com/RH12503/Triangula/render"
	"github.com/fogleman/gg"
	"math"
)

// WritePNG saves a PNG of a result
func WritePNG(filename string, points normgeom.NormPointGroup, img image.Data, pixelScale float64) error {
	imageW, imageH := img.Size()

	w := int(math.Round(float64(imageW) * pixelScale))
	h := int(math.Round(float64(imageH) * pixelScale))

	dc := gg.NewContext(w, h)
	polygons := polygonation.Polygonate(points, imageW, imageH)
	polygonsData := render.PolygonsOnImage(polygons, img)

	dc.Fill()

	for _, d := range polygonsData {
		poly := d.Polygon.Points
		col := d.Color

		dc.SetRGB(col.R, col.G, col.B)

		dc.NewSubPath()
		for _, p := range poly {
			dc.LineTo(p.X*float64(w), p.Y*float64(h))
		}
		dc.ClosePath()

		dc.SetLineWidth(1)
		dc.Stroke()

		dc.NewSubPath()
		for _, p := range poly {
			dc.LineTo(p.X*float64(w), p.Y*float64(h))
		}
		dc.ClosePath()
		dc.Fill()
	}
	err := dc.SavePNG(filename)


	return err
}
