package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pieti/gameofhive/goh"
)

func main() {
	g := goh.NewGame()

	ebiten.SetWindowSize(goh.ScreenWidth, goh.ScreenHeight)
	ebiten.SetWindowTitle("Game of Hive")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
