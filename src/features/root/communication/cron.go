package communication

import (
	"context"
	"github.com/robfig/cron/v3"
	"log"
	"matchmaking-service/src/features/root/business/usecases"
)

func StartCron(ctx context.Context, matchPlayersUsecase usecases.IMatchPlayersUsecase) {
	c := cron.New()
	_, err := c.AddFunc("@every 5s", func() {
		matchPlayersUsecase.MatchPlayers(ctx)
	})
	if err != nil {
		log.Fatalf("Failed to add cron: %v", err)
	}

	c.Start()
	log.Println("Cron-manager is ready")
}
