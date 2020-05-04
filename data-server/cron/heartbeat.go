package cron

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"object-storage-go/data-server/model"
	"object-storage-go/scheduler-server/etcd"
	"strconv"
	"time"
)

func HeartBeat(ctx context.Context) {
	prefix := model.Conf.DataServerConfig.DataServerEtcdPrefix + "/" +
		model.Conf.DataServerConfig.DataServerRegion + "/heartbeat/" +
		model.Conf.DataServerConfig.DataServerIndex
	t, _ := strconv.ParseInt(model.Conf.DataServerConfig.HeartBeatTimeInSeconds, 10, 64)
	ticker := time.NewTicker(time.Duration(t) * time.Second)
	defer ticker.Stop()
	for {
		<- ticker.C
		log.Debugf("Start to send heartbeat, data server %s", model.Conf.DataServerConfig.DataServerIndex)
		hb, _ := prepHearBeat()
		log.Debugf("prefix is %s", prefix)
		log.Debugf("heartbeat is %s", string(hb))
		etcd.EtcdClient.Put(ctx, prefix, string(hb))
		log.Debug("Send heartbeat to etcd success")
	}
}

// should not error happen
func prepHearBeat() ([]byte, error) {
	timestamp := time.Now().Unix()
	hb := model.HeartBeat{
		DataServerAddress:"111.111.111.111",
		Timestamp:timestamp,
		CpuUsage:0.1,
		MemUsage:0.2,
		DiskUsage:0.3,
	}
	return json.Marshal(hb)
}

func CpuUsage() (float64, error) {
	return 0, nil
}



