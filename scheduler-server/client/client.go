package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"object-storage-go/scheduler-server/router"
	"time"
)

func main()  {
	fmt.Println("start grpc client")
	conn, _ := grpc.Dial("localhost:8899", grpc.WithInsecure())

	client := router.NewRouteGuideClient(conn)

	p1 := &router.Point{
		Latitude:111,
		Longitude:222,
	}
	// 简单
	//f1, _ := client.GetFeature(context.Background(), p1)
	//fmt.Printf("receive from server %v\n", f1)

	// stream 返回
	//stream, _ := client.ListFeatures(context.Background(), &router.Rectangle{
	//	Lo: p1,
	//	Hi: p1,
	//})
	//
	//for {
	//	f1, err := stream.Recv()
	//	if err == io.EOF {
	//		fmt.Println("receive success")
	//		break
	//	}
	//	fmt.Printf("receive %v\n", f1)
	//}

	// stream 输入
	//stream, _ := client.RecordRoute(context.Background())
	//stream.Send(p1)
	//fmt.Println("send p1 success")
	//time.Sleep(3 * time.Second)
	//stream.Send(p1)
	//fmt.Println("send p2 success")
	//stream.CloseAndRecv()
	//fmt.Println("send success")

	// 双向stream
	stream, _ := client.RouteChan(context.Background())
	waitc := make(chan struct{})

	go func() {
		// receive
		for {
			n, err := stream.Recv()
			fmt.Printf("receive from server %v\n", n)
			if err == io.EOF {
				fmt.Println("receive finish")
				//waitc <- struct{}{}
				close(waitc)
				return
			}
		}
	}()

	n1 := &router.RouteNote{
		Loate:p1,
		Message:"n1",
	}

	stream.Send(n1)
	fmt.Println("client send n1 success")
	time.Sleep(3 * time.Second)
	stream.Send(n1)
	fmt.Println("client send n2 success")
	stream.CloseSend()
	<- waitc
	defer conn.Close()
}
