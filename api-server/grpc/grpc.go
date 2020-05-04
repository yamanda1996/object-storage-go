package grpc

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"object-storage-go/api-server/etcd"
	"object-storage-go/common/common_constant"
	"object-storage-go/common/common_grpc/objectpb"
	"object-storage-go/common/common_grpc/schedulerpb"
)


func init()  {

}

func GetSchedulerClient() (*schedulerpb.SchedulerClient, error) {

	schurl, err := etcd.GetServerUrl(common_constant.SCHEDULER_SERVER_ETCD_PREFIX)
	if err != nil {
		log.Error("get scheduler url from etcd server failed")
		return nil, err
	}
	conn, err := grpc.Dial(schurl, grpc.WithInsecure())
	if err != nil {
		log.Error("dial to scheduler server failed")
		return nil, err
	}
	client := schedulerpb.NewSchedulerClient(conn)
	return &client, nil
}

func GetDataServerClient(server *schedulerpb.DataServer) (*objectpb.ObjectClient, error) {
	conn, err := grpc.Dial(server.Address+":"+server.Port, grpc.WithInsecure())
	if err != nil {
		log.Error("dial to data server failed")
		return nil, err
	}
	client := objectpb.NewObjectClient(conn)
	return &client, nil
}