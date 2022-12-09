package app

import (
	_ "embed"
	"fmt"
	stdcolor "image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/snowzhop/verlet-fireplace/internal/color"
	"github.com/snowzhop/verlet-fireplace/internal/math"
	"github.com/snowzhop/verlet-fireplace/internal/math/geom/v3"
	"github.com/snowzhop/verlet-fireplace/internal/math/quadtree"
	"github.com/snowzhop/verlet-fireplace/internal/physics"
	"github.com/snowzhop/verlet-fireplace/resource"
)

const (
	radius     = float64(6)
	cellLen    = radius*2 + 2
	rootOffset = float64(20)
)

var (
	particleSprite ebiten.Image
)

func init() {
	particleSprite = resource.Particle()
}

type Fireplace struct {
	game
	drawOptions ebiten.DrawImageOptions

	gravity math.Vec2

	field          *geom.Field
	root           *quadtree.Node
	movableObjects []*physics.VerletObject
	hiddenObjects  []*physics.VerletObject

	heatEmitters []*physics.VerletObject

	staticMainConstraint physics.Circle
}

func NewFireplace(screenWidth, screenHeight int) *Fireplace {
	fmt.Printf("screenHeight: %d\n", screenHeight)
	fmt.Printf("screenWidth: %d\n", screenWidth)

	var (
		game = game{
			screenWidth:           screenWidth,
			screenHeight:          screenHeight,
			temperatureStep:       1.05,
			temperatureLosing:     3,
			heatEmitterEfficiency: 0.8,
		}
		screenWidthF64 = float64(screenWidth)
		// root = quadtree.New(float64(screenWidth))
		root    = quadtree.NewWithStart(-rootOffset, -rootOffset, screenWidthF64+2*rootOffset)
		objects []*physics.VerletObject

		// spawn particles
		startCount = 400
		offset     = float64(0)

		// spawn heat emitters
		mean                 = screenWidthF64 / 2
		stdDeviation         = float64(70)
		normalDistrib        = math.NewGauss(mean, stdDeviation)
		heatRadiusMultiplier = float64(600)
	)

	for y := screenWidthF64 - radius; true; y -= 2 * radius {
		if offset > 0 {
			offset = 0
		} else {
			offset = radius
		}
		for x := radius + offset; x < screenWidthF64; x += 2 * radius {
			obj := physics.NewVerletObjectWithTemp(
				math.Vec2{
					X: x,
					Y: y,
				},
				radius,
				1,
			)

			objects = append(objects, obj)
			root.Insert(obj)
		}
		if len(objects) >= startCount {
			break
		}
	}

	// for i := 0; i < 2; i++ {
	// 	obj := physics.NewVerletObjectWithTemp(
	// 		math.Vec2{
	// 			X: float64(i)*radius + 1,
	// 			Y: 0,
	// 		},
	// 		radius,
	// 		100,
	// 	)

	// 	objects = append(objects, obj)
	// 	root.Insert(obj)
	// }

	var heatEmitters []*physics.VerletObject
	for i := float64(0); i < screenWidthF64; i += radius * 7 {
		heatEmitterRadius := normalDistrib.Calculate(i) * heatRadiusMultiplier * radius * 2
		if heatEmitterRadius >= radius {
			heatEmitters = append(heatEmitters, physics.NewVerletObjectWithTemp(
				math.Vec2{
					X: i,
					Y: screenWidthF64 - heatEmitterRadius,
				},
				heatEmitterRadius,
				800,
			))
		}
	}

	fmt.Printf("movableObjects count: %d\n", len(objects))
	fmt.Printf("heat emitters count: %d\n", len(heatEmitters))
	root.Dump()

	return &Fireplace{
		game:           game,
		gravity:        math.Vec2{X: 0, Y: 0.5},
		movableObjects: objects,
		root:           root,
		heatEmitters:   heatEmitters,
	}
}

func (f *Fireplace) Update() error {
	switch {
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
	case repeatingKeyPress(ebiten.KeyF1):
		f.game.debug = !f.game.debug
	case repeatingKeyPress(ebiten.KeyF2):
		f.game.debugTemp = !f.game.debugTemp
	case repeatingKeyPress(ebiten.KeyF3):
		f.game.drawTemp = !f.game.drawTemp
	}

	var (
		dt       float64 = 1
		subSteps         = 2
		subDt            = dt / float64(subSteps)
		// subDt = float64(1)
	)

	for i := 0; i < subSteps; i++ {
		f.rebuildTree()
		f.applyHeat()
		f.applyForces()
		f.applyAllConstraints()
		f.solveCollisions6()
		f.updatePositions2(subDt)
	}

	// time.Sleep(100 * time.Millisecond)

	return nil
}

func (f *Fireplace) Draw(screen *ebiten.Image) {
	// f.withCircleConstraint(screen)

	var tempSum float64
	for _, obj := range f.movableObjects {
		tempSum += obj.Temperature()

		if f.game.drawTemp {
			ebitenutil.DrawCircle(
				screen,
				obj.CurrentPosition.X,
				obj.CurrentPosition.Y,
				obj.Radius(),
				color.ToRGBByTemperature(uint32(obj.Temperature())),
			)
		} else {
			f.drawOptions.GeoM.Reset()
			f.drawOptions.GeoM.Translate(obj.CurrentPosition.X-obj.Radius(), obj.CurrentPosition.Y-obj.Radius())
			screen.DrawImage(&particleSprite, &f.drawOptions)
		}

		if f.game.debug {
			ebitenutil.DebugPrintAt(
				screen,
				fmt.Sprintf("%d", obj.ID()),
				int(obj.CurrentPosition.X)-3,
				int(obj.CurrentPosition.Y)-10,
			)
		}
		if f.game.debugTemp {
			ebitenutil.DebugPrintAt(
				screen,
				fmt.Sprintf("%f", obj.Temperature()),
				int(obj.CurrentPosition.X)-3,
				int(obj.CurrentPosition.Y)-10,
			)
		}
	}
	if f.game.drawTemp {
		for _, emitter := range f.heatEmitters {
			ebitenutil.DrawCircle(
				screen,
				emitter.CurrentPosition.X,
				emitter.CurrentPosition.Y,
				emitter.Radius(),
				stdcolor.RGBA{
					R: 0, G: 0, B: 255, A: 255,
				},
			)
		}
	}

	if f.debug {
		f.root.Draw(screen)
	}

	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf(
			"count: %d | average temp: %f | %f",
			len(f.movableObjects),
			float32(tempSum)/float32(len(f.movableObjects)),
			ebiten.ActualFPS(),
		),
	)
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
