package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

const (
	ScreenWidth  = 1024
	ScreenHeight = 1024
)

func main() {
	g := NewGame()

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Game of Hive")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
