package math

import (
	stdmath "math"
)

type Gauss struct {
	mean         float64
	stdDeviation float64
}

func NewGauss(mean, stdDeviation float64) Gauss {
	return Gauss{
		mean:         mean,
		stdDeviation: stdDeviation,
	}
}

func (g Gauss) Calculate(x float64) float64 {
	return gauss(x, g.mean, g.stdDeviation)
}

func gauss(x, mean, stdDeviation float64) float64 {
	pow := stdmath.Pow((x-mean)/stdDeviation, 2) * (-0.5)
	return stdmath.Exp(pow) / (stdDeviation * stdmath.Sqrt2 * stdmath.SqrtPi)
}

func Sigmoid(x float64) float64 {
	return 1 / (1 + stdmath.Exp(-x))
}

func Linear(x, k float64) float64 {
	return k * x
}
