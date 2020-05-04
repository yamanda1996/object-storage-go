package sch_grpc

import (
	"context"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"object-storage-go/common/common_constant"
	"object-storage-go/common/common_grpc/schedulerpb"
	"object-storage-go/scheduler-server/etcd"
	"strings"
)

type Scheduler struct {
}

func (s *Scheduler) GetDataServer(c context.Context, request *schedulerpb.CommonRequest) (*schedulerpb.DataServer, error) {

	log.Info("scheduler start to get data server")
	list, err := etcd.DiscoveryClient.DiscoveryService(common_constant.DATA_SERVER_ETCD_PREFIX)
	if err != nil {
		log.Error("find data server list failed")
		return nil, err
	}
	splits := strings.Split(list[rand.Intn(len(list))], ":")

	return &schedulerpb.DataServer{Address:splits[0], Port:splits[1],}, nil
}