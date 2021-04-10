package export

import "math"

func scale(num float64, d int) int {
	return int(math.Round(num * float64(d)))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func multAndRound(v int, s float64) int {
	return int(math.Round(float64(v) * s))
}