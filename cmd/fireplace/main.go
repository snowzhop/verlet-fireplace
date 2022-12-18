package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/snowzhop/verlet-fireplace/internal/app"
)

func main() {
	width := 600

	fireplace := app.NewFireplace(width, width)

	ebiten.SetWindowSize(width, width)
	ebiten.SetWindowTitle("Verlet Fireplace")

	if err := ebiten.RunGame(fireplace); err != nil {
		log.Fatalf("error has occured during the game: %v", err)
	}
}
