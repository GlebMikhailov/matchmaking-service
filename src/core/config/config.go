package config

import (
	"github.com/joho/godotenv"
	"log"
	"matchmaking-service/src/core/utils"
	"os"
	"path/filepath"
	"strconv"
)

type GameMode string

const (
	Big   GameMode = "big"
	Small GameMode = "small"
)

var AllGameModes = []GameMode{
	Big,
	Small,
}

type RatingSpread struct {
	Spread          int
	SecondsInSearch int
}

type GameModeBotConfiguration struct {
	IsEnabled      bool
	SecondsForBots int64
}

type GameModeConfig struct {
	GameMode    GameMode
	Players     int
	IsEnabled   bool
	MaxSpread   []RatingSpread
	GameModeBot GameModeBotConfiguration
	SkipTime    int64
}

type SideCarAppConfig struct {
	Port int
}

type AppConfig struct {
	GameModes        []GameModeConfig
	SideCarAppConfig SideCarAppConfig
}

func GetAppConfig() AppConfig {
	currentDir, _ := os.Getwd()
	envPath := filepath.Join(currentDir, ".env")
	if _, err := os.Stat(envPath); err == nil {
		err := godotenv.Load(envPath)
		if err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
	}

	sideCarAppPort, err := strconv.Atoi(utils.CoalesceStr(os.Getenv("SIDE_CAR_APP_PORT"), "8080"))
	if err != nil {
		log.Fatal("Error parsing SIDE_CAR_APP_PORT")
	}

	return AppConfig{
		GameModes: []GameModeConfig{
			{
				GameMode:  Big,
				Players:   8,
				IsEnabled: true,
				MaxSpread: []RatingSpread{
					{Spread: 100, SecondsInSearch: 0},
					{Spread: 200, SecondsInSearch: 30},
					{Spread: 300, SecondsInSearch: 40},
				},
				GameModeBot: GameModeBotConfiguration{
					IsEnabled:      true,
					SecondsForBots: 60,
				},
				//SkipTime: 20,
				SkipTime: 0,
			},
			{
				GameMode:  Small,
				Players:   4,
				IsEnabled: true,
				MaxSpread: []RatingSpread{
					{Spread: 100, SecondsInSearch: 0},
					{Spread: 200, SecondsInSearch: 40},
				},
				GameModeBot: GameModeBotConfiguration{
					IsEnabled:      true,
					SecondsForBots: 60,
				},
				//SkipTime: 20,
				SkipTime: 0,
			},
		},
		SideCarAppConfig: SideCarAppConfig{
			Port: sideCarAppPort,
		},
	}
}
