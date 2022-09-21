package color

import (
	"image/color"

	"github.com/snowzhop/verlet-fireplace/internal/physics"
)

func ToRGBByTemperature(t uint32) color.RGBA {
	if t < physics.MaxTemperature/2 {
		r := float32(t) / physics.MaxTemperature * 2 * 255
		return color.RGBA{
			R: uint8(r),
			A: 255,
		}
	} else {
		gb := (float32(t) - physics.MaxTemperature/2) / physics.MaxTemperature * 2 * 255
		return color.RGBA{
			R: 255,
			G: uint8(gb),
			B: uint8(gb),
			A: 255,
		}
	}
}
