package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"go.mod/internal/config"
	"go.mod/internal/game"
)

func main() {
	cfg := config.Default()
	g := game.New(cfg)

	ebiten.SetWindowSize(cfg.ScreenWidth, cfg.ScreenHeight)
	ebiten.SetWindowTitle("Snake Game")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
