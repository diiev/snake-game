package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"go.mod/internal/config"
	"go.mod/internal/input"
	"go.mod/internal/render"
	"go.mod/internal/world"
)

type Game struct {
	world        *world.World
	renderer     *render.Renderer
	inputHandler *input.Handler
	config       config.Config
	gameOver     bool
	score        int
	updateCount  int
}

// New создаёт новую игру
func New(config config.Config) *Game {
	gridWidth := config.ScreenWidth / config.TileSize
	gridHeight := config.ScreenHeight / config.TileSize
	w := world.NewWorld(gridWidth, gridHeight)
	return &Game{world: w, renderer: render.NewRenderer(config), inputHandler: input.NewHandler(w.Snake), config: config, gameOver: false, score: 0}
}

// Update обновляет состояние игры
func (g *Game) Update() error {
	if g.gameOver {
		// Нажми R для перезагрузки
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			*g = *New(g.config)
		}
		return nil
	}
	// Обработка ввода
	g.inputHandler.Update()
	// Обновляем мир с нужной скоростью
	g.updateCount++
	if g.updateCount >= 60/g.config.GameSpeed {
		g.world.Update()
		g.score = g.world.Snake.Len() - 3
		// Начальная длина 3
		g.updateCount = 0
		// Проверяем окончание игры
		if g.world.CheckGameOver() {
			g.gameOver = true
		}
	}
	return nil
}

// Draw отрисовывает игру
func (g *Game) Draw(screen *ebiten.Image) {
	g.renderer.Draw(screen, g.world.Snake, g.world.Food, g.score)
	if g.gameOver {
		// TODO: Нарисовать экран Game Over
	}
}

// Layout возвращает логические размеры экрана
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.config.ScreenWidth, g.config.ScreenHeight
}
