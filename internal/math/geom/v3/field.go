package geom

import (
	"fmt"

	"github.com/dhconnelly/rtreego"
	"github.com/snowzhop/verlet-fireplace/internal/physics"
)

type Field struct {
	mainTree *rtreego.Rtree
}

func NewField(min, max int) *Field {
	tree := rtreego.NewTree(2, min, max)

	return &Field{
		mainTree: tree,
	}
}

func (f *Field) Insert(obj *physics.VerletObject) {
	f.mainTree.Insert(obj)
}

func (f *Field) Delete(obj *physics.VerletObject) {
	f.mainTree.Delete(obj)
}

func (f *Field) Neighbors(obj *physics.VerletObject) []*physics.VerletObject {
	rawNeighbors := f.mainTree.SearchIntersect(obj.Bounds())

	result := make([]*physics.VerletObject, 0, len(rawNeighbors))
	for i := range rawNeighbors {
		n, ok := rawNeighbors[i].(*physics.VerletObject)
		if !ok {
			panic(fmt.Errorf("wrong object inside rtree: %v", rawNeighbors[i]))
		}
		result = append(result, n)
	}

	return result
}

func (f *Field) UpdateObject(obj *physics.VerletObject) {
	// empty
}

func (f *Field) Dump() {
	fmt.Printf("\n%s\n", f.mainTree)
}
