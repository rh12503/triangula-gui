package polygons

import (
	"bufio"
	"fmt"
	"github.com/RH12503/Triangula-GUI/util"
	"github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/polygonation"
	"github.com/RH12503/Triangula/render"
	"os"
)

// WriteSVG saves a SVG of a result
func WriteSVG(filename string, points normgeom.NormPointGroup, img image.Data) error {

	outFile, err := os.Create(filename)

	if err != nil {
		return err
	}

	w, h := img.Size()
	polygons := polygonation.Polygonate(points, w, h)
	polygonsData := render.PolygonsOnImage(polygons, img)

	writer := bufio.NewWriter(outFile)
	writer.WriteString(fmt.Sprintf(util.SvgStart, w, h))

	for _, d := range polygonsData {
		poly := d.Polygon.Points
		col := d.Color

		writer.WriteString(`<polygon points="`)

		for _, p := range poly {
			writer.WriteString(fmt.Sprintf("%v,%v ", util.Scale(p.X, w), util.Scale(p.Y, h)))
		}

		writer.WriteString(fmt.Sprintf(`" style="fill:rgb(%v,%v,%v)"/>\n`,
			util.Scale(col.R, 255), util.Scale(col.G, 255), util.Scale(col.B, 255)))
	}
	writer.WriteString("</svg>")
	writer.Flush()
	outFile.Close()

	return nil
}
