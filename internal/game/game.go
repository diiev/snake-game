package game

import (
	"time"

	"github.com/diiev/snake-game/internal/ai"
	"github.com/diiev/snake-game/internal/audio"
	"github.com/diiev/snake-game/internal/config"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/mod/sumdb/storage"
)

type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StatePaused
	StateLevelComplete
	StateGameOver
	StateHighScores
)

type Game struct {
	config *config.Config
	state  GameState

	// Игровые объекты

	player    *entities.player
	enemies   []*entities.enemy
	food      *entities.Food
	bonuses   []*entities.Bonus
	obstacles []*entities.obstacles

	// Системы
	aiManager    *ai.Manager
	audioManager *audio.Manager
	uiManager    *ui.Manager
	storage      *storage.Storage

	// Состояние игры

	score      int
	level      int
	lives      int
	levelScore int
	lastUpdate time.Time
	gameSpeed  time.Duration
	enemyCount int

	// Таймеры
	bonusTimer *time.Timer
	spawnTimer *time.Timer

	// Ввод

	inputBuffer []ebiten.Key
}

func NewGame() (*Game, error) {
	cfg := config.Default()

	g := &Game{
		config:     cfg,
		state:      StateMenu,
		score:      0,
		level:      1,
		lives:      cfg.InitialLives,
		levelScore: 0,
		gameSpeed:  cfg.InitialSpeed,
		enemyCount: cfg.EnemiesPerLevel[0],
	}

	// Инициализация систем
	var err error
	if g.audioManager, err = audio.NewManager(cfg); err != nil {
		return nil, err
	}
	if g.uiManager, err = ui.NewManager(cfg); err != nil {
		return nil, err
	}
	if g.storage, err = storage.NewStorage(cfg.SaveFilePath); err != nil {
		return nil, err
	}

	// Инициализация AI

	g.aiManager = ai.NewManager()

	// Запуск музыки
	g.audioManager.PlayMusic("main_theme", true)
	return g, nil
}

func (g *Game) Update() error {

	// Обновление ввода
	g.handleInput()

	switch g.state {
	case StateMenu:
		return g.updateMenu()
	case StatePlaying:
		return g.updatePlaying()
	case StatePaused:
		return g.updatePaused()
	case StateGameOver:
		return g.updateGameOver()

	}
	return nil
}

func (g *Game) updatePlaying() error {
	// Фиксированный шаг времени
	now := time.Now()

	if now.Sub(g.lastUpdate) < g.gameSpeed {
		return nil
	}
	g.lastUpdate = now

	// Обновление игрока
	g.player.Update(g.inputBuffer)

	// Обновление врагов через AI

	for _, enemy := range g.enemies {
		enemy.Update(g.aiManager.GetAction(enemy, g.player, g.food))

	}

	// Проверка столкновений
	if g.checkCollisions() {
		g.handleCollision()
		return nil
	}
	// Проверка съедания еды
	if g.player.Head().Equals(g.food.Position) {
		g.handleFoodCollected()
	}

	// Проверка бонусов
	for i, bonus := range g.bonuses {
		if g.player.Head().Equals(bonus.Position) {
			g.handleBonusCollected(i)
		}
	}

	// Проверка уровня
	if g.checkLevelComplete() {
		g.handleLevelComplete()
	}
	// Cпавн бонусов
	g.handleBonusSpawn()

	return nil
}

func (g *Game) handleFoodCollected() {
	g.player.Grow()
	g.score += g.config.PointsPerFruit
	g.levelScore += g.config.PointsPerFruit

	// Звук
	g.audioManager.PlaySound("eat")

	// Новая еда
	g.food.RandomizePosition(g.getAllOccupiedPositions())
}

func (g *Game) handleBonusCollected(index int) {
	// Добавляем жизнь
	g.lives++

	// Очки за бонус

	g.score += g.config.PointsPerBonus
	g.levelScore += g.config.PointsPerBonus

	// Звук
	g.audioManager.PlaySound("bonus")

	// Удаляем бонус
	g.bonuses = append(g.bonuses[:index], g.bonuses[index+1]...)
}

func (g *Game) handleLevelComplete() {
	g.state = StateLevelComplete

	// Увеличаем уровень
	g.level++

	// Увеличиваем сложность
	g.gameSpeed -= g.config.SpeedIncrement

	if g.gameSpeed < 50*time.Millisecond {
		g.gameSpeed = 50 * time.Millisecond
	}

	// Добавляем врагов
	if g.level-1 < len(g.config.EnemiesPerLevel) {
		g.enemyCount = g.config.EnemiesPerLevel[g.level-1]

	}
	// Сброс очков уровня
	g.levelScore = 0

	// Звук

	g.audioManager.PlaySound("level_up")

	go func() {
		time.Sleep(2 * time.Second)
		g.resetLevel()
		g.state = StatePlaying
	}()

}

func (g *Game) resetLevel() {
	// Создаем нового игрока
	g.player = entities.NewPlayer(g.config.GridSize)

	// Создаем врагов
	g.enemies = make([]*entities.Enemy, g.enemyCount)
	for i := 0; i < g.enemyCount; i++ {
		g.enemies[i] = entities.NewEnemy(g.config.GridSize, i+1)
	}

	// Созадем еду
	g.food = entities.NewFood(g.config.GridSize)
	g.food.RandomizePosition(g.getAllOccupiedPositions())
	g.bonuses = nil
}
