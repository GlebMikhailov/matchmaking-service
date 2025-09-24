package usecases

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type ICheckRedisUsecase interface {
	CheckRedis(ctx context.Context) bool
}

type createCheckRedisUsecase struct {
	redisClient *redis.Client
	ctx         context.Context
}

func (checkNatsUsecase *createCheckRedisUsecase) CheckRedis(ctx context.Context) bool {
	_, err := checkNatsUsecase.redisClient.Ping(ctx).Result()
	return err == nil
}

func CheckRedisUsecase(redisClient *redis.Client) ICheckRedisUsecase {
	return &createCheckRedisUsecase{
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}
