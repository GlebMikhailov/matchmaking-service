package nats

import (
	"github.com/nats-io/nats.go"
	"log"
	"matchmaking-service/src/core/utils"
	"os"
)

func GetNats() *nats.Conn {
	natsUrl := utils.CoalesceStr(os.Getenv("NATS_URL"), nats.DefaultURL)
	log.Printf("NATS URL: %s", natsUrl)
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		log.Fatal(err)
	}
	//defer nc.Close()

	return nc
}
