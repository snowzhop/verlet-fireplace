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

	//go:embed bloom.kage
	bloomShader []byte

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

func BloomShader() []byte {
	result := make([]byte, len(bloomShader))
	n := copy(result, bloomShader)
	if n != len(bloomShader) {
		panic("failed to copy bloomShader to result")
	}
	return result
}
