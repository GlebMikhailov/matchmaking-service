package models

import "matchmaking-service/src/core/config"

type Player struct {
	Id        string            `json:"id"`
	Trophies  int               `json:"trophies"`
	IsPremium bool              `json:"is_premium"`
	JoinedAt  int64             `json:"joined_at"`
	GameModes []config.GameMode `json:"game_mode"`
}
