package app

import (
	"fmt"
	stdcolor "image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/snowzhop/verlet-fireplace/internal/color"
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
	radius := float64(8)

	fmt.Printf("screenHeight: %d\n", screenHeight)
	fmt.Printf("screenWidth: %d\n", screenWidth)

	var objects []*physics.VerletObject
	for x := int(radius); x < screenWidth-int(radius); x += 2 * int(radius) {
		for y := int(radius); y < int(float32(screenHeight)/2); y += 2 * int(radius) {
			objects = append(objects, physics.NewVerletObjectWithTemp(
				math.Vec2{
					X: float64(x),
					Y: float64(y),
				},
				radius,
				100,
			))
		}

	}

	return &Fireplace{
		game: game{
			screenWidth:  screenWidth,
			screenHeight: screenHeight,
		},
		gravity:        math.Vec2{X: 0, Y: 0.5},
		movableObjects: objects,
	}
}

func (f *Fireplace) Update() error {
	if repeatingMouseClick(ebiten.MouseButtonLeft) {
		cursorX, cursorY := ebiten.CursorPosition()

		obj := physics.NewVerletObject(
			math.Vec2{X: float64(cursorX), Y: float64(cursorY)},
			10,
		)

		f.movableObjects = append(f.movableObjects, obj)
	}
	if repeatingKeyPress(ebiten.KeyC) {
		f.movableObjects = f.movableObjects[:0]
	}

	var (
		dt       float64 = 1
		subSteps         = 2
		subDt            = dt / float64(subSteps)
	)

	for i := 0; i < subSteps; i++ {
		f.applyGravity()
		f.applyConstraint()
		f.solveCollisions()
		f.updatePositions(subDt)
	}

	return nil
}

func (f *Fireplace) Draw(screen *ebiten.Image) {
	// f.withCircleConstraint(screen)

	var tempSum float64
	for _, obj := range f.movableObjects {
		tempSum += obj.Temperature

		ebitenutil.DrawCircle(
			screen,
			obj.CurrentPosition.X,
			obj.CurrentPosition.Y,
			obj.Radius,
			color.ToRGBByTemperature(uint32(obj.Temperature)),
		)
	}

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf(
			"count: %d\naverage temp: %f",
			len(f.movableObjects),
			float32(tempSum)/float32(len(f.movableObjects)),
		))
}

func (f *Fireplace) withCircleConstraint(screen *ebiten.Image) {
	screen.Fill(stdcolor.RGBA{
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
		stdcolor.Black,
	)
	ebitenutil.DrawCircle(
		screen,
		f.staticMainConstraint.Position.X,
		f.staticMainConstraint.Position.Y,
		1,
		stdcolor.RGBA{
			R: 255,
			A: 255,
		},
	)
}

func (f *Fireplace) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return f.screenWidth, f.screenHeight
}
