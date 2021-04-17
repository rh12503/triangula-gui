package export

import (
	"bufio"
	"fmt"
	"github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/render"
	"github.com/RH12503/Triangula/triangulation"
	"os"
)

const svgStart = `<?xml version="1.0"?>
<svg width="%v" height="%v"
     xmlns="http://www.w3.org/2000/svg"
     shape-rendering="crispEdges">
`

const svgPoly = `<polygon points="%v,%v %v,%v %v,%v" style="fill:rgb(%v,%v,%v)"/>
`

// WriteSVG saves a SVG of a result
func WriteSVG(filename string, points normgeom.NormPointGroup, img image.Data) error {

	outFile, err := os.Create(filename)

	if err != nil {
		return err
	}

	w, h := img.Size()
	triangles := triangulation.Triangulate(points, w, h)
	triangleData := render.TrianglesOnImage(triangles, img)

	writer := bufio.NewWriter(outFile)
	writer.WriteString(fmt.Sprintf(svgStart, w, h))
	for _, d := range triangleData {
		tri := d.Triangle.Points
		col := d.Color

		writer.WriteString(
			fmt.Sprintf(
				svgPoly, scale(tri[0].X, w), scale(tri[0].Y, h),
				scale(tri[1].X, w), scale(tri[1].Y, h),
				scale(tri[2].X, w), scale(tri[2].Y, h),
				scale(col.R, 255), scale(col.G, 255), scale(col.B, 255),
			),
		)
	}
	writer.WriteString("</svg>")
	writer.Flush()
	outFile.Close()

	return nil
}

