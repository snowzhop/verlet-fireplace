package physics

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dhconnelly/rtreego"
	"github.com/snowzhop/verlet-fireplace/internal/math"
	"github.com/snowzhop/verlet-fireplace/internal/math/quadtree"
)

const (
	maxAcceleration = float64(3)
	maxVelocity     = float64(3)
	cuttedVelocity  = maxVelocity * 0.01
)

var (
	r           *rand.Rand
	idGenerator func() uint64
)

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixMicro()))

	var id uint64
	idGenerator = func() uint64 {
		id++
		return id
	}
}

type VerletObject struct {
	CurrentPosition math.Vec2
	OldPosition     math.Vec2
	Acceleration    math.Vec2

	radius      float64
	temperature float64

	Hidden bool

	id uint64
}

func NewVerletObject(startPos math.Vec2, radius float64) *VerletObject {
	t := r.Float64() * MaxTemperature

	return &VerletObject{
		CurrentPosition: startPos,
		OldPosition:     startPos,
		radius:          radius,
		temperature:     t,
		id:              idGenerator(),
	}
}

func NewVerletObjectWithTemp(startPos math.Vec2, radius float64, temp float64) *VerletObject {
	return &VerletObject{
		CurrentPosition: startPos,
		OldPosition:     startPos,
		radius:          radius,
		temperature:     temp,
		id:              idGenerator(),
	}
}

func (v *VerletObject) UpdatePosition(dt float64) {
	velocity := math.SubVec2(v.CurrentPosition, v.OldPosition)

	vLen := velocity.Len()
	if vLen > maxVelocity {
		velocity.X = velocity.X * cuttedVelocity / vLen
		velocity.Y = velocity.Y * cuttedVelocity / vLen

		v.CurrentPosition = math.SumVec2(v.OldPosition, velocity)
		// v.Hidden = true
	}

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

func (v *VerletObject) ID() uint64 {
	return v.id
}

func (v *VerletObject) Radius() float64 {
	return v.radius
}

func (v *VerletObject) Temperature() float64 {
	return v.temperature
}

func (v *VerletObject) IncreaseTemperature(tChange float64) {
	v.temperature += tChange

	if v.temperature < 1 {
		v.temperature = 1
	}
	if v.temperature > MaxTemperature {
		v.temperature = MaxTemperature
	}
}

func (v *VerletObject) Refresh(distribX float64) {
	curPos := math.Vec2{X: math.RandomFloat64(0, distribX), Y: v.radius}
	v.CurrentPosition = curPos
	v.OldPosition = curPos
	v.Acceleration = math.Vec2{}
	v.temperature = 1
	v.Hidden = false
}

func (v *VerletObject) SetRadius(r float64) {
	v.radius = r
}

func (v *VerletObject) Bounds() *rtreego.Rect {
	rect, err := rtreego.NewRect(
		rtreego.Point{v.CurrentPosition.X - v.radius/2, v.CurrentPosition.Y - v.radius/2},
		[]float64{v.radius, v.radius},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create new rect for object %d", v.id))
	}
	return rect
}

func (v *VerletObject) Position() math.Vec2 {
	return v.CurrentPosition
}

func (v *VerletObject) Intersects(p quadtree.Particle) bool {
	obj, ok := p.(*VerletObject)
	if !ok {
		panic("failed to cast interface to VerletObject")
	}

	collisionAxis := math.SubVec2(v.CurrentPosition, obj.CurrentPosition)
	dist := collisionAxis.Len()
	return dist <= obj.radius+v.radius
}

func (v *VerletObject) Side() float64 {
	return v.radius * 2
}

func (v *VerletObject) String() string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("{%v %v id:%d}", v.CurrentPosition, v.OldPosition, v.id)
}
