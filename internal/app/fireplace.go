package app

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/snowzhop/verlet-fireplace/internal/math"
	"github.com/snowzhop/verlet-fireplace/internal/physics"
)

type Fireplace struct {
	game

	gravity math.Vec2

	movableObjects []*physics.VerletObject
	// should be static objects
}

func NewFireplace(screenWidth, screenHeight int) *Fireplace {
	radius := float64(10)
	obj := physics.NewVerletObject(
		math.Vec2{
			X: (float64(screenWidth) + radius) / 2,
			Y: (float64(screenHeight) + radius) / 2,
		},
		radius,
	)

	return &Fireplace{
		game: game{
			screenWidth:  screenWidth,
			screenHeight: screenHeight,
		},
		gravity:        math.Vec2{X: 0, Y: 0.1},
		movableObjects: []*physics.VerletObject{obj},
	}
}

func (f *Fireplace) Update() error {
	f.applyGravity()
	f.updatePositions(1)

	return nil
}

func (f *Fireplace) Draw(screen *ebiten.Image) {
	for _, obj := range f.movableObjects {
		ebitenutil.DrawCircle(
			screen,
			obj.CurrentPosition.X,
			obj.CurrentPosition.Y,
			obj.Radius,
			obj.Color,
		)
	}

	if len(f.movableObjects) > 0 {
		ebitenutil.DebugPrint(
			screen,
			fmt.Sprintf(
				"%v, %v",
				f.movableObjects[0].CurrentPosition.X,
				f.movableObjects[0].CurrentPosition.Y,
			),
		)
	}

}

func (f *Fireplace) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return f.screenWidth, f.screenHeight
}
