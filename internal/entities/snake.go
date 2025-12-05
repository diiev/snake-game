package entities

import "errors"

type Direction int

const (
	DirectionUp Direction = iota
	DirectionRight
	DirectionDown
	DirectionLeft
)

type Snake struct {
	body      []Vector2
	direction Direction
	gridSize  int
	growing   bool
}

func NewSnake(initialLenght, gridSize int) *Snake {
	body := make([]Vector2, initialLenght)
	startX := 5

	for i := range body {
		body[i] = Vector2{startX - i, 5}

	}
	return &Snake{
		body:      body,
		direction: DirectionRight,
		gridSize:  gridSize,
		growing:   false,
	}
}
func (s *Snake) Move() {
	head := s.Head()
	var newHead Vector2

	switch s.direction {
	case DirectionUp:
		newHead = Vector2{head.X, head.Y - 1}
	case DirectionRight:
		newHead = Vector2{head.X + 1, head.Y}
	case DirectionDown:
		newHead = Vector2{head.X, head.Y + 1}
	case DirectionLeft:
		newHead = Vector2{head.X - 1, head.Y}
	}
	// Добавляем новую голову
	s.body = append([]Vector2{newHead}, s.body...)
	// Если не растем, удаляем хвост
	if !s.growing {
		s.body = s.body[:len(s.body)-1]
	} else {
		s.growing = false
	}

}

func (s *Snake) ChangeDirection(newDir Direction) error {
	// Предотвращаем разворот на 180 градусов
	if (s.direction == DirectionUp && newDir == DirectionDown) ||
		(s.direction == DirectionDown && newDir == DirectionUp) ||
		(s.direction == DirectionLeft && newDir == DirectionRight) ||
		(s.direction == DirectionRight && newDir == DirectionLeft) {
		return errors.New("Invalid direction change")
	}
	s.direction = newDir
	return nil
}

func (s *Snake) Head() Vector2 {
	return s.body[0]
}

func (s *Snake) Grow() {
	s.growing = true
}

func (s *Snake) GetBody() []Vector2 {
	return s.body
}

func (s *Snake) SelfCollision() bool {
	head := s.Head()
	// Начинаем с 1, чтобы не проверять голову с самой собой
	for i := 1; i < len(s.body); i++ {
		if head.Equal(s.body[i]) {
			return true
		}
	}
	return false
}

func (s *Snake) Length() int {
	return len(s.body)
}
