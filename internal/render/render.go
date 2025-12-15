package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"go.mod/internal/config"
	"go.mod/internal/entity"
)

type Renderer struct{ config config.Config }

// NewRenderer создаёт новый renderer
func NewRenderer(cfg config.Config) *Renderer {
	// Подстрахуемся дефолтным цветом фона
	if cfg.BgColor == nil {
		cfg.BgColor = color.RGBA{R: 0, G: 0, B: 0, A: 0xff}
	}
	if cfg.SnakeColor == nil {
		cfg.SnakeColor = color.RGBA{R: 0, G: 255, B: 0, A: 0xff}
	}
	if cfg.FoodColor == nil {
		cfg.FoodColor = color.RGBA{R: 255, G: 0, B: 0, A: 0xff}
	}

	return &Renderer{config: cfg}
}

// Draw отрисовывает всё на экран
func (r *Renderer) Draw(screen *ebiten.Image, snake *entity.Snake, food *entity.Food, score int) {
	// очистка экрана
	screen.Fill(r.config.BgColor)

	r.drawSnake(screen, snake)
	r.drawFood(screen, food)
}

func (r *Renderer) drawSnake(screen *ebiten.Image, snake *entity.Snake) {

	for i, segment := range snake.Body {
		x := float32(segment.X * r.config.TileSize)
		y := float32(segment.Y * r.config.TileSize)
		size := float32(r.config.TileSize - 2)
		// Голова светлее
		var clr color.Color = r.config.SnakeColor
		if i == 0 {
			clr = color.RGBA{R: 0, G: 255, B: 200, A: 255}
		}
		vector.DrawFilledRect(screen, x+1, y+1, size, size, clr, true)
	}
}

func (r *Renderer) drawFood(screen *ebiten.Image, food *entity.Food) {
	x := float32(food.Position.X * r.config.TileSize)
	y := float32(food.Position.Y * r.config.TileSize)
	size := float32(r.config.TileSize - 2)
	vector.DrawFilledRect(screen, x+1, y+1, size, size, r.config.FoodColor, true)
}

func (r *Renderer) drawGrid(screen *ebiten.Image) {
	// Вертикальные линии
	for x := 0; x <= r.config.ScreenWidth; x += r.config.TileSize {
		vector.StrokeLine(screen, float32(x), 0, float32(x), float32(r.config.ScreenHeight), 1, r.config.BgColor, true)
	}
	// Горизонтальные линии
	for y := 0; y <= r.config.ScreenHeight; y += r.config.TileSize {
		vector.StrokeLine(screen, 0, float32(y), float32(r.config.ScreenWidth), float32(y), 1, r.config.BgColor, true)
	}
}

func (r *Renderer) drawScore(screen *ebiten.Image, score int) {
	// Простой текст счёта (без использования текстурного шрифта)

}
