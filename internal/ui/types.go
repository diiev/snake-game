package ui

import (
	"time"
)

// MenuData содержит данные для экрана меню
type MenuData struct {
	MenuItems     []MenuItem
	SelectedIndex int
	GameVersion   string
	HighScore     int
}

// MenuItem - пункт меню
type MenuItem struct {
	Text   string
	Action string // "start", "highscores", "settings", "exit"
}

// GameData содержит данные для игрового экрана
type GameData struct {
	Player    *entities.Player
	Enemies   []*entities.Enemy
	Food      *entities.Food
	Bonuses   []*entities.Bonus
	Obstacles []*entities.Obstacle

	Score          int
	Level          int
	Lives          int
	LevelScore     int
	LevelProgress  float64 // 0.0 - 1.0
	NextLevelScore int

	GameTime time.Duration
	IsPaused bool
}

// GameOverData содержит данные для экрана Game Over
type GameOverData struct {
	FinalScore     int
	LevelReached   int
	LivesRemaining int
	TotalTime      time.Duration

	IsHighScore  bool
	PlayerName   string
	EnteringName bool

	// Для отображения
	CauseOfDeath string // "wall", "enemy", "self"
	LastWords    []string
}

// HighScoresData содержит данные для таблицы рекордов
type HighScoresData struct {
	Scores         []HighScoreEntry
	PlayerPosition int // Позиция игрока в таблице (-1 если нет)
	PlayerScore    int
}

// HighScoreEntry - запись в таблице рекордов
type HighScoreEntry struct {
	Rank  int
	Name  string
	Score int
	Level int
	Time  string
	Date  string
}

// SettingsData содержит данные для настроек
type SettingsData struct {
	MusicVolume   float64
	SoundVolume   float64
	GameSpeed     string // "slow", "normal", "fast"
	ControlScheme string // "arrows", "wasd"
	ShowGrid      bool
	ShowFPS       bool
}

// LevelData содержит данные о прогрессе уровня
type LevelData struct {
	CurrentLevel    int
	LevelName       string
	LevelProgress   float64
	EnemiesCount    int
	SpeedMultiplier float64
	IsBossLevel     bool
}

// HUDData содержит данные для HUD (статистика в реальном времени)
type HUDData struct {
	Score      int
	Lives      int
	Level      int
	Time       string
	FPS        float64
	TPS        float64
	BonusTimer *BonusTimerData
}

// BonusTimerData содержит данные о таймере бонуса
type BonusTimerData struct {
	Type     string
	TimeLeft string
	IsActive bool
	Color    string
}
