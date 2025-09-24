package main

import (
	"context"
	"log"
	"matchmaking-service/src/core/config"
	appNats "matchmaking-service/src/core/nats"
	"matchmaking-service/src/core/redis"
	"matchmaking-service/src/features/root/business/usecases"
	"matchmaking-service/src/features/root/communication"
	"matchmaking-service/src/features/side_car"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	appConfig := config.GetAppConfig()
	redisDb := redis.GetRedis()
	nc := appNats.GetNats()

	createPlayerUsecase := usecases.AddPlayerUsecase(redisDb)
	deletePlayerUsecase := usecases.DeletePlayerUsecase(redisDb)
	matchPlayersUsecase := usecases.CreateMatchUsecase(redisDb, appConfig)

	_, err := nc.Subscribe("create.player", communication.CreatePlayer(createPlayerUsecase))
	if err != nil {
		log.Fatal(err)
	}
	_, err = nc.Subscribe("delete.player", communication.DeletePlayer(deletePlayerUsecase))
	if err != nil {
		log.Fatal(err)
	}

	communication.StartCron(ctx, matchPlayersUsecase)
	side_car.InitSideCar(appConfig, nc, redisDb)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Finish...")
}
