package quadtree

import "github.com/snowzhop/verlet-fireplace/internal/math"

type Particle interface {
	Position() math.Vec2
	Intersects(p Particle) bool
	Side() float64
}
