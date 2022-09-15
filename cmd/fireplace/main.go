package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/snowzhop/verlet-fireplace/internal/app"
)

func main() {
	fireplace := app.NewFireplace(640, 480)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Verlet Fireplace")

	if err := ebiten.RunGame(fireplace); err != nil {
		log.Fatalf("error has occured during the game: %v", err)
	}
}
