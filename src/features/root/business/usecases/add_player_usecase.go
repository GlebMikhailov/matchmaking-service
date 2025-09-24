package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"matchmaking-service/src/features/root/application/dto"
	"matchmaking-service/src/features/root/application/models"
	rediskey "matchmaking-service/src/features/root/infrastructure/storage/redis"
	"time"
)

type IAddPlayerUsecase interface {
	AddPlayer(ctx context.Context, player dto.CreatePlayerDto) error
}

type addPlayerUsecase struct {
	redisClient *redis.Client
	ctx         context.Context
}

func (pu *addPlayerUsecase) AddPlayer(ctx context.Context, player dto.CreatePlayerDto) error {
	now := time.Now().Unix()
	playerModel := models.Player{
		Id:        player.Id,
		JoinedAt:  now,
		Trophies:  player.Trophies,
		IsPremium: player.IsPremium,
		GameModes: player.GameModes,
	}
	playerJSON, err := json.Marshal(playerModel)

	if err != nil {
		return err
	}

	existingPlayer, _ := pu.redisClient.HGet(ctx, "players", player.Id).Result()

	if existingPlayer != "" {
		return errors.New("player already exists")
	}

	err = pu.redisClient.HSet(ctx, "players", player.Id, playerJSON).Err()
	for _, element := range player.GameModes {

		_, err = pu.redisClient.ZAdd(pu.ctx, rediskey.GetMatchPlayersRedisKey(element), redis.Z{
			Score:  float64(now),
			Member: player.Id,
		}).Result()
	}

	if err != nil {
		return err
	}

	return nil
}

func AddPlayerUsecase(redisClient *redis.Client) IAddPlayerUsecase {
	return &addPlayerUsecase{
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}
