package triangles

import (
	"bufio"
	"fmt"
	"github.com/RH12503/Triangula-GUI/util"
	"github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/render"
	"github.com/RH12503/Triangula/triangulation"
	"os"
)

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
	writer.WriteString(fmt.Sprintf(util.SvgStart, w, h))
	for _, d := range triangleData {
		tri := d.Triangle.Points
		col := d.Color

		writer.WriteString(
			fmt.Sprintf(
				svgPoly, util.Scale(tri[0].X, w), util.Scale(tri[0].Y, h),
				util.Scale(tri[1].X, w), util.Scale(tri[1].Y, h),
				util.Scale(tri[2].X, w), util.Scale(tri[2].Y, h),
				util.Scale(col.R, 255), util.Scale(col.G, 255), util.Scale(col.B, 255),
			),
		)
	}
	writer.WriteString("</svg>")
	writer.Flush()
	outFile.Close()

	return nil
}

