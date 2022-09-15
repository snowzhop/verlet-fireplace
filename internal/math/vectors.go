package math

type Vec2 struct {
	X, Y float64
}

func SubVec2(a, b Vec2) Vec2 {
	return Vec2{
		X: a.X - b.X,
		Y: a.Y - b.Y,
	}
}

func SumVec2(a, b Vec2) Vec2 {
	return Vec2{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}

func ApplyVec2(a Vec2, num float64) Vec2 {
	return Vec2{
		X: a.X * num,
		Y: a.Y * num,
	}
}
