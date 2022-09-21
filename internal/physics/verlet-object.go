package physics

import (
	"math/rand"
	"time"

	"github.com/snowzhop/verlet-fireplace/internal/math"
)

const (
	MaxTemperature = 1000
)

var (
	r *rand.Rand
)

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixMicro()))
}

type VerletObject struct {
	CurrentPosition math.Vec2
	OldPosition     math.Vec2
	Acceleration    math.Vec2

	Radius      float64
	Temperature float64
}

func NewVerletObject(startPos math.Vec2, radius float64) *VerletObject {
	t := r.Float64() * MaxTemperature

	return &VerletObject{
		CurrentPosition: startPos,
		OldPosition:     startPos,
		Radius:          radius,
		Temperature:     t,
	}
}

func NewVerletObjectWithTemp(startPos math.Vec2, radius float64, temp float64) *VerletObject {
	return &VerletObject{
		CurrentPosition: startPos,
		OldPosition:     startPos,
		Radius:          radius,
		Temperature:     temp,
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
