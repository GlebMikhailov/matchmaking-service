package usecases

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"matchmaking-service/src/core/config"
	"matchmaking-service/src/features/root/application/dto"
)

type IDeletePlayerUsecase interface {
	DeletePlayer(ctx context.Context, player dto.DeletePlayerDto) error
}

type deletePlayerUsecase struct {
	redisClient *redis.Client
	ctx         context.Context
}

func (deletePu *deletePlayerUsecase) DeletePlayer(ctx context.Context, player dto.DeletePlayerDto) error {
	playerId := player.Id
	for _, element := range config.AllGameModes {
		deletePu.redisClient.ZRem(ctx, "players:"+string(element), playerId)
	}

	_, err := deletePu.redisClient.HDel(ctx, "players", playerId).Result()

	if err != nil {
		return errors.New("player not found")
	}
	return nil
}

func DeletePlayerUsecase(redisClient *redis.Client) IDeletePlayerUsecase {
	return &deletePlayerUsecase{
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}
