package heartbeat

import (
	"object-storage-go/data-serverage-go/data-server/rabbitmq"
	"os"
	"time"
)

func HeartBeat()  {
	mq := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer mq.Close()

	for {
		mq.Publish("apiServer", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(time.Duration(5) * time.Second)
	}
}