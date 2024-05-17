package alien

import (
	"encoding/json"
	"image/color"
	"log"
	"os"
)

type Config struct {
	ScreenWidth     int        `json:"screenWidth"`
	ScreenHeight    int        `json:"screenHeight"`
	Title           string     `json:"title"`
	BgColor         color.RGBA `json:"bgColor"`
	ShipSpeedFactor float64    `json:"shipSpeedFactor"`

	BulletWidth       int        `json:"bulletWidth"`
	BulletHeight      int        `json:"bulletHeight"`
	BulletSpeedFactor float64    `json:"bulletSpeedFactor"`
	BulletColor       color.RGBA `json:"bulletColor"`
	MaxBulletNum      int        `json:"maxBulletNum"`
	BulletInterval    int64      `json:"bulletInterval"`

	AlienSpeedFactor float64 `json:"alienSpeedFactor"`

	TitleFontSize int `json:"titleFontSize"`
	FontSize      int `json:"fontSize"`
	SmallFontSize int `json:"smallFontSize"`

	ScoreColor color.RGBA `json:"scoreColor"`
}

func LoadConfig() *Config {
	file, err := os.Open("./alien/config.json")
	if err != nil {
		log.Fatalf("os.Open failed: %v\n", err)
	}

	var cfg Config
	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		log.Fatalf("json.Decode failed: %v\n", err)
	}

	return &cfg
}
