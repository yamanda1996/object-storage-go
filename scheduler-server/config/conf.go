package config

import (
	"fmt"
	"gopkg.in/gcfg.v1"
	"object-storage-go/scheduler-server/constant"
)

var (
	Conf = &Config{}
)

type Config struct {
	RabbitMqConfig 			RabbitMqConfig
	SchedulerServerConfig 	SchedulerServerConfig
	EtcdServerConfig 		EtcdServerConfig
	LogConfig				LogConfig
}


type RabbitMqConfig struct {
	RabbitMqAddress 					string
	RabbitMqPort						string
	RabbitMqUser						string
	RabbitMqPwd 						string
}

type SchedulerServerConfig struct {
	SchedulerServerAddress 				string
	SchedulerServerPort 				string
	SchedulerServerIndex 				string
	SchedulerServerEtcdPrefix 			string
	DataServerEtcdPrefix				string
	SchedulerServerRegion 				string
	HeartBeatTimeInSeconds				string
}

type EtcdServerConfig struct {
	EtcdServerAddress					string
	EtcdServerPort						string
	EtcdServerTimeout 					string
}

type LogConfig struct {
	LogFilePath							string
	LogLevel 							string
}

func Init() error {
	err := gcfg.ReadFileInto(Conf, constant.SCHEDULER_SERVER_CONFIG_FILEPATH)
	if err != nil {
		fmt.Println("Read config file to Config failed")
		return err
	}
	return nil
}