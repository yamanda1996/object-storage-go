package main

import (
	"fmt"
	"object-storage-go/scheduler-server/cron"
	"time"

	"object-storage-go/scheduler-server/config"
	_ "object-storage-go/scheduler-server/etcd"
	"object-storage-go/scheduler-server/utils"

	log "github.com/sirupsen/logrus"
)

func main() {
	// 分布式对象系统中的调度模块，用来负责文件存储调度，模块之间使用grpc调用，元数据存储在etcd中
	// 目前模块之间使用的是restful http请求调用，元数据存储elastic search 中 TODO

	//flag.Parse()
	err := config.Init()
	if err != nil {
		fmt.Println("Init config failed")
		return
	}

	err = utils.InitLog()
	if err != nil {
		fmt.Println("Init log failed")
		return
	}
	log.Debug("Init log success")

	cron.DoCron()

	time.Sleep(1 * time.Hour)
}
