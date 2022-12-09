package app

import (
	"fmt"
	"sync"

	"github.com/snowzhop/verlet-fireplace/internal/math"
	"github.com/snowzhop/verlet-fireplace/internal/math/quadtree"
	"github.com/snowzhop/verlet-fireplace/internal/physics"
)

const (
	maxProcNum = 32
)

func (f *Fireplace) applyGravity() {
	for _, obj := range f.movableObjects {
		obj.Accelerate(f.gravity)
	}
}

func (f *Fireplace) applyForces() {
	for _, obj := range f.movableObjects {
		raiseForce := physics.RaiseForce(obj.Temperature())

		obj.Accelerate(
			math.SumVec2(
				f.gravity,
				raiseForce,
			),
		)
		// t := math.Sigmoid(obj.Temperature()/physics.MaxTemperature*12-6) * obj.Temperature()
		t := math.Linear(obj.Temperature(), 1.5) / physics.MaxTemperature * obj.Temperature()
		obj.IncreaseTemperature(-t)

		if obj.CurrentPosition.IsNaN() {
			panic(fmt.Sprintf("applyForces: pos is nan: %v", *obj))
		}
	}
}

func (f *Fireplace) updatePositions(dt float64) {
	for _, obj := range f.movableObjects {
		obj.UpdatePosition(dt)
		f.field.UpdateObject(obj)
	}
}

func (f *Fireplace) updatePositions2(dt float64) {
	for _, obj := range f.movableObjects {
		obj.UpdatePosition(dt)
		if obj.CurrentPosition.IsNaN() {
			panic(fmt.Sprintf("updatePositions2: pos is nan: %v", obj))
		}
	}
}

func (f *Fireplace) applyAllConstraints() {
	for _, obj := range f.movableObjects {
		// toObj := math.SubVec2(obj.CurrentPosition, f.staticMainConstraint.Position)
		// dist := toObj.Len()

		// if dist > f.staticMainConstraint.Radius-obj.Radius {
		// 	n := math.ApplyVec2(toObj, float64(1)/dist)
		// 	obj.CurrentPosition = math.SumVec2(
		// 		f.staticMainConstraint.Position,
		// 		math.ApplyVec2(n, f.staticMainConstraint.Radius-obj.Radius),
		// 	)
		// }

		f.applyWorldBoxConstraintDirectly(obj)

		if obj.CurrentPosition.IsNaN() {
			panic(fmt.Sprintf("applyConstraint: pos is nan: %v", *obj))
		}
	}
}

func (f *Fireplace) applyWorldBoxConstraintDirectly(obj *physics.VerletObject) {
	if obj.CurrentPosition.X < obj.Radius() {
		obj.CurrentPosition.X = obj.Radius()
	}
	if obj.CurrentPosition.X > float64(f.game.screenWidth)-obj.Radius() {
		obj.CurrentPosition.X = float64(f.game.screenWidth) - obj.Radius()
	}

	if obj.CurrentPosition.Y > float64(f.game.screenHeight)-obj.Radius() {
		obj.CurrentPosition.Y = float64(f.game.screenHeight) - obj.Radius()
	}
	if obj.CurrentPosition.Y < obj.Radius() {
		obj.CurrentPosition.Y = obj.Radius()
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
			if dist < obj1.Radius()+obj2.Radius() {
				n := math.ApplyVec2(collisionAxis, 1/dist)
				delta := obj1.Radius() + obj2.Radius() - dist
				correction := math.ApplyVec2(n, float64(0.5)*delta)

				obj1.CurrentPosition = math.SumVec2(
					obj1.CurrentPosition,
					correction,
				)
				obj2.CurrentPosition = math.SubVec2(
					obj2.CurrentPosition,
					correction,
				)

				// Temperature
				tStep := 0.5
				if obj1.Temperature() < obj2.Temperature() {
					obj1.IncreaseTemperature(tStep)
					obj2.IncreaseTemperature(-tStep)
				} else {
					obj1.IncreaseTemperature(-tStep)
					obj2.IncreaseTemperature(tStep)
				}
			}
		}
	}
}

func (f *Fireplace) solveCollisions2() {
	for _, obj1 := range f.movableObjects {

		neighbors := f.field.Neighbors(obj1)
		for _, obj2 := range neighbors {
			collisionAxis := math.SubVec2(obj1.CurrentPosition, obj2.CurrentPosition)
			dist := collisionAxis.Len()
			if dist < obj1.Radius()+obj2.Radius() {
				n := math.ApplyVec2(collisionAxis, 1/dist)
				delta := obj1.Radius() + obj2.Radius() - dist
				correction := math.ApplyVec2(n, float64(0.5)*delta)

				obj1.CurrentPosition = math.SumVec2(
					obj1.CurrentPosition,
					correction,
				)
				obj2.CurrentPosition = math.SubVec2(
					obj2.CurrentPosition,
					correction,
				)

				// Temperature
				tStep := 0.5
				if obj1.Temperature() < obj2.Temperature() {
					obj1.IncreaseTemperature(tStep)
					obj2.IncreaseTemperature(-tStep)
				} else {
					obj1.IncreaseTemperature(-tStep)
					obj2.IncreaseTemperature(tStep)
				}
			}
		}
	}
}

