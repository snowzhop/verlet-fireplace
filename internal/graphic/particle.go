package graphic

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/snowzhop/verlet-fireplace/internal/color"
)

type (
	TextureMap struct {
		textures     [][]*ebiten.Image
		radiusOffset int
		tempOffset   int
	}
)

func NewParticleTextureMap(minRadius, maxRadius, minTemp, maxTemp int) *TextureMap {
	textures := make([][]*ebiten.Image, maxRadius-minRadius+1)
	for r := 0; r <= maxRadius-minRadius; r++ {
		textRow := make([]*ebiten.Image, maxTemp-minTemp+1)
		for t := 0; t <= maxTemp-minTemp; t++ {
			radius := r + minRadius
			textRow[t] = ebiten.NewImage(radius*2, radius*2)
			ebitenutil.DrawCircle(
				textRow[t],
				float64(radius),
				float64(radius),
				float64(radius),
				color.ToRGBByTemperature(uint32(t)),
			)
		}
		textures[r] = textRow
	}

	return &TextureMap{
		textures:     textures,
		radiusOffset: minRadius,
		tempOffset:   minTemp,
	}
}

func (m *TextureMap) TemperatureImage(r, t float64) *ebiten.Image {
	return m.textures[int(r)-m.radiusOffset][int(t)-m.tempOffset]
}
