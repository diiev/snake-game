package entity

// Point представляет координату на игровом поле
type Point struct{ X, Y int }

// // Direction константы для направлений движения

var (
	DirectionUp    = Point{X: 0, Y: -1}
	DirectionDown  = Point{X: 0, Y: 1}
	DirectionLeft  = Point{X: -1, Y: 0}
	DirectionRight = Point{X: 1, Y: 0}
)
