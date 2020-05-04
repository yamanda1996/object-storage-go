package model

import (
	"fmt"
	"gopkg.in/gcfg.v1"
	"object-storage-go/data-server/utils"
	"os"
)

var (
	Conf = &Config{}
)

type Config struct {
	RabbitMqConfig 			RabbitMqConfig
	DataServerConfig 		DataServerConfig
	EtcdServerConfig 		EtcdServerConfig
	LogConfig				LogConfig
}


type RabbitMqConfig struct {
	RabbitMqAddress 					string
	RabbitMqPort						string
	RabbitMqUser						string
	RabbitMqPwd 						string
}

type DataServerConfig struct {
	DataServerAddress 				string
	DataServerPort 					string
	DataServerIndex 				string
	DataServerEtcdPrefix 			string
	DataServerRegion 				string
	StorageRoot 					string
	HeartBeatTimeInSeconds			string
	BucketNumber 					int
}

type EtcdServerConfig struct {
	EtcdServerAddress					string
	EtcdServerPort						string
	EtcdServerTimeout					string
}

type LogConfig struct {
	LogFilePath							string
	LogLevel 							string
}

func init() {
	err := gcfg.ReadFileInto(Conf, utils.DATA_SERVER_CONFIG_FILEPATH)
	if err != nil {
		fmt.Println("Read config file to Config failed")
		os.Exit(-1)
	}
}