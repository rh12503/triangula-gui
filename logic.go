package main

import (
	"errors"
	"github.com/RH12503/Triangula-GUI/polygons"
	"github.com/RH12503/Triangula-GUI/triangles"
	"github.com/RH12503/Triangula/algorithm"
	"github.com/RH12503/Triangula/algorithm/evaluator"
	"github.com/RH12503/Triangula/color"
	"github.com/RH12503/Triangula/fitness"
	"github.com/RH12503/Triangula/generator"
	"github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/mutation"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/polygonation"
	"github.com/RH12503/Triangula/render"
	"github.com/RH12503/Triangula/triangulation"
	"strings"
)

const (
	none     int = iota
	gradient int = iota
	split    int = iota
)

type Logic interface {
	NewAlgorithm(img image.Data, mutations int, mutationAmount float64, numPoints, population, cutoff, blockSize, cacheSize int) algorithm.Algorithm
	RenderData(normgeom.NormPointGroup, image.Data) RenderData
	SaveSVG(file string, points normgeom.NormPointGroup, img image.Data) error
	SavePNG(file string, points normgeom.NormPointGroup, img image.Data, scale float64, effect int) error
}

type TriangleLogic struct {
}

func (t TriangleLogic) SaveSVG(file string, points normgeom.NormPointGroup, img image.Data) error {
	filename := file
	if !strings.HasSuffix(filename, ".svg") {
		filename += ".svg"
	}

	return triangles.WriteSVG(filename, points, img)
}

func (t TriangleLogic) SavePNG(file string, points normgeom.NormPointGroup, img image.Data, scale float64, effect int) error {
	filename := file
	if !strings.HasSuffix(filename, ".png") {
		filename += ".png"
	}

	if effect == none {
		return triangles.WritePNG(filename, points, img, scale)
	} else if effect == gradient {
		return triangles.WriteEffectPNG(filename, points, img, scale, true)
	} else if effect == split {
		return triangles.WriteEffectPNG(filename, points, img, scale, false)
	}

	return errors.New("invalid effect")
}

func (t TriangleLogic) NewAlgorithm(img image.Data, mutations int, mutationAmount float64, numPoints, population, cutoff, blockSize, cacheSize int) algorithm.Algorithm {
	evaluatorFactory := func(n int) evaluator.Evaluator {
		return evaluator.NewParallel(fitness.TrianglesImageFunctions(img, blockSize, n), cacheSize)
	}
	var mutator mutation.Method
	mutator = mutation.NewGaussianMethod(float64(mutations)/float64(numPoints), mutationAmount)

	pointFactory := func() normgeom.NormPointGroup {
		return generator.RandomGenerator{}.Generate(numPoints)
	}

	algo := algorithm.NewModifiedGenetic(pointFactory, population, cutoff, evaluatorFactory, mutator)

	return algo
}

func (t TriangleLogic) RenderData(points normgeom.NormPointGroup, img image.Data) RenderData {
	w, h := img.Size()
	triangles := triangulation.Triangulate(points, w, h)
	triangleData := render.TrianglesOnImage(triangles, img)

	data := RenderData{
		Width:    w,
		Height:   h,
		Polygons: make([]normgeom.NormPolygon, len(triangleData)),
		Colors:   make([]color.RGB, len(triangleData)),
	}

	for i, d := range triangleData {
		data.Colors[i] = d.Color
		tri := d.Triangle.Points
		data.Polygons[i] = normgeom.NormPolygon{Points: []normgeom.NormPoint{tri[0], tri[1], tri[2]}}
	}

	return data
}

type PolygonLogic struct {
}

func (p PolygonLogic) SaveSVG(file string, points normgeom.NormPointGroup, img image.Data) error {
	filename := file
	if !strings.HasSuffix(filename, ".svg") {
		filename += ".svg"
	}

	return polygons.WriteSVG(filename, points, img)
}

func (p PolygonLogic) SavePNG(file string, points normgeom.NormPointGroup, img image.Data, scale float64, effect int) error {
	filename := file
	if !strings.HasSuffix(filename, ".png") {
		filename += ".png"
	}

	if effect == none {
		return polygons.WritePNG(filename, points, img, scale)
	} else if effect == gradient {
		return polygons.WriteEffectPNG(filename, points, img, scale, true)
	} else if effect == split {
		return polygons.WriteEffectPNG(filename, points, img, scale, false)
	}

	return errors.New("invalid effect")
}

func (p PolygonLogic) NewAlgorithm(img image.Data, mutations int, mutationAmount float64, numPoints, population, cutoff, blockSize, cacheSize int) algorithm.Algorithm {
	evaluatorFactory := func(n int) evaluator.Evaluator {
		return evaluator.NewParallel(fitness.PolygonsImageFunctions(img, blockSize, n), cacheSize)
	}
	var mutator mutation.Method
	mutator = mutation.NewGaussianMethod(float64(mutations)/float64(numPoints), mutationAmount)

	pointFactory := func() normgeom.NormPointGroup {
		return generator.RandomGenerator{}.Generate(numPoints)
	}

	algo := algorithm.NewModifiedGenetic(pointFactory, population, cutoff, evaluatorFactory, mutator)

	return algo
}

func (p PolygonLogic) RenderData(points normgeom.NormPointGroup, img image.Data) RenderData {
	w, h := img.Size()
	polygons := polygonation.Polygonate(points, w, h)
	triangleData := render.PolygonsOnImage(polygons, img)

	data := RenderData{
		Width:    w,
		Height:   h,
		Polygons: make([]normgeom.NormPolygon, len(triangleData)),
		Colors:   make([]color.RGB, len(triangleData)),
	}

	for i, d := range triangleData {
		data.Colors[i] = d.Color
		data.Polygons[i] = d.Polygon
	}

	return data
}
