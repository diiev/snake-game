package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"go.mod/internal/entity"
)

type Handler struct{ snake *entity.Snake }

// NewHandler создаёт новый input handler
func NewHandler(snake *entity.Snake) *Handler { return &Handler{snake: snake} }

// Update обрабатывает ввод пользователя
func (h *Handler) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		h.snake.SetDirection(entity.DirectionUp)
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		h.snake.SetDirection(entity.DirectionDown)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		h.snake.SetDirection(entity.DirectionLeft)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		h.snake.SetDirection(entity.DirectionRight)
	}
}