func (f *Fireplace) solveCollisions3() {
	for _, obj1 := range f.movableObjects {
		neighbors := f.field.Neighbors(obj1)
		for _, obj2 := range neighbors {
			collisionAxis := math.SubVec2(obj1.CurrentPosition, obj2.CurrentPosition)
			dist := collisionAxis.Len()
			if dist < obj1.Radius()+obj2.Radius() {

				n := math.ApplyVec2(collisionAxis, 1/dist)
				delta := obj1.Radius() + obj2.Radius() - dist
				correction := math.ApplyVec2(n, float64(0.5)*delta)

				obj1.CurrentPosition = math.SumVec2(
					obj1.CurrentPosition,
					correction,
				)
				obj2.CurrentPosition = math.SubVec2(
					obj2.CurrentPosition,
					correction,
				)

				// Temperature
				tStep := 0.5
				if obj1.Temperature() < obj2.Temperature() {
					obj1.IncreaseTemperature(tStep)
					obj2.IncreaseTemperature(-tStep)
				} else {
					obj1.IncreaseTemperature(-tStep)
					obj2.IncreaseTemperature(tStep)
				}
			}
		}
	}
}

func (f *Fireplace) solveCollisions4() {
	intersections := f.root.FindIntersections()
	for _, pair := range intersections {
		obj1, ok := pair.First.(*physics.VerletObject)
		if !ok {
			panic(fmt.Sprintf("failed to type cast obj1 (%T)", pair.First))
		}
		obj2, ok := pair.Second.(*physics.VerletObject)
		if !ok {
			panic(fmt.Sprintf("failed to type cast obj2 (%T)", pair.Second))
		}

		collisionAxis := math.SubVec2(obj1.CurrentPosition, obj2.CurrentPosition)
		dist := collisionAxis.Len()

		n := math.ApplyVec2(collisionAxis, 1/dist)
		delta := obj1.Radius() + obj2.Radius() - dist
		correction := math.ApplyVec2(n, float64(0.5)*delta)

		obj1.CurrentPosition = math.SumVec2(
			obj1.CurrentPosition,
			correction,
		)
		obj2.CurrentPosition = math.SubVec2(
			obj2.CurrentPosition,
			correction,
		)

		// Temperature
		tStep := 0.5
		if obj1.Temperature() < obj2.Temperature() {
			obj1.IncreaseTemperature(tStep)
			obj2.IncreaseTemperature(-tStep)
		} else {
			obj1.IncreaseTemperature(-tStep)
			obj2.IncreaseTemperature(tStep)
		}
	}
}

func (f *Fireplace) solveCollisions5() {
	for _, obj1 := range f.movableObjects {
		intersections := f.root.Query(obj1)
		for _, rawObj2 := range intersections {
			obj2, ok := rawObj2.(*physics.VerletObject)
			if !ok {
				panic(fmt.Sprintf("failed to convert obj2 (%T) to VerletObject: %[1]v", obj2))
			}

			collisionAxis := math.SubVec2(obj1.CurrentPosition, obj2.CurrentPosition)
			dist := collisionAxis.Len()

			n := math.ApplyVec2(collisionAxis, 1/dist)
			delta := obj1.Radius() + obj2.Radius() - dist
			correction := math.ApplyVec2(n, float64(0.5)*delta)

			obj1.CurrentPosition = math.SumVec2(
				obj1.CurrentPosition,
				correction,
			)
			obj2.CurrentPosition = math.SubVec2(
				obj2.CurrentPosition,
				correction,
			)

			// Temperature
			tStep := 0.5
			if obj1.Temperature() < obj2.Temperature() {
				obj1.IncreaseTemperature(tStep)
				obj2.IncreaseTemperature(-tStep)
			} else {
				obj1.IncreaseTemperature(-tStep)
				obj2.IncreaseTemperature(tStep)
			}
		}
	}
}

