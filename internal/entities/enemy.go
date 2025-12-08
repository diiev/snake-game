package entities

import (
	"math/rand"
	"time"
)

type Enemy struct {
	body      []Vector2
	direction Direction
	gridSize  int
	id        int
	color     string
	speed     int // 1 = обычная, 2 = быстрая
	lastMove  time.Time
}

func NewEnemy(gridSize, id int) *Enemy {
	rand.Seed(time.Now().UnixNano() * int64(id+1))

	enemy := &Enemy{
		body:      []Vector2{{10 + id, 10 + id}},
		direction: Direction(rand.Intn(4)),
		gridSize:  gridSize,
		id:        id,
		color:     "#8b5cf6",
		speed:     1,
	}
	return enemy
}

func (e *Enemy) Update(newDirection Direction) {
	// Проверяем таймер для скорости
	if e.speed > 1 {
		if time.Since(e.lastMove) < time.Second/time.Duration(e.speed*2) {
			return
		}
	}
	e.lastMove = time.Now()

	// Меняем направление если нужно
	if newDirection != e.direction {
		e.direction = newDirection
	}
	head := e.Head()
	var newHead Vector2

	switch e.direction {
	case DirectionUp:
		newHead = Vector2{head.X, head.Y - 1}
	case DirectionRight:
		newHead = Vector2{head.X + 1, head.Y}
	case DirectionDown:
		newHead = Vector2{head.X, head.Y + 1}
	case DirectionLeft:
		newHead = Vector2{head.X - 1, head.Y}

	}
	e.body = []Vector2{newHead}

}

func (e *Enemy) CanMove(dir Direction) bool {
	// Простая проверка не выходить за границы
	head := e.Head()

	switch dir {
	case DirectionUp:
		return head.Y > 0
	case DirectionRight:
		return head.X < 39
	case DirectionDown:
		return head.Y < 29
	case DirectionLeft:
		return head.X > 0

	}
	return true
}

func (e *Enemy) Head() Vector2 {
	return e.body[0]

}
func (e *Enemy) Body() []Vector2 {
	return e.body
}

func (e *Enemy) Color() string {
	return e.color
}
func (e *Enemy) ID() int {
	return e.id
}

func (e *Enemy) CurrentDirection() Direction {
	return e.direction
}
