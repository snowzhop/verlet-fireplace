package physics

import "github.com/snowzhop/verlet-fireplace/internal/math"

type VerletSolver struct {
	gravity math.Vec2
}

func (v *VerletSolver) Update(dt float64) {}

func (v *VerletSolver) ApplyGravity() {

}
