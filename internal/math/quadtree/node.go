package quadtree

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/snowzhop/verlet-fireplace/internal/container"
	"github.com/snowzhop/verlet-fireplace/internal/math"
)

type NodeType byte

const (
	Leaf NodeType = iota
	Internal
)

type Quadrant int

const (
	NE Quadrant = iota
	NW
	SW
	SE
	UnknownQuadrant = -1
)

const (
	threshold = 8
)

type Node struct {
	QuadrantPos math.Vec2
	TotalPos    math.Vec2
	Count       uint32
	Width       float64

	Type      NodeType
	Children  []*Node
	Particles []Particle
}

func New(width float64) *Node {
	return newNode(0, 0, width)
}

func NewWithStart(x, y, width float64) *Node {
	return newNode(x, y, width)
}

func newNode(x, y, width float64) *Node {
	return &Node{
		QuadrantPos: math.Vec2{
			X: x,
			Y: y,
		},
		Children: make([]*Node, 4),
		Width:    width,
		Type:     Leaf,
	}
}

func (n *Node) defineQuadrant(pos math.Vec2) Quadrant {
	halfWidth := n.Width * 0.5
	// if pos.Y < n.QuadrantPos.Y+halfWidth {
	// 	if pos.X < n.QuadrantPos.X+halfWidth {
	// 		return NW
	// 	}
	// 	return NE
	// }
	// if pos.X < n.QuadrantPos.X+halfWidth {
	// 	return SW
	// }
	// return SE
	switch {
	case pos.X < n.QuadrantPos.X+halfWidth && pos.Y < n.QuadrantPos.Y+halfWidth:
		return NW
	case pos.X >= n.QuadrantPos.X+halfWidth && pos.Y < n.QuadrantPos.Y+halfWidth:
		return NE
	case pos.X < n.QuadrantPos.X+halfWidth && pos.Y >= n.QuadrantPos.Y+halfWidth:
		return SW
	case pos.X >= n.QuadrantPos.X+halfWidth && pos.Y >= n.QuadrantPos.Y+halfWidth:
		return SE
	}
	return UnknownQuadrant
}

func (n *Node) defineQuadrant2(p Particle) Quadrant {
	halfWidth := n.Width * 0.5
	pSide := p.Side()
	pPos := p.Position()

	if pPos.X+pSide < n.QuadrantPos.X+halfWidth {
		if pPos.Y+pSide < n.QuadrantPos.Y+halfWidth {
			return NW
		} else if pPos.Y >= n.QuadrantPos.Y+halfWidth {
			return SW
		}
	} else if pPos.X >= n.QuadrantPos.X+halfWidth {
		if pPos.Y+pSide < n.QuadrantPos.Y+halfWidth {
			return NE
		} else if pPos.Y >= n.QuadrantPos.Y+halfWidth {
			return SE
		}
	}

	return UnknownQuadrant
}

func (n *Node) split() {
	halfWidth := n.Width * 0.5
	n.Children[NE] = newNode(n.QuadrantPos.X+halfWidth, n.QuadrantPos.Y, halfWidth)
	n.Children[NW] = newNode(n.QuadrantPos.X, n.QuadrantPos.Y, halfWidth)
	n.Children[SW] = newNode(n.QuadrantPos.X, n.QuadrantPos.Y+halfWidth, halfWidth)
	n.Children[SE] = newNode(n.QuadrantPos.X+halfWidth, n.QuadrantPos.Y+halfWidth, halfWidth)
	n.Type = Internal
}

func (n *Node) Insert(p Particle) {
	if n.Type == Leaf {
		if len(n.Particles) < threshold {
			n.Particles = append(n.Particles, p)
		} else {
			n.splitLeaf()
			n.Insert(p)
		}
	} else {
		q := n.defineQuadrant2(p)
		if q != UnknownQuadrant {
			n.Children[q].Insert(p)
		} else {
			n.Particles = append(n.Particles, p)
		}
	}
}

func (n *Node) splitLeaf() {
	n.split()
	parentParticles := make([]Particle, 0, threshold)
	for _, p := range n.Particles {
		q := n.defineQuadrant2(p)
		if q != UnknownQuadrant {
			n.Children[q].Particles = append(n.Children[q].Particles, p)
		} else {
			parentParticles = append(parentParticles, p)
		}
	}
	n.Particles = parentParticles
}

func (n *Node) Dump() {
	// TODO: print tree info
}

func (n *Node) FindIntersections() []*container.Pair[Particle, Particle] {
	return findAllIntersections(n)
}

func findAllIntersections(node *Node) []*container.Pair[Particle, Particle] {
	var result []*container.Pair[Particle, Particle]
	for i := 0; i < len(node.Particles); i++ {
		for j := 0; j < i; j++ {
			if node.Particles[i].Intersects(node.Particles[j]) {
				result = append(result, &container.Pair[Particle, Particle]{
					First:  node.Particles[i],
					Second: node.Particles[j],
				})
			}
		}
	}
	if node.Type != Leaf {
		for _, ch := range node.Children {
			for _, p := range node.Particles {
				result = append(result, findIntersectionsInDescendants(ch, p)...)
			}
		}
		for _, ch := range node.Children {
			result = append(result, findAllIntersections(ch)...)
		}
	}
	return result
}

func findIntersectionsInDescendants(node *Node, p Particle) []*container.Pair[Particle, Particle] {
	var result []*container.Pair[Particle, Particle]

	for _, val := range node.Particles {
		if p.Intersects(val) {
			result = append(result, &container.Pair[Particle, Particle]{
				First:  p,
				Second: val,
			})
		}
	}
	if node.Type != Leaf {
		for _, ch := range node.Children {
			result = append(result, findIntersectionsInDescendants(ch, p)...)
		}
	}

	return result
}

func (n *Node) Query(p Particle) []Particle {
	return n.query(p)
}

func (n *Node) query(p Particle) []Particle {
	var result []Particle
	for _, val := range n.Particles {
		if val.Intersects(p) {
			result = append(result, val)
		}
	}
	if n.Type != Leaf {
		for _, ch := range n.Children {
			pos := p.Position()
			side := p.Side()

			inside := pos.X > ch.QuadrantPos.X ||
				pos.X+side < ch.QuadrantPos.X+ch.Width ||
				pos.Y > ch.QuadrantPos.Y ||
				pos.Y+side < ch.QuadrantPos.Y+ch.Width

			if inside {
				result = append(result, ch.query(p)...)
			}
		}
	}

	return result
}

func (n *Node) Clear() {
	if n == nil {
		return
	}
	for _, ch := range n.Children {
		ch.Clear()
	}
	n.Particles = nil
	n.Children = nil
}

func (n *Node) Draw(screen *ebiten.Image) {
	if n == nil {
		return
	}

	ebitenutil.DrawLine(
		screen,
		n.QuadrantPos.X,
		n.QuadrantPos.Y+n.Width/2,
		n.QuadrantPos.X+n.Width,
		n.QuadrantPos.Y+n.Width/2,
		color.White,
	)
	ebitenutil.DrawLine(
		screen,
		n.QuadrantPos.X+n.Width/2,
		n.QuadrantPos.Y,
		n.QuadrantPos.X+n.Width/2,
		n.QuadrantPos.Y+n.Width,
		color.White,
	)
	for _, ch := range n.Children {
		ch.Draw(screen)
	}
}
