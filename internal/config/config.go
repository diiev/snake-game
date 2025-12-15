package config

import "image/color"

type Config struct {
	ScreenWidth  int
	ScreenHeight int
	TileSize     int

	SnakeColor color.Color
	FoodColor  color.Color
	BgColor    color.Color

	GameSpeed int
}

func Default() Config {
	return Config{
		ScreenWidth:  640,
		ScreenHeight: 480,
		TileSize:     20,
		GameSpeed:    10,
	}
}