func (f *Fireplace) solveCollisions6() {
	var (
		objSliceLen  = len(f.movableObjects)
		numOfThreads = computeProcessorNum(objSliceLen)

		step  = objSliceLen / numOfThreads
		start = -step
		end   = 0

		wg sync.WaitGroup
	)

	proc := func(objects []*physics.VerletObject) {
		for _, obj1 := range objects {
			intersections := f.root.Query(obj1)
			for _, rawObj2 := range intersections {
				obj2, ok := rawObj2.(*physics.VerletObject)
				if !ok {
					panic(fmt.Sprintf("failed to convert obj2 (%T) to VerletObject at proc: %[1]v", obj2))
				}

				collisionAxis := math.SubVec2(obj1.CurrentPosition, obj2.CurrentPosition)
				dist := collisionAxis.Len()

				if dist == 0 {
					obj1.CurrentPosition = alignObject(obj1)
					obj2.CurrentPosition = alignObject(obj2)
				} else {
					n := math.ApplyVec2(collisionAxis, 1/dist)
					delta := obj1.Radius() + obj2.Radius() - dist
					correction := math.ApplyVec2(n, float64(0.5)*delta)

					obj1.CurrentPosition = math.SumVec2(
						obj1.CurrentPosition,
						correction,
					)
					obj2.CurrentPosition = math.SubVec2(
						obj2.CurrentPosition,
						correction,
					)
				}

				// Temperature
				if obj1.Temperature() < obj2.Temperature() {
					obj1.IncreaseTemperature(f.game.temperatureStep)
					obj2.IncreaseTemperature(-f.game.temperatureStep)
				} else {
					obj1.IncreaseTemperature(-f.game.temperatureStep)
					obj2.IncreaseTemperature(f.game.temperatureStep)
				}
			}
		}
	}

	for end < len(f.movableObjects) {
		start += step
		end += step
		if objSliceLen-end < step {
			end = objSliceLen
		}

		wg.Add(1)
		go func(s, e int) {
			defer wg.Done()

			proc(f.movableObjects[s:e])
		}(start, end)
	}

	wg.Wait()

	// ---------------------
	hiddenParticles := make([]*physics.VerletObject, 0, len(f.movableObjects)+len(f.hiddenObjects))
	particles := make(map[math.Vec2]int) // position -> index
	for i, obj := range f.movableObjects {
		var (
			ok bool
		)
		if _, ok = particles[obj.CurrentPosition]; !ok {
			particles[obj.CurrentPosition] = i
			continue
		}

		fmt.Println("COLLISION!!!")

		index := i - len(hiddenParticles)
		f.movableObjects = append(f.movableObjects[:index], f.movableObjects[index+1:]...)
		hiddenParticles = append(hiddenParticles, obj)
		// obj.IncreaseTemperature(-physics.MaxTemperature)
	}

	for _, obj1 := range f.hiddenObjects {
		ok := f.canHiddenObjectBeRestored(obj1)
		if ok {
			f.movableObjects = append(f.movableObjects, obj1)
		} else {
			hiddenParticles = append(hiddenParticles, obj1)
		}
	}
	f.hiddenObjects = hiddenParticles

	// ---------------------
}

func (f *Fireplace) canHiddenObjectBeRestored(obj *physics.VerletObject) bool {
	intersections := f.root.Query(obj)
	for _, rawObj2 := range intersections {
		obj2, ok := rawObj2.(*physics.VerletObject)
		if !ok {
			panic(fmt.Sprintf("failed to convert obj2 (%T) to VerletObject(hidden): %[1]v", obj2))
		}
		if obj.CurrentPosition == obj2.CurrentPosition {
			return false
		}
	}
	return true
}

func (f *Fireplace) rebuildTree() {
	f.root.Clear()
	// f.root = quadtree.New(float64(f.game.screenWidth))
	f.root = quadtree.NewWithStart(-rootOffset, -rootOffset, float64(f.game.screenWidth)+2*rootOffset)
	for _, obj := range f.movableObjects {
		f.root.Insert(obj)
	}
}

func (f *Fireplace) applyHeat() {
	for _, heatEmitter := range f.heatEmitters {
		heatedObjects := f.root.Query(heatEmitter)
		for _, rawObj := range heatedObjects {
			obj, ok := rawObj.(*physics.VerletObject)
			if !ok {
				panic(fmt.Sprintf("applyHeat: failed to type cast obj %T", rawObj))
			}

			obj.IncreaseTemperature(heatEmitter.Temperature() * f.game.heatEmitterEfficiency)
			if obj.CurrentPosition.IsNaN() {
				panic(fmt.Sprintf("applyHeat: pos is nan: %v", *obj))
			}
		}
	}
}

func alignObject(obj *physics.VerletObject) math.Vec2 {
	actionAxis := math.SubVec2(obj.OldPosition, obj.CurrentPosition)
	length := actionAxis.Len()
	if length == 0 {
		return obj.CurrentPosition
	}

	n := math.ApplyVec2(actionAxis, 1/length)
	delta := length - obj.Radius()*2
	newPos := math.SubVec2(
		obj.CurrentPosition,
		math.ApplyVec2(n, -delta),
	)

	return newPos
}

func computeProcessorNum(l int) int {
	switch {
	case l <= 10:
		return 1
	case l > 10 && l <= 50:
		return 2
	case l > 50 && l <= 300:
		return 8
	case l > 300 && l <= 700:
		return 32
	case l > 700:
		return 32
	}
	return 1
}
