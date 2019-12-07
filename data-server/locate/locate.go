package locate

import (
	"object-storage-go/data-serverage-go/data-server/rabbitmq"
	"os"
	"strconv"
)

func Locate(name string) bool {
	_, err := os.Stat(name)
	return os.IsExist(err)
}

func StartLocate()  {
	mq := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer mq.Close()

	mq.Bind("dataServer")
	ch := mq.Consume()

	for msg := range ch{
		//去掉传入的json的引号
		str, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}

		if Locate(os.Getenv("STORAGE_ROOT") + "/objects/" + str) {
			mq.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}

