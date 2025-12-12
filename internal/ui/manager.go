package ui

import (
	"fmt"

	"github.com/diiev/snake-game/internal/config"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type Manager struct {
	config    *config.Config
	fonts     map[string]font.Face
	showDebug bool
}

func NewManager(cfg *config.Config) (*Manager, error) {
	m := &Manager{
		config:    cfg,
		fonts:     make(map[string]font.Face),
		showDebug: false,
	}

	if err := m.loadFonts(); err != nil {
		return nil, err
	}

	return m, nil
}
func (m *Manager) Draw(screen *ebiten.Image, state gamestate.State, data interface{}) {
	// Сначала рисуем фон
	m.drawBackground(screen)

	// Затем рисуем контент в зависимости от состояния
	switch state {
	case gamestate.StateMenu:
		if menuData, ok := data.(MenuData); ok {
			m.drawMenu(screen, menuData)
		}
	case gamestate.StatePlaying:
		if gameData, ok := data.(GameData); ok {
			m.drawGame(screen, gameData)
			// Отдельно рисуем HUD поверх всего
			m.drawHUD(screen, m.GetHUDData())
		}
	case gamestate.StatePaused:
		if gameData, ok := data.(GameData); ok {
			m.drawGame(screen, gameData)
			m.drawPauseOverlay(screen)
		}
	case gamestate.StateGameOver:
		if gameOverData, ok := data.(GameOverData); ok {
			m.drawGameOver(screen, gameOverData)
		}
	case gamestate.StateHighScores:
		if highScoresData, ok := data.(HighScoresData); ok {
			m.drawHighScores(screen, highScoresData)
		}
	case gamestate.StateSettings:
		if settingsData, ok := data.(SettingsData); ok {
			m.drawSettings(screen, settingsData)
		}
	}

	// Отладка
	if m.showDebug {
		m.drawDebugInfo(screen)
	}
}

func (m *Manager) drawGame(screen *ebiten.Image, data GameData) {
	// Рисуем сетку
	if m.config.ShowGrid {
		m.drawGrid(screen)
	}

	// Рисуем игровые объекты
	m.drawPlayer(screen, data.Player)
	m.drawEnemies(screen, data.Enemies)
	m.drawFood(screen, data.Food)
	m.drawBonuses(screen, data.Bonuses)
	m.drawObstacles(screen, data.Obstacles)

	// Рисуем статистику уровня
	m.drawLevelInfo(screen, data)
}

func (m *Manager) drawHUD(screen *ebiten.Image, data HUDData) {
	// Верхняя панель
	m.drawTopPanel(screen, data)

	// Панель бонусов
	if data.BonusTimer != nil && data.BonusTimer.IsActive {
		m.drawBonusTimer(screen, data.BonusTimer)
	}

	// Отображение FPS если включено
	if m.config.ShowFPS {
		m.drawFPS(screen, data.FPS, data.TPS)
	}
}

func (m *Manager) drawTopPanel(screen *ebiten.Image, data HUDData) {
	width := screen.Bounds().Dx()
	panelHeight := 40

	// Фон панели
	m.drawRect(screen, 0, 0, width, panelHeight, m.config.Colors["ui"])

	// Счет
	scoreText := fmt.Sprintf("Счет: %d", data.Score)
	m.drawText(screen, scoreText, 10, 20, m.config.Colors["text"], m.fonts["medium"])

	// Жизни
	livesText := fmt.Sprintf("Жизни: %d", data.Lives)
	m.drawText(screen, livesText, width/4, 20, m.config.Colors["text"], m.fonts["medium"])

	// Уровень
	levelText := fmt.Sprintf("Уровень: %d", data.Level)
	m.drawText(screen, levelText, width/2, 20, m.config.Colors["text"], m.fonts["medium"])

	// Время
	timeText := fmt.Sprintf("Время: %s", data.Time)
	m.drawText(screen, timeText, width*3/4, 20, m.config.Colors["text"], m.fonts["medium"])
}

func (m *Manager) drawLevelInfo(screen *ebiten.Image, data GameData) {
	width := screen.Bounds().Dx()
	y := 50

	// Прогресс-бар уровня
	if data.LevelProgress > 0 {
		m.drawProgressBar(screen,
			width/2-100, y, 200, 20,
			data.LevelProgress,
			fmt.Sprintf("Уровень %d: %d/%d",
				data.Level, data.LevelScore, data.NextLevelScore+data.LevelScore))
	}

	// Информация о врагах
	if len(data.Enemies) > 0 {
		enemyText := fmt.Sprintf("Врагов: %d", len(data.Enemies))
		m.drawText(screen, enemyText, width-100, y+30,
			m.config.Colors["text"], m.fonts["small"])
	}
}

func (m *Manager) drawGameOver(screen *ebiten.Image, data GameOverData) {
	width := screen.Bounds().Dx()
	height := screen.Bounds().Dy()

	// Полупрозрачный оверлей
	m.drawRect(screen, 0, 0, width, height, "rgba(0,0,0,0.7)")

	// Заголовок
	title := "ИГРА ОКОНЧЕНА"
	m.drawText(screen, title, width/2, height/4,
		m.config.Colors["text"], m.fonts["large"], true)

	// Причина смерти
	causeText := fmt.Sprintf("Причина: %s", data.CauseOfDeath)
	m.drawText(screen, causeText, width/2, height/4+50,
		"#ff6b6b", m.fonts["medium"], true)

	// Статистика
	stats := []string{
		fmt.Sprintf("Финальный счет: %d", data.FinalScore),
		fmt.Sprintf("Достигнут уровень: %d", data.LevelReached),
		fmt.Sprintf("Осталось жизней: %d", data.LivesRemaining),
		fmt.Sprintf("Время игры: %s", formatDuration(int(data.TotalTime.Seconds()))),
	}

	y := height/2 - 50
	for _, stat := range stats {
		m.drawText(screen, stat, width/2, y,
			m.config.Colors["text"], m.fonts["medium"], true)
		y += 30
	}

	// Последние слова
	if len(data.LastWords) > 0 {
		y += 20
		for _, word := range data.LastWords {
			m.drawText(screen, word, width/2, y,
				"#94a3b8", m.fonts["small"], true)
			y += 25
		}
	}

	// Если это рекорд
	if data.IsHighScore {
		y += 40
		recordText := "НОВЫЙ РЕКОРД!"
		m.drawText(screen, recordText, width/2, y,
			"#fbbf24", m.fonts["large"], true)

		// Ввод имени
		if data.EnteringName {
			y += 50
			namePrompt := fmt.Sprintf("Введите имя: %s", data.PlayerName)
			m.drawText(screen, namePrompt, width/2, y,
				m.config.Colors["text"], m.fonts["medium"], true)
		}
	}

	// Кнопки
	y = height - 100
	buttons := []string{"Новая игра", "Таблица рекордов", "В меню"}
	for i, button := range buttons {
		m.drawButton(screen, button, width/2, y+(i*60), 200, 40,
			i == 0, m.config.Colors["ui"])
	}
}
