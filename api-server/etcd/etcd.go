package etcd

import (
	"go-common/library/log"
	"math/rand"
	"object-storage-go/api-server/model"
	"object-storage-go/common/common_constant"
	"object-storage-go/common/common_etcd"
	"os"
)

var ApiServerUrl 	string
var ApiServerPrefix string

var Register 		*common_etcd.ServiceRegister
var DiscoveryClient *common_etcd.DiscoveryClient

var EtcdEndpoints 	[]string

func init()  {
	ApiServerPrefix = model.Config.ApiServerConfig.ApiServerEtcdPrefix + model.Config.ApiServerConfig.ApiServerIndex
	ApiServerUrl = model.Config.ApiServerConfig.ApiServerAddress + ":" +
		model.Config.ApiServerConfig.ApiServerPort
	EtcdEndpoints = []string{model.Config.EtcdServerConfig.EtcdServerAddress + ":" +
		model.Config.EtcdServerConfig.EtcdServerPort}
	err := RegisterToEtcdServer(ApiServerPrefix, ApiServerUrl)
	if err != nil {
		log.Error("register api server to etcd server failed")
		os.Exit(1)
	}

	DiscoveryClient, err = common_etcd.NewDiscoveryClient(EtcdEndpoints)
}

func RegisterToEtcdServer(k,v string) error {

	var err error
	Register, err = common_etcd.NewServiceRegister(EtcdEndpoints, common_constant.ETCD_DIAL_TIMEOUT)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	err = Register.RegisterService(k, v)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func GetServerUrl(prefix string) (string, error) {
	urls, err := DiscoveryClient.DiscoveryService(prefix)
	if err != nil {
		return "", err
	}
	return urls[rand.Intn(len(urls))], nil
}