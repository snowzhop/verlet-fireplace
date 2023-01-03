package physics

import "github.com/snowzhop/verlet-fireplace/internal/math"

const (
	MaxTemperature = 1000

	maxRaiseForce = 2
)

func RaiseForce(t float64) math.Vec2 {
	raiseForce := t / MaxTemperature * maxRaiseForce

	return math.Vec2{
		X: 0,
		Y: -raiseForce,
	}
}

func NewHeatEmitters(width, step, baseRadius float64) (heatEmitters []*VerletObject) {
	for i := float64(0); i < width; i += baseRadius * step {
		heatEmitterRadius := math.NormFloat64() * 10
		if heatEmitterRadius >= baseRadius {
			heatEmitters = append(heatEmitters, NewVerletObjectWithTemp(
				math.Vec2{
					X: i,
					Y: width - heatEmitterRadius,
				},
				heatEmitterRadius,
				1000,
			))
		}
	}
	return
}
