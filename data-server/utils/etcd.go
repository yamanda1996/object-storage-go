package utils

import (
	"go.etcd.io/etcd/clientv3"
	"object-storage-go/data-server/model"
	"strconv"
	"time"
)

var (
	EtcdClient clientv3.KV

	)

func init() {
	endpoint := model.Conf.EtcdServerConfig.EtcdServerAddress +
		":" + model.Conf.EtcdServerConfig.EtcdServerPort
	timeout, _ := strconv.ParseInt(
		model.Conf.EtcdServerConfig.EtcdServerTimeout, 10, 64)

	c := clientv3.Config{
		Endpoints: []string{endpoint},
		DialTimeout: time.Duration(timeout) * time.Second,
	}
	client, _ := clientv3.New(c)
	EtcdClient = clientv3.NewKV(client)
}