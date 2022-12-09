package geom

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func Fprint(screen *ebiten.Image, field *Field) {

	for i := 1; i < int(field.maxX); i++ {
		_, h := screen.Size()

		x1 := float64(i * int(field.cellSideLen))
		y1 := float64(0)
		x2 := float64(i * int(field.cellSideLen))
		y2 := float64(h)

		ebitenutil.DrawLine(
			screen,
			x1,
			y1,
			x2,
			y2,
			color.White,
		)
	}
	for i := 1; i < int(field.maxY); i++ {
		w, _ := screen.Size()

		x1 := float64(0)
		y := float64(i * int(field.cellSideLen))
		x2 := float64(w)

		ebitenutil.DrawLine(
			screen,
			x1, y,
			x2, y,
			color.White,
		)
	}
}
