package math

import "math/rand"

func RandomOffset() (float64, float64) {
	var a, b float64
	for a == b {
		a = rand.Float64()*20 - 10
		b = rand.Float64()*20 - 10
	}
	return a, b
}

func RandomFloat64(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}
