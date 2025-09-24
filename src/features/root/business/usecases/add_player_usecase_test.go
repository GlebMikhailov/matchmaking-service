package usecases

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"matchmaking-service/src/core/config"
	"matchmaking-service/src/features/root/application/dto"
	"matchmaking-service/src/features/root/application/models"
	"matchmaking-service/src/test/unit"
	"testing"
	"time"
)

func TestAddPlayerUsecase(t *testing.T) {
	redisClient, cleanup := unit.SetupTestRedis(t)
	defer cleanup()

	usecase := AddPlayerUsecase(redisClient)
	ctx := context.Background()

	t.Run("Successfully create a new player", func(t *testing.T) {
		testPlayer := dto.CreatePlayerDto{
			Id:        uuid.New().String(),
			Trophies:  1500,
			IsPremium: true,
			GameModes: []config.GameMode{config.Big, config.Small},
		}

		err := usecase.AddPlayer(ctx, testPlayer)
		assert.NoError(t, err)

		playerJSON, err := redisClient.HGet(ctx, "players", testPlayer.Id).Result()
		assert.NoError(t, err)

		var storedPlayer models.Player
		err = json.Unmarshal([]byte(playerJSON), &storedPlayer)
		assert.NoError(t, err)

		assert.Equal(t, testPlayer.Id, storedPlayer.Id)
		assert.Equal(t, testPlayer.Trophies, storedPlayer.Trophies)
		assert.Equal(t, testPlayer.IsPremium, storedPlayer.IsPremium)
		assert.ElementsMatch(t, testPlayer.GameModes, storedPlayer.GameModes)
		assert.WithinDuration(t, time.Unix(storedPlayer.JoinedAt, 0), time.Now(), 5*time.Second)
	})

	t.Run("Throw an error because player already exist", func(t *testing.T) {
		testPlayer := dto.CreatePlayerDto{
			Id:        uuid.New().String(),
			Trophies:  1500,
			IsPremium: true,
			GameModes: []config.GameMode{config.Big, config.Small},
		}

		err := usecase.AddPlayer(ctx, testPlayer)
		assert.NoError(t, err)

		err = usecase.AddPlayer(ctx, testPlayer)
		assert.Error(t, err)
	})
}
