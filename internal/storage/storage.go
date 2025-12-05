package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type HighScore struct {
	Name     string    `json:"name"`
	Score    int       `json:"score"`
	Level    int       `json:"level"`
	Date     time.Time `json:"date"`
	Duration int       `json:"duration"` // в секундах
}

type Storage struct {
	filePath   string
	highScores []HighScore
	mu         sync.RWMutex
	maxEntries int
}

func NewStorage(filepath string) (*Storage, error) {
	s := &Storage{
		filePath:   filepath,
		highScores: make([]HighScore, 0),
		maxEntries: 10,
	}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

func (s *Storage) AddScore(name string, score, level, duration int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	newScore := HighScore{
		Name:     name,
		Score:    score,
		Level:    level,
		Date:     time.Now(),
		Duration: duration,
	}

	// Добавляем и сортируем
	s.highScores = append(s.highScores, newScore)

	// Обрезаем до максимального количества
	if len(s.highScores) > s.maxEntries {
		s.highScores = s.highScores[:s.maxEntries]
	}

	// Сохраняем
	if err := s.save(); err != nil {
		return false
	}

	return true
}

func (s *Storage) sortScores() {
	sort.Slice(s.highScores, func(i, j int) bool {
		// Сначала по очкам, потом по уровню, потом по времени
		if s.highScores[i].Score != s.highScores[j].Score {
			return s.highScores[i].Score > s.highScores[j].Score
		}
		if s.highScores[i].Level != s.highScores[j].Level {
			return s.highScores[i].Level > s.highScores[j].Level
		}
		return s.highScores[i].Duration < s.highScores[j].Duration
	})
}

func (s *Storage) GetHighScores() []HighScore {
	s.mu.RLock()
	defer s.mu.RUnlock()

	scores := make([]HighScore, len(s.highScores))
	copy(scores, s.highScores)
	return scores
}

func (s *Storage) IsHighScore(score int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.highScores) < s.maxEntries {
		return true
	}
	return score > s.highScores[len(s.highScores)-1].Score
}

func (s *Storage) load() error {
	file, err := os.Open(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&s.highScores)
}

func (s *Storage) save() error {
	// Создаем директорию если не существует
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	return encoder.Encode(s.highScores)
}

// Утилиты для форматирования

func (s *Storage) FormatHighScores() []string {
	scores := s.GetHighScores()
	formatted := make([]string, len(scores))

	for i, score := range scores {
		formatted[i] = fmt.Sprintf("%d. %s - %d pts (Level %d, %s)",
			i+1,
			score.Name,
			score.Score,
			score.Level,
			formatDuration(score.Duration),
		)
	}
	return formatted
}

func formatDuration(seconds int) string {
	minutes := seconds / 60
	secs := seconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, secs)
}
