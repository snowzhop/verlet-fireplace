package quadtree

import (
	"fmt"
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
	Unknown = -1
)

type Node struct {
	QuadrantPos math.Vec2
	TotalPos    math.Vec2
	Count       uint32
	Width       float64

	Type     NodeType
	Children []*Node
	Particle Particle
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
	return Unknown
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
		n.splitLeaf(p)
		return
	}

	n.TotalPos = math.SumVec2(n.TotalPos, p.Position())
	n.Count++
	n.Children[n.defineQuadrant(p.Position())].Insert(p)
}

func (n *Node) splitLeaf(newParticle Particle) {
	if n.Particle != nil {
		a := n.Particle
		b := newParticle

		n.TotalPos = math.SumVec2(n.TotalPos, b.Position())
		n.Count++

		current := n
		qA := n.defineQuadrant(a.Position())
		qB := n.defineQuadrant(b.Position())

		loopCounter := 0
		for qA == qB {
			loopCounter++
			if loopCounter > 200 {
				fmt.Println("!!!    loopCounter > 200    !!!")
				fmt.Println("a:", a)
				fmt.Println("b:", b)
			}
			current.split()
			current = current.Children[qA]
			qA = current.defineQuadrant(a.Position())
			qB = current.defineQuadrant(b.Position())

			current.TotalPos = math.SumVec2(current.TotalPos, a.Position(), b.Position())
			current.Count += 2
		}

		current.split()

		current.Children[qA].Particle = a
		current.Children[qA].TotalPos = math.SumVec2(
			current.Children[qA].TotalPos,
			a.Position(),
		)
		current.Children[qA].Count++

		current.Children[qB].Particle = b
		current.Children[qB].TotalPos = math.SumVec2(
			current.Children[qB].TotalPos,
			b.Position(),
		)
		current.Children[qB].Count++

		n.Particle = nil
		return
	}

	n.Particle = newParticle
	n.TotalPos = math.SumVec2(newParticle.Position())
	n.Count++
}

func (n *Node) Dump() {
	// TODO: print tree info
}

func (n *Node) FindIntersections() []*container.Pair[Particle, Particle] {
	return findAllIntersections(n)
}

func findAllIntersections(node *Node) []*container.Pair[Particle, Particle] {
	var result []*container.Pair[Particle, Particle]
	if node.Type != Leaf {
		for _, ch := range node.Children {
			if node.Particle != nil {
				result = append(result, findIntersectionsInDescendants(ch, node.Particle)...)
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
	switch node.Type {
	case Leaf:
		if node.Particle != nil {
			if node.Particle.Intersects(p) {
				result = append(result, &container.Pair[Particle, Particle]{
					First:  p,
					Second: node.Particle,
				})
			}
		}
	case Internal:
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
	switch n.Type {
	case Leaf:
		if n.Particle == nil || n.Particle == p {
			break
		}
		if n.Particle.Intersects(p) {
			result = append(result, n.Particle)
		}
	case Internal:
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
	n.Particle = nil
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
