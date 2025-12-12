package main

import (
	"image"
	"log"

	"github.com/diiev/snake-game/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Инициализация игры
	g, err := game.NewGame()
	if err != nil {
		log.Fatal("Ошибка инициализации игры:", err)
	}

	// Настройка окна
	ebiten.SetWindowSize(g.Config().ScreenWidth, g.Config().ScreenHeight)
	ebiten.SetWindowTitle("Snake Evolution")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowIcon([]image.Image{loadIcon()})
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(60)

	// Запуск игры
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal("Ошибка выполнения игры:", err)
	}
}

func loadIcon() image.Image {
	// Создаем простую иконку
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	drawSnakeIcon(img)
	return img
}
