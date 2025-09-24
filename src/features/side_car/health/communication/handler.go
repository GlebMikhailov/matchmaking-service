package communication

import (
	"github.com/gin-gonic/gin"
	"matchmaking-service/src/features/side_car/health/business/usecase"
	"net/http"
)

// swagger:model ServiceHealthyResponse
type ServiceHealthyResponse struct {
	ServiceName string `json:"serviceName" example:"redis"`
	Status      string `json:"status" example:"up/down"`
}

// CheckHealthHandler godoc
// @Summary Check external services statuses
// @Tags Health
// @Produce json
// @Success 200 {array} ServiceHealthyResponse "Service healthy"
// @Failure 500 {array} ServiceHealthyResponse "Service non-healthy"
// @Router /health/check [get]
func CheckHealthHandler(router *gin.Engine, checkNatsUsecase usecases.ICheckNatsUsecase, checkRedisUsecase usecases.ICheckRedisUsecase) {
	router.GET("/health/check", func(c *gin.Context) {
		var serviceStatuses []ServiceHealthyResponse
		natsStatus := checkNatsUsecase.CheckNats(c.Request.Context())
		redisStatus := checkRedisUsecase.CheckRedis(c.Request.Context())
		serviceStatuses = append(serviceStatuses, transformToResponse("nats", natsStatus))
		serviceStatuses = append(serviceStatuses, transformToResponse("redis", redisStatus))

		if natsStatus && redisStatus {
			c.JSON(http.StatusOK, serviceStatuses)
		} else {
			c.JSON(http.StatusInternalServerError, serviceStatuses)
		}
	})
}

func transformToResponse(serviceName string, healthy bool) ServiceHealthyResponse {
	status := "up"
	if !healthy {
		status = "down"
	}
	return ServiceHealthyResponse{
		ServiceName: serviceName,
		Status:      status,
	}
}
