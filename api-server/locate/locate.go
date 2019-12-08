package locate

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"object-storage-go/api-server/rabbitmq"
	"strconv"
	"strings"
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

func Handler(w http.ResponseWriter, r *http.Request)  {

	m := r.Method

	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	info := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])

	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b, _ := json.Marshal(info)

	w.Write(b)
}

func Locate(name string) string {
	mq := rabbitmq.New(rabbitmq.GetRabbitMqDialUrl())
	mq.Publish("dataServer", name)


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
