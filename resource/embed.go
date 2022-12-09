package resource

import (
	"bytes"
	_ "embed"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed particle.png
	rawParticle []byte
	particleImg ebiten.Image
)

func init() {
	var (
		err error
		img image.Image
	)
	img, _, err = image.Decode(bytes.NewReader(rawParticle))
	if err != nil {
		log.Fatalf("failed to decode particle img: %v", err)
	}
	particleImg = *ebiten.NewImageFromImage(img)
}

func Particle() ebiten.Image {
	return particleImg
}
