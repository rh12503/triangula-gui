package main

import (
	"github.com/RH12503/Triangula/color"
	"github.com/RH12503/Triangula/normgeom"
)

// RenderData is sent to the frontend for rendering
type RenderData struct {
	Width, Height int
	Polygons      []normgeom.NormPolygon
	Colors        []color.RGB
}

