package geom

import (
	"github.com/snowzhop/verlet-fireplace/internal/math"
	"github.com/snowzhop/verlet-fireplace/internal/physics"
)

type Field struct {
	Cells    [][]*Cell
	CellsMap map[Rect]*Cell

	cellSideLen float64
}

type Rect struct {
	UpperLeftCorner  math.Vec2
	LowerRightCorner math.Vec2
}

type Cell struct {
	Objects  map[uint64]*physics.VerletObject
	Position Rect
}

func (f *Field) Neighbors(cell *Cell) []*Cell {
	if f == nil {
		return nil
	}
	if cell == nil {
		return nil
	}

	upperLeftLeftX := cell.Position.UpperLeftCorner.X - f.cellSideLen
	lowerRightLeftX := cell.Position.LowerRightCorner.X - f.cellSideLen
	upperLeftRightX := cell.Position.UpperLeftCorner.X + f.cellSideLen
	lowerRightRigthX := cell.Position.LowerRightCorner.X + f.cellSideLen

	upperLeftUpperY := cell.Position.UpperLeftCorner.Y - f.cellSideLen
	lowerRightUpperY := cell.Position.LowerRightCorner.Y - f.cellSideLen
	upperLeftLowerY := cell.Position.UpperLeftCorner.Y + f.cellSideLen
	lowerRightLowerY := cell.Position.LowerRightCorner.Y + f.cellSideLen

	var res []*Cell
	// |.|.|.|
	// |*|o|.|
	// |.|.|.|
	if cell, ok := f.CellsMap[Rect{
		UpperLeftCorner:  math.Vec2{X: upperLeftLeftX, Y: cell.Position.UpperLeftCorner.Y},
		LowerRightCorner: math.Vec2{X: lowerRightLeftX, Y: cell.Position.LowerRightCorner.Y},
	}]; ok {
		res = append(res, cell)
	}
	// |*|.|.|
	// |.|o|.|
	// |.|.|.|
	if cell, ok := f.CellsMap[Rect{
		UpperLeftCorner:  math.Vec2{X: upperLeftLeftX, Y: upperLeftUpperY},
		LowerRightCorner: math.Vec2{X: lowerRightLeftX, Y: lowerRightUpperY},
	}]; ok {
		res = append(res, cell)
	}
	// |.|*|.|
	// |.|o|.|
	// |.|.|.|
	if cell, ok := f.CellsMap[Rect{
		UpperLeftCorner:  math.Vec2{X: cell.Position.UpperLeftCorner.X, Y: upperLeftUpperY},
		LowerRightCorner: math.Vec2{X: cell.Position.LowerRightCorner.X, Y: lowerRightUpperY},
	}]; ok {
		res = append(res, cell)
	}
	// |.|.|*|
	// |.|o|.|
	// |.|.|.|
	if cell, ok := f.CellsMap[Rect{
		UpperLeftCorner:  math.Vec2{X: upperLeftRightX, Y: upperLeftUpperY},
		LowerRightCorner: math.Vec2{X: lowerRightRigthX, Y: lowerRightUpperY},
	}]; ok {
		res = append(res, cell)
	}
	// |.|.|.|
	// |.|o|*|
	// |.|.|.|
	if cell, ok := f.CellsMap[Rect{
		UpperLeftCorner:  math.Vec2{X: upperLeftRightX, Y: cell.Position.UpperLeftCorner.Y},
		LowerRightCorner: math.Vec2{X: lowerRightRigthX, Y: cell.Position.LowerRightCorner.Y},
	}]; ok {
		res = append(res, cell)
	}
	// |.|.|.|
	// |.|o|.|
	// |.|.|*|
	if cell, ok := f.CellsMap[Rect{
		UpperLeftCorner:  math.Vec2{X: upperLeftRightX, Y: upperLeftLowerY},
		LowerRightCorner: math.Vec2{X: lowerRightRigthX, Y: lowerRightLowerY},
	}]; ok {
		res = append(res, cell)
	}
	// |.|.|.|
	// |.|o|.|
	// |.|*|.|
	if cell, ok := f.CellsMap[Rect{
		UpperLeftCorner:  math.Vec2{X: cell.Position.UpperLeftCorner.X, Y: upperLeftLowerY},
		LowerRightCorner: math.Vec2{X: cell.Position.LowerRightCorner.X, Y: lowerRightLowerY},
	}]; ok {
		res = append(res, cell)
	}
	// |.|.|.|
	// |.|o|.|
	// |*|.|.|
	if cell, ok := f.CellsMap[Rect{
		UpperLeftCorner:  math.Vec2{X: upperLeftLeftX, Y: upperLeftLowerY},
		LowerRightCorner: math.Vec2{X: lowerRightLeftX, Y: lowerRightLowerY},
	}]; ok {
		res = append(res, cell)
	}

	return res
}
