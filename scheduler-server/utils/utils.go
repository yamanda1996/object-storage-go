package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"object-storage-go/common/common_etcd"
	"object-storage-go/scheduler-server/config"
	"os"
)

var SchedulerServerUrl 		string
var SchedulerServerPrefix 	string

var Register 		*common_etcd.ServiceRegister

//func init () {
//	err := InitLog()
//	if err != nil {
//		fmt.Println("init log failed")
//		os.Exit(1)
//	}
//
//	err = InitConfig()
//	if err != nil {
//		log.Error("init config file failed")
//		os.Exit(1)
//	}
//
//	SchedulerServerPrefix = model.Config.SchedulerServerConfig.SchedulerServerEtcdPrefix +
//		model.Config.SchedulerServerConfig.SchedulerServerIndex
//	SchedulerServerUrl = model.Config.SchedulerServerConfig.SchedulerServerAddress + ":" +
//		model.Config.SchedulerServerConfig.SchedulerServerPort
//	err = RegisterToEtcdServer(SchedulerServerPrefix, SchedulerServerUrl)
//	if err != nil {
//		log.Error("register data server to etcd server failed")
//		os.Exit(1)
//	}
//}

func InitLog() error {
	log.SetFormatter(&log.TextFormatter{})
	file, err := os.OpenFile(config.Conf.LogConfig.LogFilePath,
		os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Init log failed")
		return fmt.Errorf("Init log failed.")
	}
	log.SetOutput(file)
	log.SetLevel(log.DebugLevel)
	return nil
}

//func InitConfig() error {
//
//	err := gcfg.ReadFileInto(&model.Config, constant.SCHEDULER_SERVER_CONFIG_FILEPATH)
//	if err != nil {
//		log.Errorf("read config file [%s] failed", constant.SCHEDULER_SERVER_CONFIG_FILEPATH)
//		os.Exit(1)
//	}
//	return nil
//}
//
//func RegisterToEtcdServer(k,v string) error {
//
//	endpoints := []string{model.Config.EtcdServerConfig.EtcdServerAddress + ":" +
//		model.Config.EtcdServerConfig.EtcdServerPort}
//
//	var err error
//	Register, err = common_etcd.NewServiceRegister(endpoints, common_constant.ETCD_DIAL_TIMEOUT)
//	if err != nil {
//		log.Error(err.Error())
//		return err
//	}
//
//	err = Register.RegisterService(k, v)
//	if err != nil {
//		log.Error(err.Error())
//		return err
//	}
//	return nil
//}