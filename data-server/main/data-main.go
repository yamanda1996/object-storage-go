package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"object-storage-go/data-server/heartbeat"
	"object-storage-go/data-server/model"
	"object-storage-go/data-server/objects"
	"object-storage-go/data-server/utils"
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

func main() {
	log.Debug("start main")
	go heartbeat.HeartBeat()

	router := gin.Default()
	router.GET("/objects/:filename", objects.GetObject)
	router.POST("/objects/:filename", objects.UploadObject)

	router.Run(":" + strconv.Itoa(model.Config.DataServerConfig.DataServerPort))

}


