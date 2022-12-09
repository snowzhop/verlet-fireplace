package geom

import (
	"testing"

	"github.com/snowzhop/verlet-fireplace/internal/math"
)

func TestRect(t *testing.T) {
	testcases := []struct {
		A   Rect
		B   Rect
		Res bool
	}{
		{
			A:   Rect{UpperLeftCorner: math.Vec2{X: 3, Y: 6}, LowerRightCorner: math.Vec2{X: 4, Y: -12}},
			B:   Rect{UpperLeftCorner: math.Vec2{X: 0, Y: 0}, LowerRightCorner: math.Vec2{X: 10, Y: -10}},
			Res: false,
		},
		{
			A:   Rect{UpperLeftCorner: math.Vec2{X: -1, Y: -1}, LowerRightCorner: math.Vec2{X: 4, Y: -8}},
			B:   Rect{UpperLeftCorner: math.Vec2{X: -1, Y: -1}, LowerRightCorner: math.Vec2{X: 4, Y: -8}},
			Res: true,
		},
	}

	for _, tc := range testcases {
		if (tc.A == tc.B) != tc.Res {
			t.Errorf("A (%v) == B (%v). Real result: %v", tc.A, tc.B, tc.Res)
		}
	}
}
