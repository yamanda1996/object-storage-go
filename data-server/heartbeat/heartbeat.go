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
	log.Debugf("data server start heart beat")

	for {
		mq.Publish("apiServer", listenAddress)
		time.Sleep(time.Duration(5) * time.Second)
	}
}