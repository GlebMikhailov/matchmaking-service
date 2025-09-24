package side_car

import (
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"matchmaking-service/src/core/config"
	_ "matchmaking-service/src/docs"
	usecases "matchmaking-service/src/features/side_car/health/business/usecase"
	"strconv"
)
import "matchmaking-service/src/features/side_car/health/communication"

// @title Matchmaking Side app API
// @version 1.0
// @description Side car app's docs for health check, metrics
// @host localhost:8080
// @BasePath /
// @schemes http
func InitSideCar(config config.AppConfig, conn *nats.Conn, redisClient *redis.Client) {
	g := gin.Default()

	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	checkNatsUsecase := usecases.CheckNatsUsecase(conn)
	checkRedisUsecase := usecases.CheckRedisUsecase(redisClient)
	communication.CheckHealthHandler(g, checkNatsUsecase, checkRedisUsecase)

	err := g.Run("0.0.0.0:" + strconv.Itoa(config.SideCarAppConfig.Port))
	if err != nil {
		panic(err)
	}
}
