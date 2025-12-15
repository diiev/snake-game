package world

import "go.mod/internal/entity"

type World struct {
	Snake  *entity.Snake
	Food   *entity.Food
	Width  int
	Height int
}

// NewWorld создаёт новый мир с змеёй и едой

func NewWorld(width, height int) *World {
	snake := entity.NewSnake(width/2, height/2)
	food := entity.NewFood()
	food.Spawn(width, height)
	return &World{Snake: snake, Food: food, Width: width, Height: height}
}

// Update обновляет состояние мира

func (w *World) Update() {
	w.Snake.Move()
	// Проверяем, съедена ли еда
	if w.Food.IsEaten(w.Snake.GetHead()) {
		w.Snake.Grow()
		w.Food.Spawn(w.Width, w.Height)
	}
}

// CheckGameOver проверяет условия окончания игры
func (w *World) CheckGameOver() bool {
	head := w.Snake.GetHead()
	// Столкновение со стенами
	if head.X < 0 || head.X >= w.Width || head.Y < 0 || head.Y >= w.Height {
		return true
	}
	// Столкновение с собой
	for i := 1; i < w.Snake.Len(); i++ {
		if head == w.Snake.Body[i] {
			return true
		}
	}
	return false
}
