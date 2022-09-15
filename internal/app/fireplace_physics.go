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

func (f *Fireplace) solveCollisions() {
	objectCount := len(f.movableObjects)
	for i := 0; i < objectCount; i++ {
		obj1 := f.movableObjects[i]
		for k := i + 1; k < objectCount; k++ {
			obj2 := f.movableObjects[k]
			collisionAxis := math.SubVec2(obj1.CurrentPosition, obj2.CurrentPosition)
			dist := collisionAxis.Len()
			if dist < obj1.Radius+obj2.Radius {
				n := math.ApplyVec2(collisionAxis, 1/dist)
				delta := obj1.Radius + obj2.Radius - dist

				obj1.CurrentPosition = math.SumVec2(
					obj1.CurrentPosition,
					math.ApplyVec2(n, float64(0.5)*delta),
				)
				obj2.CurrentPosition = math.SubVec2(
					obj2.CurrentPosition,
					math.ApplyVec2(n, float64(0.5)*delta),
				)
			}
		}
	}
}
