package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/snowzhop/verlet-fireplace/internal/math"
	"github.com/snowzhop/verlet-fireplace/internal/physics"
)

func repeatingMouseClick(button ebiten.MouseButton) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.MouseButtonPressDuration(button)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func repeatingKeyPress(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}

	return false
}

func (f *Fireplace) readInputs() {
	switch {
	case repeatingKeyPress(ebiten.KeyF1):
		f.game.debug = !f.game.debug
	case repeatingKeyPress(ebiten.KeyF2):
		f.game.debugTemp = !f.game.debugTemp
	case repeatingKeyPress(ebiten.KeyF3):
		f.game.drawTemp = !f.game.drawTemp
	case repeatingMouseClick(ebiten.MouseButtonLeft):
		cursorX, cursorY := ebiten.CursorPosition()

		offsetX, offsetY := math.RandomOffset()

		obj := physics.NewVerletObject(
			math.Vec2{
				X: float64(cursorX) + offsetX,
				Y: float64(cursorY) + offsetY,
			},
			radius,
		)

		f.root.Insert(obj)

		f.movableObjects = append(f.movableObjects, obj)
	case repeatingKeyPress(ebiten.KeyC):
		f.movableObjects = f.movableObjects[:0]
	case repeatingKeyPress(ebiten.KeyP):
		f.game.pause = !f.game.pause
	case repeatingKeyPress(ebiten.KeyB):
		f.game.bloom = !f.game.bloom
	}
}
