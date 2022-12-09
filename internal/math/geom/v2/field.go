package geom

import (
	"fmt"

	"github.com/snowzhop/verlet-fireplace/internal/math"
	"github.com/snowzhop/verlet-fireplace/internal/physics"
)

type Field struct {
	Map          [][]*Cell
	objectsCells map[uint64]*Cell
	maxX         uint32
	maxY         uint32
	cellSideLen  uint32
}

func neighborField(r uint32) map[[2]int]struct{} {
	res := make(map[[2]int]struct{})

	for i := 0; i <= int(r); i++ {
		for j := 0; j <= int(r); j++ {
			if i == 0 && j == 0 {
				continue
			}
			res[[2]int{i, j}] = struct{}{}
			res[[2]int{i, -j}] = struct{}{}
			res[[2]int{-i, j}] = struct{}{}
			res[[2]int{-i, -j}] = struct{}{}
		}
	}
	return res
}

var neighborsOffsets = []*math.Vec2{
	// inner circle
	{X: -1, Y: 0},
	{X: -1, Y: -1},
	{X: 0, Y: -1},
	{X: 1, Y: -1},
	{X: 1, Y: 0},
	{X: 1, Y: 1},
	{X: 0, Y: 1},
	{X: -1, Y: 1},
	// outer circle
	{X: -2, Y: 0},
	{X: -2, Y: -1},
	{X: -2, Y: -2},
	{X: -1, Y: -2},
	{X: 0, Y: -2},
	{X: 1, Y: -2},
	{X: 2, Y: -2},
	{X: 2, Y: -1},
	{X: 2, Y: 0},
	{X: 2, Y: 1},
	{X: 2, Y: 2},
	{X: 1, Y: 2},
	{X: 0, Y: 2},
	{X: -1, Y: 2},
	{X: -2, Y: 2},
	{X: -2, Y: 1},
}

func NewField(width, height, cellSideLen uint32) *Field {
	maxX := width / cellSideLen
	maxY := height / cellSideLen

	m := make([][]*Cell, maxX)
	for i := range m {
		m[i] = make([]*Cell, maxY)
		for j := range m[i] {
			m[i][j] = NewCell(uint32(i), uint32(j))
		}
	}
	objects := make(map[uint64]*Cell, maxX+maxY)

	return &Field{
		Map:          m,
		objectsCells: objects,
		maxX:         maxX,
		maxY:         maxY,
		cellSideLen:  cellSideLen,
	}
}

func (f *Field) Insert(obj *physics.VerletObject) {
	x := int(obj.CurrentPosition.X / float64(f.cellSideLen))
	y := int(obj.CurrentPosition.Y / float64(f.cellSideLen))

	cell := f.Map[x][y]
	cell.Add(obj)
	f.objectsCells[obj.ID()] = cell
}

func (f *Field) Neighbors(obj *physics.VerletObject) []*physics.VerletObject {
	var result []*physics.VerletObject

	cell := f.objectsCells[obj.ID()]

	for _, offset := range neighborsOffsets {
		x := int(cell.x) + int(offset.X)
		switch {
		case x < 0:
			continue
		case x >= int(f.maxX):
			continue
		}
		y := int(cell.y) + int(offset.Y)
		switch {
		case y < 0:
			continue
		case y >= int(f.maxY):
			continue
		}

		n := f.Map[x][y]
		for _, rObj := range n.objects {
			result = append(result, rObj)
		}
	}
	return result
}

func (f *Field) UpdateObject(obj *physics.VerletObject) {
	x := int(obj.CurrentPosition.X / float64(f.cellSideLen))
	if x >= int(f.maxX) {
		x = int(f.maxX) - 1
	}
	if x < 0 {
		x = 0
	}
	y := int(obj.CurrentPosition.Y / float64(f.cellSideLen))
	if y >= int(f.maxY) {
		y = int(f.maxY) - 1
	}
	if y < 0 {
		y = 0
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("\nobj before panic: %v\n", obj)
			fmt.Printf("x: %v\ny: %v\n\n", x, y)
			panic(err)
		}
	}()

	oldCell := f.objectsCells[obj.ID()]
	newCell := f.Map[x][y]
	if oldCell == newCell {
		return
	}

	oldCell.Delete(obj)
	delete(f.objectsCells, obj.ID())

	if newCell != nil {
		newCell.Add(obj)
	}
	f.objectsCells[obj.ID()] = newCell
}

func (f *Field) Dump() {
	var count int
	for _, row := range f.Map {
		for _, cell := range row {
			if cell != nil {
				count += len(cell.objects)
			}
		}
	}

	fmt.Printf("Cell count on map: %d\n", count)
	fmt.Printf("len(objectCells): %d\n", len(f.objectsCells))
	fmt.Printf("maxX: %d\n", f.maxX)
	fmt.Printf("maxY: %d\n", f.maxY)
	fmt.Printf("cellSideLen: %d\n", f.cellSideLen)
}
