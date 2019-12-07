package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gcfg.v1"
	"object-storage-go/data-server/constant"
	"object-storage-go/data-server/model"
	"os"
)

func InitLog() error {
	log.SetFormatter(&log.JSONFormatter{})
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

	err := gcfg.ReadFileInto(&model.Config, constant.DATA_SERVER_CONFIG_FILEPATH)
	if err != nil {
		log.Errorf("read config file [%s] failed", constant.DATA_SERVER_CONFIG_FILEPATH)
		os.Exit(1)
	}
	return nil
}
