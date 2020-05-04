package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gcfg.v1"
	"object-storage-go/api-server/constant"
	"object-storage-go/api-server/model"
	"os"
)

func init() {
	err := InitLog()
	if err != nil {
		fmt.Println("init log failed")
		os.Exit(1)
	}

	err = InitConfig()
	if err != nil {
		log.Error("init config file failed")
		os.Exit(1)
	}
}

func InitLog() error {
	log.SetFormatter(&log.TextFormatter{})
	file, err := os.OpenFile(constant.LOG_FILEPATH, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("init log failed")
		return fmt.Errorf("init log failed")
	}
	log.SetOutput(file)
	log.SetLevel(log.DebugLevel)
	return nil
}

func InitConfig() error {

	err := gcfg.ReadFileInto(&model.Config, constant.API_SERVER_CONFIG_FILEPATH)
	if err != nil {
		log.Errorf("read config file [%s] failed", constant.API_SERVER_CONFIG_FILEPATH)
		os.Exit(1)
	}
	return nil
}



func GetFromHeader(context *gin.Context, h string) []string {
	for k, v := range context.Request.Header {
		if k == h {
			return v
		}
	}
	return nil
}
