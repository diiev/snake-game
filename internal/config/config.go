package config

import "time"

type Config struct {
	// Графика
	ScreenWidth  int
	ScreenHeight int
	GridSize     int

	// Игровой процесс

	InitialSpeed   time.Duration
	SpeedIncrement time.Duration
	InitialLives   int
	PointsPerFruit int
	PointsPerBonus int

	// Уровни

	LevelThresholds []int // Очки для перехода на уровень
	EnemiesPerLevel []int //Количество врагов на уровне

	//Внешний вид
	Colors    map[string]string
	FontSizes map[string]int

	// Аудио

	MusicVolume float64
	SoundVolume float64

	// Система
	MaxHighScores int
	SaveFilePath  string
}

func Default() *Config {
	return &Config{
		ScreenWidth:  800,
		ScreenHeight: 600,
		GridSize:     20,

		InitialSpeed:   150 * time.Millisecond,
		SpeedIncrement: 10 * time.Millisecond,
		InitialLives:   3,
		PointsPerFruit: 10,
		PointsPerBonus: 50,

		LevelThresholds: []int{100, 250, 500, 1000, 2000},
		EnemiesPerLevel: []int{0, 1, 1, 2, 2, 3},

		Colors: map[string]string{
			"background": "#0f172a",
			"grid":       "#1e293b",
			"snake":      "#10b981",
			"snake_head": "#34d399",
			"food":       "#ef4444",
			"enemy":      "#8b5cf6",
			"bonus":      "#fbbf24",
			"text":       "#f8fafc",
			"ui":         "#475569",
		},

		FontSizes: map[string]int{
			"small":  16,
			"medium": 24,
			"large":  36,
		},

		MusicVolume:   0.5,
		SoundVolume:   0.7,
		MaxHighScores: 10,
		SaveFilePath:  "assets/data/scores.json",
	}
}
