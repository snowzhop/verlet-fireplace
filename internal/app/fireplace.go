package app

import (
	_ "embed"
	"fmt"
	stdcolor "image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/snowzhop/verlet-fireplace/internal/graphic"
	"github.com/snowzhop/verlet-fireplace/internal/math"
	"github.com/snowzhop/verlet-fireplace/internal/math/geom/v3"
	"github.com/snowzhop/verlet-fireplace/internal/math/quadtree"
	"github.com/snowzhop/verlet-fireplace/internal/physics"
	"github.com/snowzhop/verlet-fireplace/resource"
)

const (
	minRadius  = float64(2)
	maxRadius  = float64(11)
	radius     = float64(3)
	cellLen    = radius*2 + 2
	rootOffset = float64(20)
)

var (
	particleSprite *ebiten.Image
	bloomShader    *ebiten.Shader
)

func init() {
	particleSprite = ebiten.NewImage(int(2*radius), int(2*radius))
	ebitenutil.DrawCircle(
		particleSprite,
		radius,
		radius,
		radius,
		stdcolor.RGBA{
			R: 255, G: 0, B: 0, A: 255,
		},
	)
}

type Fireplace struct {
	game

	gravity        math.Vec2
	root           *quadtree.Node
	movableObjects []*physics.VerletObject

	heatEmitters []*physics.VerletObject

	// graphic
	particleTextureMap *graphic.TextureMap

	// obsolete
	field                *geom.Field
	hiddenObjects        []*physics.VerletObject
	staticMainConstraint physics.Circle
}

func NewFireplace(screenWidth, screenHeight int) *Fireplace {
	fmt.Printf("screenHeight: %d\n", screenHeight)
	fmt.Printf("screenWidth: %d\n", screenWidth)

	var (
		game = game{
			screenWidth:           screenWidth,
			screenHeight:          screenHeight,
			temperatureStep:       10,
			temperatureLosing:     0.8,
			heatEmitterEfficiency: 0.0025,
			bloom:                 true,
		}
		screenWidthF64 = float64(screenWidth)
		// root           = quadtree.New(float64(screenWidth))
		root    = quadtree.NewWithStart(-rootOffset, -rootOffset, screenWidthF64+2*rootOffset)
		objects []*physics.VerletObject

		// spawn particles
		startCount = 2000
		offset     = float64(0)
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

	// --------------- TEST ----------------
	// objects = nil
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
	// -------------------------------------

	var heatEmitters []*physics.VerletObject
	for i := float64(0); i < screenWidthF64; i += radius * 7 {
		heatEmitterRadius := math.NormFloat64() * 10
		if heatEmitterRadius >= radius {
			heatEmitters = append(heatEmitters, physics.NewVerletObjectWithTemp(
				math.Vec2{
					X: i,
					Y: screenWidthF64 - heatEmitterRadius,
				},
				heatEmitterRadius,
				1000,
			))
		}
	}

	var (
		err            error
		rawBloomShader = resource.BloomShader()
	)
	fmt.Printf("%s\n", string(rawBloomShader))
	bloomShader, err = ebiten.NewShader(rawBloomShader)
	if err != nil {
		panic(fmt.Sprintf("failed to load bloomShader: %v", err))
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
		particleTextureMap: graphic.NewParticleTextureMap(
			int(minRadius),
			int(maxRadius),
			0,
			physics.MaxTemperature,
		),
	}
}

func (f *Fireplace) Update() error {
	f.readInputs()

	if !f.game.pause {
		var (
			dt       float64 = 0.2
			subSteps         = 3
			subDt            = dt / float64(subSteps)
		)

		for i := 0; i < subSteps; i++ {
			f.rebuildTree()
			f.applyHeat()
			f.applyForces()
			f.applyAllConstraints()
			f.solveCollisions7()
			f.updatePositions2(subDt)
			f.recalculateRadiuses()
		}
	}

	return nil
}

func (f *Fireplace) Draw(screen *ebiten.Image) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	var (
		tempSum            float64
		hiddenObjectsCount int

		commonImageOptions = &ebiten.DrawImageOptions{}

		w, h = screen.Size()
	)

	ballSource := ebiten.NewImage(w, h)
	for _, obj := range f.movableObjects {
		if obj.Hidden {
			hiddenObjectsCount++
			continue
		}
		if obj.Temperature() == 0 {
			continue
		}

		tempSum += obj.Temperature()

		commonImageOptions.GeoM.Reset()
		commonImageOptions.GeoM.Translate(obj.CurrentPosition.X-obj.Radius(), obj.CurrentPosition.Y-obj.Radius())
		ballSource.DrawImage(
			f.particleTextureMap.TemperatureImage(obj.Radius(), obj.Temperature()),
			commonImageOptions,
		)

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

	if f.game.bloom {
		shaderOptions := &ebiten.DrawRectShaderOptions{}
		shaderOptions.Uniforms = map[string]interface{}{
			"Horizontal": float32(0),
		}
		shaderOptions.Images[0] = ballSource
		w, h := ballSource.Size()
		imgBuffer := ebiten.NewImage(w, h)
		imgBuffer.DrawRectShader(
			f.game.screenWidth,
			f.game.screenHeight,
			bloomShader,
			shaderOptions,
		)

		shaderOptions.Uniforms = map[string]interface{}{
			"Horizontal": float32(1),
		}
		shaderOptions.Images[0] = imgBuffer
		screen.DrawRectShader(
			f.game.screenWidth,
			f.game.screenHeight,
			bloomShader,
			shaderOptions,
		)
	} else {
		commonImageOptions.GeoM.Reset()
		screen.DrawImage(ballSource, commonImageOptions)
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
			"count: %d | hidden: %d | average temp: %f | %f",
			len(f.movableObjects)-hiddenObjectsCount,
			hiddenObjectsCount,
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
