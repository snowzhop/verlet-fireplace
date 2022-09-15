package math

import (
	"testing"
)

func TestSubVec2(t *testing.T) {
	testcases := []struct {
		a Vec2
		b Vec2
		c Vec2
	}{
		{
			a: Vec2{X: 1, Y: 2},
			b: Vec2{X: 1, Y: 1},
			c: Vec2{X: 0, Y: 1},
		},
		{
			a: Vec2{X: 3, Y: -1},
			b: Vec2{X: -1, Y: 4},
			c: Vec2{X: 4, Y: -5},
		},
	}

	for _, vectors := range testcases {
		c := SubVec2(vectors.a, vectors.b)
		if c.X != vectors.c.X || c.Y != vectors.c.Y {
			t.Errorf("got wrong vector: {%v} != {%v}", c, vectors.c)
		}
	}
}

func TestSumVec2(t *testing.T) {
	testcases := []struct {
		a Vec2
		b Vec2
		c Vec2
	}{
		{
			a: Vec2{X: 1, Y: 2},
			b: Vec2{X: 1, Y: 1},
			c: Vec2{X: 2, Y: 3},
		},
		{
			a: Vec2{X: 3, Y: -1},
			b: Vec2{X: -1, Y: 4},
			c: Vec2{X: 2, Y: 3},
		},
	}

	for _, vectors := range testcases {
		c := SumVec2(vectors.a, vectors.b)
		if c.X != vectors.c.X || c.Y != vectors.c.Y {
			t.Errorf("got wrong vector: {%v} != {%v}", c, vectors.c)
		}
	}
}

func TestApplyVec2(t *testing.T) {
	testcases := []struct {
		a   Vec2
		num float64
		c   Vec2
	}{
		{
			a:   Vec2{X: 1, Y: 2},
			num: 2,
			c:   Vec2{X: 2, Y: 4},
		},
		{
			a:   Vec2{X: 3, Y: -1},
			num: -3,
			c:   Vec2{X: -9, Y: 3},
		},
	}

	for _, vectors := range testcases {
		c := ApplyVec2(vectors.a, vectors.num)
		if c.X != vectors.c.X || c.Y != vectors.c.Y {
			t.Errorf("got wrong vector: {%v} != {%v}", c, vectors.c)
		}
	}
}
