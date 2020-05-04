package main

import (
	"github.com/gin-gonic/gin"
	"object-storage-go/api-server/model"
	"object-storage-go/api-server/object"
	_ "object-storage-go/api-server/etcd"
	_ "object-storage-go/api-server/grpc"
)

func main()  {
	router := gin.Default()

	//router.GET("/objects/:filename", object.DownloadFile)
	router.POST("/objects/:filename", object.UploadFile)
	//router.DELETE("/objects/:filename", object.DeleteFile)

	//router.GET("/version/:filename", version.GetVersion)

	//router.GET("/locate/:filename", locate.LocateFile)

	router.Run(":" + model.Config.ApiServerConfig.ApiServerPort)
}
