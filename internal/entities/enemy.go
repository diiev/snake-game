package entities

import "time"

type Enemy struct {
	body      []Vector2
	direction Direction
	gridSize  int
	id        int
	color     string
	speed     int // 1 = обычная, 2 = быстрая
	lastMove  time.Time
}
