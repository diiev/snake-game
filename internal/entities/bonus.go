package entities

import (
	"math/rand"
	"time"
)

type BonusType int

const (
	BonusLife BonusType = iota
	BonusSpeed
	BonusShield
	BonusPoints
)

type Bonus struct {
	Position  Vector2
	Type      BonusType
	CreatedAt time.Time
	LifeTime  time.Duration
	Activate  bool
	Value     int
}

func NewBonus(bonusType BonusType, gridsize int) *Bonus {
	lifetime := 10 * time.Second
	var value int
	switch bonusType {
	case BonusLife:
		value = 1
	case BonusPoints:
		value = 100
	default:
		value = 0
	}
	return &Bonus{
		Type:     bonusType,
		LifeTime: lifetime,
		Activate: false,
		Value:    value,
	}

}

func (b *Bonus) Spawn(occupied []Vector2, width, height int) {
	maxX := width / 20
	maxY := height / 20

	for attempts := 0; attempts < 100; attempts++ {
		pos := Vector2{rand.Intn(maxX), rand.Intn(maxY)}

		// 	Проверяем что позиция свободна

		isOccupied := false
		for _, o := range occupied {
			if pos.Equals(o) {
				isOccupied = true
				break
			}
		}
		if !isOccupied {
			b.Position = pos
			b.CreatedAt = time.Now()
			b.Activate = true
			break
		}
	}
}

func (b *Bonus) Update() {
	if b.Activate && time.Since(b.CreatedAt) > b.LifeTime {
		b.Activate = false
	}
}

func (b *Bonus) GetColor() string {
	switch b.Type {
	case BonusLife:
		return "#fbbf24" // Желтый
	case BonusSpeed:
		return "#60a5fa" // Синий
	case BonusShield:
		return "#a78bfa" // Фиолетовый
	case BonusPoints:
		return "#10b981" // Зеленый
	default:
		return "#ffffff"
	}
}
