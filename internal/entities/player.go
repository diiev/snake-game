package entities

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	body       []Vector2
	direction  Direction
	gridSize   int
	color      string
	headColor  string
	name       string
	isShielded bool
	shieldTime time.Time
	score      int
	lives      int
}

func NewPlayer(gridSize int) *Player {
	return &Player{
		body: []Vector2{
			{5, 5},
			{4, 5},
			{3, 5},
		},
		direction:  DirectionRight,
		gridSize:   gridSize,
		color:      "#10b981", // Зеленый
		headColor:  "#34d399", // Светло-зеленый
		name:       "Player",
		isShielded: false,
		lives:      3,
	}
}

func (p *Player) Update(keys []ebiten.Key) {
	// Обработка ввода
	for _, key := range keys {
		switch key {
		case ebiten.KeyArrowUp, ebiten.KeyW:
			p.ChangeDirection(DirectionUp)
		case ebiten.KeyArrowRight, ebiten.KeyD:
			p.ChangeDirection(DirectionRight)
		case ebiten.KeyArrowDown, ebiten.KeyS:
			p.ChangeDirection(DirectionDown)
		case ebiten.KeyArrowLeft, ebiten.KeyA:
			p.ChangeDirection(DirectionLeft)
		}
	}

	// Проверка щита
	if p.isShielded && time.Since(p.shieldTime) > 10*time.Second {
		p.isShielded = false
	}
}

func (p *Player) ChangeDirection(newDir Direction) {
	// Предотвращение разворота на 180 градусов
	if (p.direction == DirectionUp && newDir == DirectionDown) ||
		(p.direction == DirectionDown && newDir == DirectionUp) ||
		(p.direction == DirectionLeft && newDir == DirectionRight) ||
		(p.direction == DirectionRight && newDir == DirectionLeft) {
		return
	}
	p.direction = newDir
}

func (p *Player) Move() {
	head := p.Head()
	var newHead Vector2

	switch p.direction {
	case DirectionUp:
		newHead = Vector2{head.X, head.Y - 1}
	case DirectionRight:
		newHead = Vector2{head.X + 1, head.Y}
	case DirectionDown:
		newHead = Vector2{head.X, head.Y + 1}
	case DirectionLeft:
		newHead = Vector2{head.X - 1, head.Y}
	}

	p.body = append([]Vector2{newHead}, p.body...)

	// Если не нужно расти, удаляем хвост
	if !p.growing {
		p.body = p.body[:len(p.body)-1]
	} else {
		p.growing = false
	}
}

func (p *Player) Grow() {
	p.growing = true
}

func (p *Player) AddLife() {
	p.lives++
}

func (p *Player) LoseLife() bool {
	p.lives--
	return p.lives > 0
}

func (p *Player) ActivateShield() {
	p.isShielded = true
	p.shieldTime = time.Now()
}

func (p *Player) IsShielded() bool {
	return p.isShielded
}

func (p *Player) Head() Vector2 {
	return p.body[0]
}

func (p *Player) Body() []Vector2 {
	return p.body
}

func (p *Player) GetLives() int {
	return p.lives
}

func (p *Player) GetScore() int {
	return p.score
}

func (p *Player) AddScore(points int) {
	p.score += points
}

func (p *Player) SetName(name string) {
	p.name = name
}

func (p *Player) GetName() string {
	return p.name
}
