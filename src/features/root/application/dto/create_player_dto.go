package dto

import (
	"matchmaking-service/src/core/config"
)

type CreatePlayerDto struct {
	Id        string            `json:"id"`
	Trophies  int               `json:"trophies"`
	IsPremium bool              `json:"is_premium"`
	GameModes []config.GameMode `json:"game_mode"`
}
