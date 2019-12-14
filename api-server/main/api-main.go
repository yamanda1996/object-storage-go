package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"object-storage-go/api-server/heartbeat"
	"object-storage-go/api-server/locate"
	"object-storage-go/api-server/model"
	"object-storage-go/api-server/object"
	"object-storage-go/api-server/utils"
	"object-storage-go/api-server/version"
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
	log.Debug("start api server")
	log.Debugf("rabbit mq address [%s]", model.Config.RabbitMqConfig.RabbitMqAddress)
	log.Debug("api server address: " + model.Config.ApiServerConfig.ApiServerAddress)
	go heartbeat.ListenHeartBeat()
	router := gin.Default()

	router.GET("/objects/:filename", object.DownloadFile)
	router.POST("/objects/:filename", object.UploadFile)
	router.DELETE("/objects/:filename", object.DeleteFile)

	router.GET("/version/:filename", version.GetVersion)

	router.GET("/locate/:filename", locate.LocateFile)

	router.Run(":" + strconv.Itoa(model.Config.ApiServerConfig.ApiServerPort))

}
