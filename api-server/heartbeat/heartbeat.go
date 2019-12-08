package heartbeat

import (
	log "github.com/sirupsen/logrus"
	"math/rand"
	"object-storage-go/api-server/rabbitmq"
	"strconv"
	"sync"
	"time"
)

var dataServers = make(map[string]time.Time)
var mutex sync.Mutex

func ListenHeartBeat()  {

	mq := rabbitmq.New(rabbitmq.GetRabbitMqDialUrl())
	defer mq.Close()

	mq.Bind("apiServer")
	ch := mq.Consume()

	go removeExpiredDataServer()

	//收到消息之后，更新map中对应的时间
	for msg := range ch {
		dataServer, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			log.Error("heart beat receive msg failed")
			continue
		}

		//操作公有对象，必须加锁
		mutex.Lock()
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}

func removeExpiredDataServer()  {
	for {
		time.Sleep(time.Duration(5) * time.Second)
		mutex.Lock()

		for s, t := range dataServers {
			//上一次心跳的时间是10秒钟之前，则任务该data-server已经宕机
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}

		mutex.Unlock()
	}
}

//获取所有的dataServer信息
func getDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()

	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	return ds
}

//随机选择一个dataServer
func ChooseRandomDataServer() string {
	ds := getDataServers()

	n := len(ds)
	if n == 0 {
		return ""
	}
	return ds[rand.Intn(n)]
}