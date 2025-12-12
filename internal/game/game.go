package game

import (
	"fmt"
	"time"

	"github.com/diiev/snake-game/internal/ai"
	"github.com/diiev/snake-game/internal/audio"
	"github.com/diiev/snake-game/internal/config"
	"github.com/diiev/snake-game/internal/ui"
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

// GetUIData возвращает данные для UI в зависимости от состояния
func (g *Game) GetUIData() interface{} {
	switch g.state {
	case StateMenu:
		return g.getMenuData()
	case StatePlaying, StatePaused:
		return g.getGameData()
	case StateGameOver:
		return g.getGameOverData()
	case StateHighScores:
		return g.getHighScoresData()
	case StateSettings:
		return g.getSettingsData()
	default:
		return nil
	}
}

func (g *Game) getMenuData() ui.MenuData {
	highScores := g.storage.GetHighScores()
	var highScore int
	if len(highScores) > 0 {
		highScore = highScores[0].Score
	}

	return ui.MenuData{
		MenuItems: []ui.MenuItem{
			{Text: "Новая игра", Action: "start"},
			{Text: "Продолжить", Action: "continue"},
			{Text: "Таблица рекордов", Action: "highscores"},
			{Text: "Настройки", Action: "settings"},
			{Text: "Выход", Action: "exit"},
		},
		SelectedIndex: g.menuSelection,
		GameVersion:   "Snake v2.0",
		HighScore:     highScore,
	}
}

func (g *Game) getGameData() ui.GameData {
	// Рассчитываем прогресс уровня
	var progress float64
	var nextLevelScore int

	if g.level-1 < len(g.config.LevelThresholds) && g.config.LevelThresholds[g.level-1] > 0 {
		threshold := g.config.LevelThresholds[g.level-1]
		progress = float64(g.levelScore) / float64(threshold)
		nextLevelScore = threshold - g.levelScore
	}

	return ui.GameData{
		Player:    g.player,
		Enemies:   g.enemies,
		Food:      g.food,
		Bonuses:   g.bonuses,
		Obstacles: g.obstacles,

		Score:          g.score,
		Level:          g.level,
		Lives:          g.lives,
		LevelScore:     g.levelScore,
		LevelProgress:  progress,
		NextLevelScore: nextLevelScore,

		GameTime: time.Since(g.gameStartTime),
		IsPaused: g.state == StatePaused,
	}
}

func (g *Game) getGameOverData() ui.GameOverData {
	// Определяем причину смерти
	cause := "unknown"
	if g.collisionType == "wall" {
		cause = "Стена"
	} else if g.collisionType == "enemy" {
		cause = "Враг"
	} else if g.collisionType == "self" {
		cause = "Сам себя"
	}

	// Генерируем "последние слова" в зависимости от счета
	lastWords := []string{"Игра окончена!"}
	if g.score < 100 {
		lastWords = append(lastWords, "Могло быть и хуже...")
	} else if g.score < 500 {
		lastWords = append(lastWords, "Неплохой результат!")
	} else {
		lastWords = append(lastWords, "Отличная игра!")
	}

	// Проверяем, является ли счет рекордным
	isHighScore := g.storage.IsHighScore(g.score)

	return ui.GameOverData{
		FinalScore:     g.score,
		LevelReached:   g.level,
		LivesRemaining: g.lives,
		TotalTime:      time.Since(g.gameStartTime),

		IsHighScore:  isHighScore,
		PlayerName:   g.playerName,
		EnteringName: g.enteringName,

		CauseOfDeath: cause,
		LastWords:    lastWords,
	}
}

func (g *Game) getHighScoresData() ui.HighScoresData {
	scores := g.storage.GetHighScores()
	formattedScores := make([]ui.HighScoreEntry, len(scores))

	// Определяем позицию текущего игрока
	playerPosition := -1
	for i, score := range scores {
		if score.Name == g.playerName && score.Score == g.score {
			playerPosition = i
		}

		formattedScores[i] = ui.HighScoreEntry{
			Rank:  i + 1,
			Name:  score.Name,
			Score: score.Score,
			Level: score.Level,
			Time:  formatDuration(score.Duration),
			Date:  score.Date.Format("02.01.2006"),
		}
	}

	return ui.HighScoresData{
		Scores:         formattedScores,
		PlayerPosition: playerPosition,
		PlayerScore:    g.score,
	}
}

func (g *Game) getSettingsData() ui.SettingsData {
	return ui.SettingsData{
		MusicVolume:   g.audioManager.GetMusicVolume(),
		SoundVolume:   g.audioManager.GetSoundVolume(),
		GameSpeed:     g.getSpeedString(),
		ControlScheme: g.controlScheme,
		ShowGrid:      g.showGrid,
		ShowFPS:       g.showFPS,
	}
}

// GetHUDData возвращает данные для HUD (можно вызывать каждый кадр)
func (g *Game) GetHUDData() ui.HUDData {
	// Получаем активный бонус
	var bonusTimer *ui.BonusTimerData
	if len(g.bonuses) > 0 && g.bonuses[0].Active {
		bonus := g.bonuses[0]
		timeLeft := bonus.Lifetime - time.Since(bonus.CreatedAt)

		bonusTimer = &ui.BonusTimerData{
			Type:     bonus.Type.String(),
			TimeLeft: formatDuration(int(timeLeft.Seconds())),
			IsActive: true,
			Color:    bonus.GetColor(),
		}
	}

	return ui.HUDData{
		Score:      g.score,
		Lives:      g.lives,
		Level:      g.level,
		Time:       formatDuration(int(time.Since(g.gameStartTime).Seconds())),
		FPS:        ebiten.ActualFPS(),
		TPS:        ebiten.ActualTPS(),
		BonusTimer: bonusTimer,
	}
}

// Вспомогательные методы
func (g *Game) getSpeedString() string {
	baseSpeed := g.config.InitialSpeed
	currentSpeed := g.gameSpeed

	if currentSpeed <= baseSpeed/2 {
		return "fast"
	} else if currentSpeed <= baseSpeed {
		return "normal"
	}
	return "slow"
}

func formatDuration(seconds int) string {
	minutes := seconds / 60
	secs := seconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, secs)
}
