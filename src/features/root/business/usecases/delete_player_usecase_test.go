package usecases

import (
	"context"
	"github.com/google/uuid"
	"matchmaking-service/src/core/config"
	"matchmaking-service/src/features/root/application/dto"
	"matchmaking-service/src/test/unit"
	"testing"
)

func TestDeletePlayerUsecase(t *testing.T) {
	redisClient, cleanup := unit.SetupTestRedis(t)
	defer cleanup()

	addPlayerUsecase := AddPlayerUsecase(redisClient)
	deletePlayerUsecase := DeletePlayerUsecase(redisClient)
	ctx := context.Background()

	t.Run("Successfully delete an existing player", func(t *testing.T) {
		testPlayer := dto.CreatePlayerDto{
			Id:        uuid.New().String(),
			Trophies:  1500,
			IsPremium: true,
			GameModes: []config.GameMode{config.Big, config.Small},
		}

		err := addPlayerUsecase.AddPlayer(ctx, testPlayer)
		if err != nil {
			t.Fatalf("Error adding player %v", err)
			return
		}
		err = deletePlayerUsecase.DeletePlayer(ctx, dto.DeletePlayerDto{Id: testPlayer.Id})
		if err != nil {
			t.Fatalf("Error deleting player %v", err)
		}
	})
}
