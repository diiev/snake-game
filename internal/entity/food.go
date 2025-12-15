package entity

import "math/rand"

type Food struct{ Position Point }

// NewFood создаёт новую еду на случайной позиции

func NewFood() *Food { return &Food{Position: Point{X: 10, Y: 10}} }

// Spawn размещает еду на новой случайной позиции

func (f *Food) Spawn(maxX, maxY int) { f.Position = Point{X: rand.Intn(maxX), Y: rand.Intn(maxY)} }

// IsEaten проверяет, съедена ли еда

func (f *Food) IsEaten(head Point) bool { return f.Position == head }
