package main

import (
	"log"

	"github.com/diiev/snake-game/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(g.Config().ScreenWidth, g.Config().ScreenHeight)
	ebiten.SetWindowTitle("Snake Game")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
