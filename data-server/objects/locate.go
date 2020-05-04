package objects

import (
	log "github.com/sirupsen/logrus"
	"object-storage-go/common/rabbitmq"
	"object-storage-go/data-server/model"
	"os"
	"strconv"
)

func Locate(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func StartLocate()  {
	mq := rabbitmq.New(rabbitmq.GetRabbitMqDialUrl())
	defer mq.Close()

	mq.Bind("dataServer")
	ch := mq.Consume()

	for msg := range ch{
		//去掉传入的json的引号
		str, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}
		log.Infof("data server receive %s from rmq", str)
		if Locate(model.Config.DataServerConfig.StorageRoot + "/objects/" + str) {
			listenAddress := model.Config.DataServerConfig.DataServerAddress + ":" +
				model.Config.DataServerConfig.DataServerPort
			log.Infof("data server find object, address %s", listenAddress)
			mq.Send(msg.ReplyTo, listenAddress)
		}
	}
}

