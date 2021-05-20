package util

import (
	"math"
	"strings"
)

const SvgStart = `<?xml version="1.0"?>
<svg width="%v" height="%v"
     xmlns="http://www.w3.org/2000/svg"
     shape-rendering="crispEdges">
`

//https://stackoverflow.com/a/25959527/15283541
var types = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
}

func MimeType(incipit []byte) string {
	incipitStr := string(incipit)
	for magic, mime := range types {
		if strings.HasPrefix(incipitStr, magic) {
			return mime
		}
	}

	return ""
}

func Scale(num float64, d int) int {
	return int(math.Round(num * float64(d)))
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MultAndRound(v int, s float64) int {
	return int(math.Round(float64(v) * s))
}
