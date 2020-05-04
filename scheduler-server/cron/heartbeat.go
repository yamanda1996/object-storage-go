package cron

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"object-storage-go/data-server/model"
	"object-storage-go/scheduler-server/config"
	"object-storage-go/scheduler-server/etcd"
	"strconv"
	"time"
)

var DataServerMap map[string]string

func HeartBeat(ctx context.Context) {
	prefix := config.Conf.SchedulerServerConfig.DataServerEtcdPrefix + "/" +
		config.Conf.SchedulerServerConfig.SchedulerServerRegion + "/heartbeat/"
	t, _ := strconv.ParseInt(config.Conf.SchedulerServerConfig.HeartBeatTimeInSeconds, 10, 64)
	ticker := time.NewTicker(time.Duration(t) * time.Second)
	defer ticker.Stop()
	for {
		<- ticker.C
		log.Debugf("Start to deal heartbeat")
		log.Debugf("Prefix is %s", prefix)
		timestamp := time.Now().Unix()
		log.Debugf("Timestamp %s", string(timestamp))

		resp, err := etcd.EtcdClient.Get(ctx, prefix, clientv3.WithPrefix())
		if err != nil {
			log.Debug("Receive from etcd failed...")
		}

		log.Debugf("Receive from etcd %d", len(resp.Kvs))

		for _, server := range resp.Kvs {
			hb := model.HeartBeat{}
			json.Unmarshal(server.Value, &hb)

			if timestamp - hb.Timestamp > 15 {
				log.Warn("data server failed")
			} else {
				log.Info("data server is normal")
			}

		}
	}
}
