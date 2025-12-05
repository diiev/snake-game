package audio

import (
	"embed"
	"log"

	"github.com/diiev/snake-game/internal/config"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
	"github.com/hajimehoshi/ebiten/audio/wav"
)

var audioFS embed.FS

type Manager struct {
	audioContext *audio.Context
	musicPlayers map[string]*audio.Player
	soundPlayers map[string]*audio.Player
	config       *config.Config
}

func NewManager(cfg *config.Config) (*Manager, error) {
	context := audio.NewContext(44100)
	m := &Manager{
		audioContext: context,
		musicPlayers: make(map[string]*audio.Player),
		soundPlayers: make(map[string]*audio.Player),
		config:       cfg,
	}

	// Загрузка музыки
	if err := m.loadMusic("main_theme", "assets/audio/music/main_theme.mp3"); err != nil {
		return nil, err
	}
	if err := m.loadMusic("menu", "assets/audio/music/menu.mp3"); err != nil {
		return nil, err
	}
	sounds := []string{"eat", "bonus", "level_up", "game_over", "hit"}
	for _, sound := range sounds {
		if err := m.loadMusic(sound, "assets/audio/sfx/"+sound+".waw"); err != nil {
			log.Printf("Warning: failed to load sound %s: %v", sound, err)
		}
	}
	return m, nil
}

func (m *Manager) loadMusic(name, path string) error {
	file, err := audioFS.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	stream, err := mp3.Decode(m.audioContext, file)
	if err != nil {
		return err
	}
	player, err := m.audioContext.NewPlayer(stream)
	if err != nil {
		return err
	}
	player.SetVolume(m.config.MusicVolume)
	m.musicPlayers[name] = player

	return nil
}

func (m *Manager) loadSound(name, path string) error {
	file, err := audioFS.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	stream, err := wav.Decode(m.audioContext, file)
	if err != nil {
		return err
	}
	player, err := m.audioContext.NewPlayer(stream)
	if err != nil {
		return err
	}
	player.SetVolume(m.config.SoundVolume)
	m.soundPlayers[name] = player
	return nil
}

func (m *Manager) PlayMusic(name string, loop bool) {
	if player, exists := m.musicPlayers[name]; exists {
		if loop {
			player.Rewind()
		}
		player.Play()
	}
}

func (m *Manager) PlaySound(name string) {
	if player, exists := m.soundPlayers[name]; exists {
		player.Rewind()
		player.Play()
	}

}

func (m *Manager) StopMusic(name string) {
	if player, exists := m.musicPlayers[name]; exists {
		player.Pause()
	}
}

func (m *Manager) SetMusicVolume(volume float64) {
	for _, player := range m.musicPlayers {
		player.SetVolume(volume)
	}
}

func (m *Manager) SetSoundVolume(volume float64) {
	for _, player := range m.soundPlayers {
		player.SetVolume(volume)
	}
}
