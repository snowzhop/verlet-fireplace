package geom

import "github.com/snowzhop/verlet-fireplace/internal/physics"

type position struct {
	x, y uint32
}

type Cell struct {
	position
	objects map[uint64]*physics.VerletObject
}

func NewCell(x, y uint32) *Cell {
	return &Cell{
		position: position{x: x, y: y},
		objects:  make(map[uint64]*physics.VerletObject),
	}
}

func (c *Cell) Add(obj *physics.VerletObject) {
	c.objects[obj.ID()] = obj
}

func (c *Cell) Delete(obj *physics.VerletObject) {
	delete(c.objects, obj.ID())
}
