package main

import (
	"common-library/log"
	"fmt"
	"object-storage-go/scheduler-server/conf"
)

func main()  {
	if err := conf.Init(); err != nil {
		fmt.Println(err.Error())
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	log.Info("start scheduler server")
	log.Init(conf.Conf.Log)
}
