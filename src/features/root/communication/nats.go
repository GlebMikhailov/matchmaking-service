package communication

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"matchmaking-service/src/features/root/application/dto"
	"matchmaking-service/src/features/root/business/usecases"
	"time"
)

func CreatePlayer(playerUsecase usecases.IAddPlayerUsecase) nats.MsgHandler {
	return func(m *nats.Msg) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		var message dto.CreatePlayerDto
		err := json.Unmarshal(m.Data, &message)
		if err != nil {
			log.Printf("Seialization error: %v", err)
			return
		}
		err = playerUsecase.AddPlayer(ctx, message)
		if err != nil {
			return
		}
	}
}

func DeletePlayer(playerUsecase usecases.IDeletePlayerUsecase) nats.MsgHandler {
	return func(m *nats.Msg) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		var message dto.DeletePlayerDto
		err := json.Unmarshal(m.Data, &message)
		if err != nil {
			log.Printf("Seialization error: %v", err)
			return
		}
		err = playerUsecase.DeletePlayer(ctx, message)
		if err != nil {
			log.Printf("Cannot delete player: %v", err)
		}
	}
}
