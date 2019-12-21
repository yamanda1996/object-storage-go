package rpc

import (
	"google.golang.org/grpc"
	"object-storage-go/common/grpc/schedulerpb"
)

var SchedulerClient schedulerpb.SchedulerClient

func init()  {
	conn, _ := grpc.Dial("localhost:8992", grpc.WithInsecure())
	SchedulerClient = schedulerpb.NewSchedulerClient(conn)
}
