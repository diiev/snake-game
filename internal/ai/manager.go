package ai

import "math/rand"

type Behaivor int

const (
	BehaivorRandom Behaivor = iota
	BehaivorChase
	BehaivorEvade
	BehaivorPatrol
)

type Manager struct {
	behaviors map[int]Behaivor
}

func NewManager() *Manager {
	return &Manager{
		behaviors: make(map[int]Behaivor),
	}
}

func (m *Manager) GetAction(enemy *entities.Enemy, player *entities.Player, food entities.Food) entities.Direction {
	behaivor := m.getBehavior(enemy.ID())

	switch behaivor {
	case BehaivorChase:
		return m.chasePlayer(enemy, player)
	case BehaivorEvade:
		return m.evadePlayer(enemy, player)
	case BehaivorPatrol:
		return m.patrol(enemy)
	default:
		return m.randomMove(enemy)
	}
}

func (m *Manager) chasePlayer(enemy *entities.Enemy, player *entities.Player) entities.Direction {
	enemyPos := enemy.Head()
	playerPos := player.Head()

	// Простой алгорит преследования
	dx := playerPos.X - enemyPos.X
	dy := playerPos.Y - enemyPos.Y

	// Предпочитаем движение по оси с большей разницей
	if abs(dx) > abs(dy) {
		if dx > 0 && enemy.CanMove(entities.DirectionRight) {
			return entities.DirectionRight
		} else if dx < 0 && enemy.CanMove(entities.DirectionLeft) {
			return entities.DirectionLeft
		}
	}

	if dy > 0 && enemy.CanMove(entities.DirectionDown) {
		return entities.DirectionDown
	} else if dy < 0 && enemy.CanMove(entities.DirectionUp) {
		return entities.DirectionUp
	}

	// Если не можем двигаться к игроку, двигаемся случайно
	return m.randomMove(enemy)
}

func (m *Manager) evadePlayer(enemy *entities.Enemy, player *entities.Player) entities.Direction {
	enemyPos := enemy.Head()
	playerPos := player.Head()

	dx := playerPos.X - enemyPos.X
	dy := playerPos.Y - enemyPos.Y

	if abs(dx) > abs(dy) {
		if dx > 0 && enemy.CanMove(entities.DirectionLeft) {
			return entities.DirectionLeft
		} else if dx < 0 && enemy.CanMove(entities.DirectionRight) {
			return entities.DirectionRight
		}
	}
	if dy > 0 && enemy.CanMove(entities.DirectionUp) {
		return entities.DirectionUp
	} else if dy < 0 && enemy.CanMove(entities.DirectionDown) {
		return entities.DirectionDown
	}

	return m.randomMove(enemy)
}

func (m *Manager) randomMove(enemy *entities.Enemy) entities.Direction {
	// Пробуем случайные направления, пока не найдем валидное

	directions := []entities.Direction{
		entities.DirectionUp,
		entities.DirectionRight,
		entities.DirectionDown,
		entities.DirectionLeft,
	}
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	for _, dir := range directions {
		if enemy.CanMove(dir) {
			return dir
		}
	}
	return enemy.CurrentDirection()
}

func (m *Manager) patrol(enemy *entities.Enemy) entities.Direction {
	// Простой патрулирование - меняем направление через случайные интервалы

	if rand.Intn(10) == 0 {
		return m.randomMove(enemy)
	}
	return enemy.CurrentDirection()
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (m *Manager) getBehavior(id int) Behaivor {
	if behavior, exists := m.behaviors[id]; exists {
		return behavior
	}

	// Назначаем случайное поведение

	behaviors := []Behaivor{BehaivorRandom, BehaivorChase, BehaivorEvade, BehaivorPatrol}
	behavior := behaviors[rand.Intn(len(behaviors))]
	m.behaviors[id] = behavior

	return behavior
}
