package locate

import (
	"object-storage-go/data-server/rabbitmq"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

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
	mq := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer mq.Close()

	mq.Publish("dataServer", name)

	ch := mq.Consume()
	//1秒之后关闭临时队列，避免无限期的等待，如果不写，则下面的msg可能就一致接收不到，就一直等下去
	go func() {
		time.Sleep(time.Duration(1) * time.Second)
		mq.Close()
	}()
	msg := <- ch

	s, _ := strconv.Unquote(string(msg.Body))

	return s
}

func Exist(name string) bool {
	return Locate(name) != ""
}
