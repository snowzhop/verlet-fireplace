package physics

import (
	"image/color"

	"github.com/snowzhop/verlet-fireplace/internal/math"
)

type VerletObject struct {
	CurrentPosition math.Vec2
	OldPosition     math.Vec2
	Acceleration    math.Vec2

	Radius float64
	Color  color.Color
}

func NewVerletObject(startPos math.Vec2, radius float64) *VerletObject {
	return &VerletObject{
		CurrentPosition: startPos,
		OldPosition:     startPos,
		Radius:          radius,
		Color:           color.White,
	}
}

func (v *VerletObject) UpdatePosition(dt float64) {
	velocity := math.SubVec2(v.CurrentPosition, v.OldPosition)

	v.OldPosition = v.CurrentPosition

	// CurrentPosition = CurrentPosition + velocity + acceleration * dt^2
	v.CurrentPosition = math.SumVec2(
		v.CurrentPosition,
		velocity,
		math.ApplyVec2(v.Acceleration, dt*dt),
	)

	v.Acceleration = math.Vec2{}
}

func (v *VerletObject) Accelerate(acc math.Vec2) {
	v.Acceleration = math.SumVec2(v.Acceleration, acc)
}
