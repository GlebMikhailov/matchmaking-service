package redis

import (
	"github.com/redis/go-redis/v9"
	"matchmaking-service/src/core/utils"
	"os"
	"strconv"
)

func GetRedis() *redis.Client {
	db, _ := strconv.Atoi(utils.CoalesceStr(os.Getenv("REDIS_DB"), "3"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	return rdb
}
