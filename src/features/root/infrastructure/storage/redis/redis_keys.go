package redis

import "matchmaking-service/src/core/config"

func GetMatchPlayersRedisKey(mode config.GameMode) string {
	return "players:" + string(mode)
}
