package api

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"object-storage-go/common/grpc/schedulerpb"
	"object-storage-go/scheduler-server/heartbeat"
	"strings"
)

type Scheduler struct {
}

func (s *Scheduler) GetDataServer(c context.Context, request *schedulerpb.CommonRequest) (*schedulerpb.DataServer, error) {
	log.Info("scheduler start to get data server")

	server := new (schedulerpb.DataServer)
	server.Address = ""
	server.Port = ""
	ds := heartbeat.GetDataServers()

	n := len(ds)
	if n == 0 {
		return server, fmt.Errorf("data server is null")
	}
	splits := strings.Split(ds[rand.Intn(n)], ":")
	server.Address = splits[0]
	server.Port = splits[1]

	return server, nil
}


