package main

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"object-storage-go/data-server/utils"
	"object-storage-go/data-server/model"
	"os"
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

func main() {
	log.Debug("start main")
	log.Debug(model.Config.DataServerConfig.RabbitMqAddress)


	//go heartbeat.HeartBeat()
	//go locate.StartLocate()
	//
	//http.HandleFunc("/object/", objects.Handler)
	//log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))

}


