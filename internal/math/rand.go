package math

import (
	"math/rand"
	"sync"
	"time"
)

var (
	onceRandomizerInitter sync.Once
	globalRandomizer      *rand.Rand
)

func getRandomizer() *rand.Rand {
	onceRandomizerInitter.Do(
		func() {
			globalRandomizer = rand.New(rand.NewSource(time.Now().UnixNano()))
		},
	)
	return globalRandomizer
}

func RandomOffset() (float64, float64) {
	randomizer := getRandomizer()
	var a, b float64
	for a == b {
		a = randomizer.Float64()*20 - 10
		b = randomizer.Float64()*20 - 10
	}
	return a, b
}

func RandomFloat64(min, max float64) float64 {
	randomizer := getRandomizer()
	return randomizer.Float64()*(max-min) + min
}

func NormFloat64() float64 {
	randomizer := getRandomizer()
	return randomizer.NormFloat64()
}
