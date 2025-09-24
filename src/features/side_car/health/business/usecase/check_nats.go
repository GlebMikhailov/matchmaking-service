package usecases

import (
	"context"
	"github.com/nats-io/nats.go"
)

type ICheckNatsUsecase interface {
	CheckNats(ctx context.Context) bool
}

type createCheckNatsUsecase struct {
	natsClient *nats.Conn
	ctx        context.Context
}

func (checkNatsUsecase *createCheckNatsUsecase) CheckNats(ctx context.Context) bool {
	return checkNatsUsecase.natsClient.Status() == nats.CONNECTED
}

func CheckNatsUsecase(natsClient *nats.Conn) ICheckNatsUsecase {
	return &createCheckNatsUsecase{
		natsClient: natsClient,
		ctx:        context.Background(),
	}
}
