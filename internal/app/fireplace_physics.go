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

func (f *Fireplace) applyConstraint() {
	for _, obj := range f.movableObjects {
		toObj := math.SubVec2(obj.CurrentPosition, f.staticMainConstraint.Position)
		dist := toObj.Len()

		if dist > f.staticMainConstraint.Radius-obj.Radius {
			n := math.ApplyVec2(toObj, float64(1)/dist)
			obj.CurrentPosition = math.SumVec2(
				f.staticMainConstraint.Position,
				math.ApplyVec2(n, f.staticMainConstraint.Radius-obj.Radius),
			)
		}
	}
}
