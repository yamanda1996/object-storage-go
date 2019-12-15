package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"object-storage-go/data-server/heartbeat"
	"object-storage-go/data-server/locate"
	"object-storage-go/data-server/model"
	"object-storage-go/data-server/objects"
	"object-storage-go/data-server/temp"
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
	// 向apiServer中发送心跳告知数据服务的存活
	go heartbeat.HeartBeat()
	// 监听dataServer中的消息
	go locate.StartLocate()

	router := gin.Default()
	router.GET("/objects/:filename", objects.GetObject)
	router.POST("/objects/:filename", objects.PutObject)

	router.POST("/temp/:hash", temp.PostTemp)
	router.PATCH("/temp/:uuid", temp.PatchTemp)
	router.PUT("/temp/:uuid", temp.PutTemp)
	router.DELETE("/temp/:uuid", temp.DeleteTemp)

	router.Run(":" + strconv.Itoa(model.Config.DataServerConfig.DataServerPort))

}


