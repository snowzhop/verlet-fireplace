package physics

import (
	"testing"

	"github.com/snowzhop/verlet-fireplace/internal/math"
	"github.com/stretchr/testify/assert"
)

func TestVerletObjectIntesects(t *testing.T) {
	testCases := []struct {
		name string
		in1  VerletObject
		in2  VerletObject
		out  bool
	}{
		{
			name: "intersects 1",
			in1: VerletObject{
				CurrentPosition: math.Vec2{X: 0, Y: 0},
				radius:          4,
			},
			in2: VerletObject{
				CurrentPosition: math.Vec2{X: 1, Y: 1},
				radius:          4,
			},
			out: true,
		},
		{
			name: "not intersect 1",
			in1: VerletObject{
				CurrentPosition: math.Vec2{X: 0, Y: 0},
				radius:          1,
			},
			in2: VerletObject{
				CurrentPosition: math.Vec2{X: 15, Y: 10},
				radius:          1,
			},
			out: false,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.in1.Intersects(&tc.in2), tc.out, tc.name)
	}
}
