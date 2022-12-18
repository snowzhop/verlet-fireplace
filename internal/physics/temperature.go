package physics

import "github.com/snowzhop/verlet-fireplace/internal/math"

const (
	MaxTemperature = 1000

	maxRaiseForce = 10
)

func RaiseForce(t float64) math.Vec2 {
	raiseForce := t / MaxTemperature * maxRaiseForce

	return math.Vec2{
		X: 0,
		Y: -raiseForce,
	}
}
