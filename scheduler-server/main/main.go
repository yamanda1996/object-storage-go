package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"object-storage-go/scheduler-server/heartbeat"
	"object-storage-go/scheduler-server/model"
	"object-storage-go/scheduler-server/utils"
	"os"
	"strconv"
)

func init() {
	err := utils.InitLog()
	if err != nil {
		fmt.Println("init log failed")
		os.Exit(1)
	}

	err = utils.InitConfig()
	if err != nil {
		fmt.Println("init config file failed")
		os.Exit(1)
	}
}

func main()  {
	// 分布式对象系统中的调度模块，用来负责文件存储调度，模块之间使用grpc调用，元数据存储在etcd中
	// 目前模块之间使用的是restful http请求调用，元数据存储elastic search 中 TODO
	//fmt.Println("this is a scheduler module for total project")

	flag.Parse()
	net.Listen("tcp", ":" + strconv.Itoa(model.Config.SchedulerServerConfig.SchedulerServerPort))
	go heartbeat.ListenHeartBeat()

	grpcServer := grpc.NewServer()



}
