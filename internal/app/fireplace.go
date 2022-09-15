package app

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/snowzhop/verlet-fireplace/internal/math"
	"github.com/snowzhop/verlet-fireplace/internal/physics"
)

type Fireplace struct {
	game

	gravity math.Vec2

	movableObjects []*physics.VerletObject

	staticMainConstraint physics.Circle
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

	circle := physics.Circle{
		Radius: 180,
		Position: math.Vec2{
			X: 250,
			Y: 200,
		},
	}

	return &Fireplace{
		game: game{
			screenWidth:  screenWidth,
			screenHeight: screenHeight,
		},
		gravity:              math.Vec2{X: 0, Y: 0.5},
		movableObjects:       []*physics.VerletObject{obj},
		staticMainConstraint: circle,
	}
}

func (f *Fireplace) Update() error {
	f.applyGravity()
	f.applyConstraint()
	f.updatePositions(1)

	return nil
}

func (f *Fireplace) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{
		R: 128,
		G: 128,
		B: 128,
		A: 255,
	})

	ebitenutil.DrawCircle(
		screen,
		f.staticMainConstraint.Position.X,
		f.staticMainConstraint.Position.Y,
		f.staticMainConstraint.Radius,
		color.Black,
	)
	ebitenutil.DrawCircle(
		screen,
		f.staticMainConstraint.Position.X,
		f.staticMainConstraint.Position.Y,
		1,
		color.RGBA{
			R: 255,
			A: 255,
		},
	)

	for _, obj := range f.movableObjects {
		ebitenutil.DrawCircle(
			screen,
			obj.CurrentPosition.X,
			obj.CurrentPosition.Y,
			obj.Radius,
			obj.Color,
		)
	}
}

func (f *Fireplace) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return f.screenWidth, f.screenHeight
}
