package unit

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

func SetupTestRedis(t *testing.T) (*redis.Client, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections").WithOccurrence(1),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err, "Failed to start Redis container")

	port, err := redisC.MappedPort(ctx, "6379")
	assert.NoError(t, err, "Failed to get Redis mapped port")

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("localhost:%s", port.Port()),
		DB:   0,
	})

	pong, err := redisClient.Ping(ctx).Result()
	assert.NoError(t, err, "Failed to ping Redis client")
	assert.Equal(t, "PONG", pong, "Redis did not respond with PONG")

	return redisClient, func() {
		err := redisClient.Close()
		if err != nil {
			t.Fatalf("Failed to close Redis connection: %v", err)
		}
		err = redisC.Terminate(ctx)
		if err != nil {
			t.Fatalf("Failed to terminate Redis container: %v", err)
		}
	}
}
