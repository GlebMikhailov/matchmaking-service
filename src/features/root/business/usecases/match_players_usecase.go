package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"matchmaking-service/src/core/config"
	"matchmaking-service/src/features/root/application/models"
	redis_key "matchmaking-service/src/features/root/infrastructure/storage/redis"
	"sort"
	"strconv"
	"time"
)

type IMatchPlayersUsecase interface {
	MatchPlayers(ctx context.Context)
}

type matchPlayersUsecase struct {
	redisClient *redis.Client
	ctx         context.Context
	appConfig   config.AppConfig
}

func (matchPlayersUsecase *matchPlayersUsecase) MatchPlayers(ctx context.Context) {
	for _, gameMode := range config.AllGameModes {
		var gameModeConfig config.GameModeConfig

		for _, v := range matchPlayersUsecase.appConfig.GameModes {
			if v.GameMode == gameMode {
				gameModeConfig = v
				break
			}
		}

		zsetKey := "players:" + string(gameMode)
		currentTime := time.Now().Unix()

		playerIdsWithScores, err := matchPlayersUsecase.redisClient.ZRangeByScore(ctx, zsetKey, &redis.ZRangeBy{
			Min: "-inf",
			Max: strconv.FormatInt(time.Now().Unix()-gameModeConfig.SkipTime, 10),
		}).Result()

		if err != nil {
			return
		}

		var playerIDs []string
		playerIDs = append(playerIDs, playerIdsWithScores...)

		if len(playerIDs) == 0 {
			return
		}

		playersData, err := matchPlayersUsecase.redisClient.HMGet(ctx, "players", playerIDs...).Result()
		if err != nil {
			return
		}

		var allPlayers []models.Player
		for _, data := range playersData {
			if data == nil {
				continue
			}
			var player models.Player
			if err := json.Unmarshal([]byte(data.(string)), &player); err != nil {
				fmt.Printf("cannot parse player: %v\n", err)
				continue
			}
			allPlayers = append(allPlayers, player)
		}

		var waitingLong []models.Player
		var waitingShort []models.Player

		if gameModeConfig.GameModeBot.IsEnabled {
			for _, player := range allPlayers {
				if time.Now().Unix()-currentTime-player.JoinedAt > gameModeConfig.GameModeBot.SecondsForBots {
					waitingLong = append(waitingLong, player)
				} else {
					waitingShort = append(waitingShort, player)
				}
			}
			processLongWaitingPlayers(ctx, matchPlayersUsecase.redisClient, waitingLong, gameModeConfig, matchPlayersUsecase.appConfig)

			processShortWaitingPlayers(ctx, matchPlayersUsecase.redisClient, waitingShort, gameModeConfig, matchPlayersUsecase.appConfig)
		} else {
			processLongWaitingPlayers(ctx, matchPlayersUsecase.redisClient, waitingLong, gameModeConfig, matchPlayersUsecase.appConfig)

			processShortWaitingPlayers(ctx, matchPlayersUsecase.redisClient, waitingShort, gameModeConfig, matchPlayersUsecase.appConfig)
		}
	}
}

func getMaxAllowedSpread(player *models.Player, gameModeConfig config.GameModeConfig) int {
	timeInQueue := time.Now().Unix() - player.JoinedAt
	maxSpread := gameModeConfig.MaxSpread[0].Spread

	sort.Slice(gameModeConfig.MaxSpread, func(i, j int) bool {
		return gameModeConfig.MaxSpread[i].SecondsInSearch < gameModeConfig.MaxSpread[j].SecondsInSearch
	})

	for _, spreadConfig := range gameModeConfig.MaxSpread {
		if int(timeInQueue) >= spreadConfig.SecondsInSearch {
			maxSpread = spreadConfig.Spread
		} else {
			break
		}
	}
	return maxSpread
}

func filterAndGroupPlayers(players []models.Player, targetGroupSize int, maxSpread int) ([]models.Player, []models.Player) {
	if len(players) < targetGroupSize {
		return nil, players
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].Trophies < players[j].Trophies
	})

	for i := 0; i <= len(players)-targetGroupSize; i++ {
		group := players[i : i+targetGroupSize]
		minTrophies := group[0].Trophies
		maxTrophies := group[targetGroupSize-1].Trophies

		if maxTrophies-minTrophies <= maxSpread {
			return group, append(players[:i], players[i+targetGroupSize:]...)
		}
	}
	return nil, players
}

func processLongWaitingPlayers(ctx context.Context, rdb *redis.Client, players []models.Player, gameModeConfig config.GameModeConfig, appConfig config.AppConfig) {
	currentPlayers := make([]models.Player, len(players))
	copy(currentPlayers, players)

	for len(currentPlayers) > 0 {
		var effectiveMaxSpread int
		if len(currentPlayers) > 0 {
			effectiveMaxSpread = getMaxAllowedSpread(&currentPlayers[0], gameModeConfig)
		} else {
			effectiveMaxSpread = gameModeConfig.MaxSpread[0].Spread
		}

		if len(currentPlayers) >= gameModeConfig.Players {
			matchGroup, remaining := filterAndGroupPlayers(currentPlayers, gameModeConfig.Players, effectiveMaxSpread)
			if matchGroup != nil {
				createMatch(matchGroup, 0, gameModeConfig.Players)
				removePlayersFromQueue(ctx, rdb, appConfig, matchGroup)
				currentPlayers = remaining
				continue
			}
		}

		if len(currentPlayers) > 0 {
			botsNeeded := gameModeConfig.Players - len(currentPlayers)
			createMatch(currentPlayers, botsNeeded, gameModeConfig.Players)
			removePlayersFromQueue(ctx, rdb, appConfig, currentPlayers)
			currentPlayers = nil
		}
		break
	}
}

func processShortWaitingPlayers(ctx context.Context, rdb *redis.Client, players []models.Player, gameModeConfig config.GameModeConfig, appConfig config.AppConfig) {
	currentPlayers := make([]models.Player, len(players))
	copy(currentPlayers, players)
	var effectiveMaxSpread int
	if len(gameModeConfig.MaxSpread) > 0 {
		effectiveMaxSpread = gameModeConfig.MaxSpread[0].Spread
	} else {
		effectiveMaxSpread = 100
	}

	for len(currentPlayers) >= gameModeConfig.Players {
		matchGroup, remaining := filterAndGroupPlayers(currentPlayers, gameModeConfig.Players, effectiveMaxSpread)
		if matchGroup != nil {
			createMatch(matchGroup, 0, gameModeConfig.Players)
			removePlayersFromQueue(ctx, rdb, appConfig, matchGroup)
			currentPlayers = remaining
		} else {
			break
		}
	}
}

func removePlayersFromQueue(ctx context.Context, rdb *redis.Client, appConfig config.AppConfig, players []models.Player) {
	for _, element := range config.AllGameModes {
		var playerIDs []interface{}
		for _, player := range players {
			rdb.HDel(ctx, "players", player.Id)
			playerIDs = append(playerIDs, player.Id)
		}

		if len(playerIDs) > 0 {
			rdb.ZRem(ctx, redis_key.GetMatchPlayersRedisKey(element), playerIDs...)
		}
	}

}

func createMatch(players []models.Player, botsNeeded int, targetCount int) {
	// TODO: handle match creation
}

func CreateMatchUsecase(redisClient *redis.Client, appConfig config.AppConfig) IMatchPlayersUsecase {
	return &matchPlayersUsecase{
		redisClient: redisClient,
		ctx:         context.Background(),
		appConfig:   appConfig,
	}
}
