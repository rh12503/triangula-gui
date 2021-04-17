package export

import (
	color2 "github.com/RH12503/Triangula/color"
	"github.com/RH12503/Triangula/geom"
	"github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/rasterize"
	"github.com/RH12503/Triangula/render"
	"github.com/RH12503/Triangula/triangulation"
	"github.com/fogleman/gg"
	"image/color"
	"math"
)

// WriteEffectPNG saves a PNG of a result with an effect applied
func WriteEffectPNG(filename string, points normgeom.NormPointGroup, img image.Data, pixelScale float64, gradient bool) error {
	imageW, imageH := img.Size()

	w := multAndRound(imageW, pixelScale)
	h := multAndRound(imageH, pixelScale)

	dc := gg.NewContext(w, h)

	dc.SetColor(color.White)
	dc.DrawRectangle(0, 0, float64(w), float64(h))
	dc.Fill()

	triangles := triangulation.Triangulate(points, imageW, imageH)
	triangleData := render.TrianglesOnImage(triangles, img)

	for i, _ := range triangleData {
		tri := triangles[i]
		points := tri.Points

		y2y3 := points[1].Y - points[2].Y
		x3x2 := points[2].X - points[1].X
		x1x3 := points[0].X - points[2].X
		y1y3 := points[0].Y - points[2].Y
		y3y1 := points[2].Y - points[0].Y

		dcol := float64(y2y3*x1x3 + x3x2*y1y3)
		avg0 := color2.AverageRGB{}
		avg1 := color2.AverageRGB{}
		avg2 := color2.AverageRGB{}

		rasterize.DDATriangle(tri, func(x, y int) {
			xx3 := x - points[2].X
			yy3 := y - points[2].Y

			// Calculate Barymetric coordinates for filling gradients
			l0 := math.Max(float64(y2y3*xx3+x3x2*yy3)/dcol, 0)
			l1 := math.Max(float64(y3y1*xx3+x1x3*yy3)/dcol, 0)

			l2 := math.Max(1-l0-l1, 0)

			max := math.Max(l0, math.Max(l1, l2))

			col := img.RGBAt(x, y)

			if max == l0 {
				avg0.Add(col)
			} else if max == l1 {
				avg1.Add(col)
			} else {
				avg2.Add(col)
			}
		})

		c0 := avg0.Average()
		c1 := avg1.Average()
		c2 := avg2.Average()

		// Prevent blank triangles
		if avg0.Count() == 0 {
			c0 = img.RGBAt(min(points[0].X, imageW-1), min(points[0].Y, imageH-1))
		}

		if avg1.Count() == 0 {
			c1 = img.RGBAt(min(points[1].X, imageW-1), min(points[1].Y, imageH-1))
		}

		if avg2.Count() == 0 {
			c2 = img.RGBAt(min(points[2].X, imageW-1), min(points[2].Y, imageH-1))
		}

		scaledTri := geom.NewTriangle(
			multAndRound(points[0].X, pixelScale),
			multAndRound(points[0].Y, pixelScale),
			multAndRound(points[1].X, pixelScale),
			multAndRound(points[1].Y, pixelScale),
			multAndRound(points[2].X, pixelScale),
			multAndRound(points[2].Y, pixelScale),
		)

		rasterize.DDATriangle(scaledTri, func(x, y int) {
			xx3 := float64(x)/pixelScale - float64(points[2].X)
			yy3 := float64(y)/pixelScale - float64(points[2].Y)

			l0 := math.Max((float64(y2y3)*xx3+float64(x3x2)*yy3)/dcol, 0)
			l1 := math.Max((float64(y3y1)*xx3+float64(x1x3)*yy3)/dcol, 0)

			l2 := math.Max(1-l0-l1, 0)

			if gradient {
				dc.SetColor(color.RGBA{
					R: uint8(scale(math.Min(c0.R*l0+c1.R*l1+c2.R*l2, 1), 255)),
					G: uint8(scale(math.Min(c0.G*l0+c1.G*l1+c2.G*l2, 1), 255)),
					B: uint8(scale(math.Min(c0.B*l0+c1.B*l1+c2.B*l2, 1), 255)),
					A: 255,
				})
			} else {
				max := math.Max(l0, math.Max(l1, l2))

				if max == l0 {
					dc.SetColor(color.RGBA{
						R: uint8(scale(c0.R, 255)),
						G: uint8(scale(c0.G, 255)),
						B: uint8(scale(c0.B, 255)),
						A: 255,
					})
				} else if max == l1 {
					dc.SetColor(color.RGBA{
						R: uint8(scale(c1.R, 255)),
						G: uint8(scale(c1.G, 255)),
						B: uint8(scale(c1.B, 255)),
						A: 255,
					})
				} else {
					dc.SetColor(color.RGBA{
						R: uint8(scale(c2.R, 255)),
						G: uint8(scale(c2.G, 255)),
						B: uint8(scale(c2.B, 255)),
						A: 255,
					})
				}
			}

			dc.SetPixel(x, y)
		})
	}

	err := dc.SavePNG(filename)

	return err
}