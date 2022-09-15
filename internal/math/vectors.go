package math

import "math"

type Vec2 struct {
	X, Y float64
}

func (v *Vec2) Len() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func SubVec2(a, b Vec2) Vec2 {
	return Vec2{
		X: a.X - b.X,
		Y: a.Y - b.Y,
	}
}

func SumVec2(vectors ...Vec2) Vec2 {
	var res Vec2

	for _, v := range vectors {
		res.X += v.X
		res.Y += v.Y
	}

	return res
}

func ApplyVec2(a Vec2, num float64) Vec2 {
	return Vec2{
		X: a.X * num,
		Y: a.Y * num,
	}
}
