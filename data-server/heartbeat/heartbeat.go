package heartbeat

import (
	log "github.com/sirupsen/logrus"
	"object-storage-go/data-server/model"
	"object-storage-go/data-server/rabbitmq"
	"strconv"
	"strings"
	"time"
)

func HeartBeat()  {
	mq := rabbitmq.New(rabbitmq.GetRabbitMqDialUrl())
	defer mq.Close()
	var builder strings.Builder
	builder.WriteString(model.Config.DataServerConfig.DataServerAddress + ":")
	builder.WriteString(strconv.Itoa(model.Config.DataServerConfig.DataServerPort))
	listenAddress := builder.String()

	for {
		mq.Publish("apiServer", listenAddress)
		log.Debugf("data server [%s] send heart beat", model.Config.DataServerConfig.DataServerAddress)
		time.Sleep(time.Duration(5) * time.Second)
	}
}