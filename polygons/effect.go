package polygons

import (
	"github.com/RH12503/Triangula-GUI/util"
	color2 "github.com/RH12503/Triangula/color"
	"github.com/RH12503/Triangula/geom"
	"github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/rasterize"
	"github.com/RH12503/Triangula/triangulation/incrdelaunay"
	"github.com/fogleman/gg"
	"image/color"
	"math"
)

// WriteEffectPNG saves a PNG of a result with an effect applied
func WriteEffectPNG(filename string, points normgeom.NormPointGroup, img image.Data, pixelScale float64, gradient bool) error {
	imageW, imageH := img.Size()

	w := util.MultAndRound(imageW, pixelScale)
	h := util.MultAndRound(imageH, pixelScale)

	dc := gg.NewContext(w, h)

	dc.SetColor(color.White)
	dc.DrawRectangle(0, 0, float64(w), float64(h))
	dc.Fill()

	edgeDist := 1/(pixelScale*pixelScale)

	triangulation := incrdelaunay.NewDelaunay(imageW, imageH)
	for _, p := range points {
		triangulation.Insert(incrdelaunay.Point{
			X: int16(math.Round(p.X * float64(imageW))),
			Y: int16(math.Round(p.Y * float64(imageH))),
		})
	}

	incrdelaunay.Voronoi(triangulation, func(pts []incrdelaunay.FloatPoint) {
		var poly geom.Polygon
		var points []incrdelaunay.FloatPoint

		for _, p := range pts {
			new := geom.Point{
				X: int(math.Round(p.X)),
				Y: int(math.Round(p.Y)),
			}

			if len(poly.Points) == 0 || poly.Points[len(poly.Points)-1] != new {
				poly.Points = append(poly.Points, new)
				points = append(points, p)
			}
		}

		n := len(poly.Points)

		vertexAverageColors := make([]color2.AverageRGB, n)

		rasterize.DDAPolygon(poly, func(x, y int) {
			weightSum := 0.
			weights := make([]float64, n)
			point := incrdelaunay.FloatPoint{X: float64(x), Y: float64(y)}

			for i, p := range points {
				nextI := (i + 1) % n

				prev := points[(i+n-1)%n]
				next := points[nextI]

				c := cross(sub(next,p), sub(point,p))
				if c*c <= float64(distSq(next,p)) {
					for j := 0; j < i; j++ {
						weights[j] = 0
					}
					weightSum = 1

					totalDist := math.Sqrt(float64(distSq(p,next)))
					dist := math.Sqrt(float64(distSq(p,point)))

					lerp := dist / totalDist

					lerp = math.Min(1, math.Max(lerp, 0))

					weights[i] = 1 - lerp
					weights[nextI] = lerp

					break
				}

				weights[i] = (cot(point, p, prev) + cot(point, p, next)) / float64(distSq(point,p))
				weightSum += weights[i]
			}



			highest := 0
			highestWeight := -1.
			for i, w := range weights {
				weight := w / weightSum
				if weight > highestWeight {
					highest = i
					highestWeight = weight
				}
			}

			vertexAverageColors[highest].Add(img.RGBAt(x, y))
		})

		vertexColors := make([]color2.RGB, n)

		for i, a := range vertexAverageColors {
			if a.Count() == 0 {
				p := poly.Points[i]
				vertexColors[i] = img.RGBAt(util.Min(p.X, imageW-1), util.Min(p.Y, imageH-1))
			} else {
				vertexColors[i] = a.Average()
			}
		}

		scaledPoly := geom.Polygon{Points: make([]geom.Point, len(poly.Points))}

		for i := range scaledPoly.Points {
			scaledPoly.Points[i].X = util.MultAndRound(poly.Points[i].X, pixelScale)
			scaledPoly.Points[i].Y = util.MultAndRound(poly.Points[i].Y, pixelScale)
		}


		rasterize.DDAPolygon(scaledPoly, func(x, y int) {
			weightSum := 0.
			weights := make([]float64, n)
			point := incrdelaunay.FloatPoint{X: float64(x)/pixelScale, Y: float64(y)/pixelScale}

			for i, p := range points {
				nextI := (i + 1) % n

				prev := points[(i+n-1)%n]
				next := points[nextI]

				c := cross(sub(next,p), sub(point,p))
				if c*c <= edgeDist*float64(distSq(next,p)) {
					for j := 0; j < i; j++ {
						weights[j] = 0
					}
					weightSum = 1

					totalDist := math.Sqrt(float64(distSq(p,next)))
					dist := math.Sqrt(float64(distSq(p,point)))

					lerp := dist / totalDist

					lerp = math.Min(1, math.Max(lerp, 0))

					weights[i] = 1 - lerp
					weights[nextI] = lerp

					break
				}

				weights[i] = (cot(point, p, prev) + cot(point, p, next)) / float64(distSq(point,p))
				weightSum += weights[i]
			}

			r, g, b := 0., 0., 0.

			if gradient {
				for i, w := range weights {
					weight := w / weightSum

					col := vertexColors[i]
					r += col.R * weight
					g += col.G * weight
					b += col.B * weight
				}
			} else {
				highest := 0
				highestWeight := -1.
				for i, w := range weights {
					weight := w / weightSum
					if weight > highestWeight {
						highest = i
						highestWeight = weight
					}
				}
				col := vertexColors[highest]
				r = col.R
				g = col.G
				b = col.B
			}

			r = math.Max(0, math.Min(r, 1))
			g = math.Max(0, math.Min(g, 1))
			b = math.Max(0, math.Min(b, 1))

			dc.SetColor(color.RGBA{
				R: uint8(util.Scale(r, 255)),
				G: uint8(util.Scale(g, 255)),
				B: uint8(util.Scale(b, 255)),
				A: 255,
			})

			dc.SetPixel(x, y)
		})

	}, imageW, imageH)


	err := dc.SavePNG(filename)


	return err
}

func sub(a, b incrdelaunay.FloatPoint) incrdelaunay.FloatPoint {
	return incrdelaunay.FloatPoint{
		X: a.X - b.X,
		Y: a.Y - b.Y,
	}
}

func distSq(a, b incrdelaunay.FloatPoint) float64 {
	dX := a.X - b.X
	dY := a.Y - b.Y
	return dX*dX + dY*dY
}

func cot(a, b, c incrdelaunay.FloatPoint) float64 {
	ba := sub(a,b)
	bc := sub(c,b)

	// Dot product / Cross product
	return (ba.X*bc.X + ba.Y*bc.Y) / cross(ba, bc)
}

func cross(a, b incrdelaunay.FloatPoint) float64 {
	return math.Abs(float64((b.X * a.Y) - (b.Y * a.X)))
}
