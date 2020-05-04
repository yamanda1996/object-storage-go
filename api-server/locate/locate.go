package locate

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"object-storage-go/common/rabbitmq"
	"strconv"
	"time"
)

func LocateFile(context *gin.Context)  {
	name := context.Param("filename")
	l := Locate(name)
	if l == "" {
		context.String(http.StatusNotFound, "can not find file %s", name)
	} else {
		context.String(http.StatusOK, "find %s success in %s", name, l)
	}
}

func Locate(name string) string {
	mq := rabbitmq.New(rabbitmq.GetRabbitMqDialUrl())
	mq.Publish("dataServer", name)
	log.Infof("api server put %s to rmq", name)
	select {
	case msg := <- mq.Consume():
		s, _ := strconv.Unquote(string(msg.Body))
		return s
	case <- time.After(time.Duration(1) * time.Second):
		log.Debugf("data server can not find file [%s]", name)
		return ""
	}
}

func Exist(name string) bool {
	return Locate(name) != ""
}
