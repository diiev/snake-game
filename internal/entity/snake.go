package entity

type Snake struct {
	Body      []Point
	Direction Point
	NextDir   Point
}

// NewSnake создаёт новую змею в центре поля
func NewSnake(x, y int) *Snake {
	return &Snake{Body: []Point{{X: x, Y: y}, {X: x - 1, Y: y}, {X: x - 2, Y: y}}, Direction: DirectionRight, NextDir: DirectionRight}
}

// Move двигает змею на один шаг

func (s *Snake) Move() {
	s.Direction = s.NextDir
	newHead := Point{X: s.Body[0].X + s.Direction.X, Y: s.Body[0].Y + s.Direction.Y}
	s.Body = append([]Point{newHead}, s.Body[:len(s.Body)-1]...)

}

// Grow увеличивает длину змеи
func (s *Snake) Grow() {
	tail := s.Body[len(s.Body)-1]
	s.Body = append(s.Body, tail)
}

// SetDirection устанавливает направление движения
func (s *Snake) SetDirection(d Point) {
	// Не позволяем развернуться на 180 градусов
	if (s.Direction.X+d.X) != 0 || (s.Direction.Y+d.Y) != 0 {
		s.NextDir = d
	}
}

// GetHead возвращает голову змеи
func (s *Snake) GetHead() Point { return s.Body[0] } // Len возвращает длину змеи func (s *Snake) Len() int { return len(s.Body) }

// Len возвращает длину змеи
func (s *Snake) Len() int { return len(s.Body) }
