package locate

import (
	"object-storage-go/data-server/model"
	"object-storage-go/data-server/rabbitmq"
	"os"
	"strconv"
)

func Locate(name string) bool {
	_, err := os.Stat(name)
	return os.IsExist(err)
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

		if Locate(model.Config.DataServerConfig.StorageRoot + "/objects/" + str) {
			listenAddress := model.Config.DataServerConfig.DataServerAddress + ":" +
				string(model.Config.DataServerConfig.DataServerPort)
			mq.Send(msg.ReplyTo, listenAddress)
		}
	}
}

