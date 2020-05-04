package etcd

import (
	"go.etcd.io/etcd/clientv3"
	"object-storage-go/scheduler-server/config"
	"strconv"
	"time"
)

var (
	EtcdClient clientv3.KV

)

func init() {
	endpoint := config.Conf.EtcdServerConfig.EtcdServerAddress +
		":" + config.Conf.EtcdServerConfig.EtcdServerPort
	timeout, _ := strconv.ParseInt(
		config.Conf.EtcdServerConfig.EtcdServerTimeout, 10, 64)

	c := clientv3.Config{
		Endpoints: []string{endpoint},
		DialTimeout: time.Duration(timeout) * time.Second,
	}
	client, _ := clientv3.New(c)
	EtcdClient = clientv3.NewKV(client)
}