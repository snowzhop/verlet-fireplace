package app

import (
	"github.com/snowzhop/verlet-fireplace/internal/math"
)

var (
	oldGravity math.Vec2
)

func (f *Fireplace) applyGravity() {
	for _, obj := range f.movableObjects {
		if oldGravity != f.gravity {
			oldGravity = f.gravity
		}

		obj.Accelerate(f.gravity)
	}
}

func (f *Fireplace) updatePositions(dt float64) {
	for _, obj := range f.movableObjects {
		obj.UpdatePosition(dt)
	}
}
