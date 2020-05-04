package server

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"net"
	"object-storage-go/scheduler-server/router"
	"time"
)

type RouteServer struct {}

func (s *RouteServer) GetFeature(context context.Context, point *router.Point) (*router.Feature, error) {
	fmt.Println("server start get feature")
	p := &router.Point{
		Latitude:12,
		Longitude:24,
	}
	return &router.Feature{Location:p}, nil
}

func (s *RouteServer) ListFeatures(rectangle *router.Rectangle, stream router.RouteGuide_ListFeaturesServer) error {
	fmt.Println("server start list feature")
	p1 := &router.Point{
		Latitude:11,
		Longitude:12,
	}

	f1 := &router.Feature{
		Name:"f1",
		Location:p1,
	}

	f2 := &router.Feature{
		Name:"f2",
		Location:p1,
	}

	stream.Send(f1)
	fmt.Println("send f1 success")
	time.Sleep(3 * time.Second)
	stream.Send(f2)
	fmt.Println("send f2 success")
	return nil
}

func (s *RouteServer) RecordRoute(stream router.RouteGuide_RecordRouteServer) error {
	startTime := time.Now()
	fmt.Printf("start record route %v\n", startTime)
	for {
		point, err := stream.Recv()
		fmt.Printf("receive point %v\n", point)
		if err == io.EOF {
			endTime := time.Now()
			fmt.Printf("end record route %v\n", endTime)
			return stream.SendAndClose(&router.RouteSummary{
				PointCount:1,
				FeatureCount:1,
				Distance:1,
				ElapsedTime:1,
			})
		}
	}
}

func (s *RouteServer) RouteChan(stream router.RouteGuide_RouteChanServer) error {
	for {
		in, err := stream.Recv()
		fmt.Printf("receive from route chat %v\n", in)
		if err == io.EOF {
			return nil
		}
		p1 := &router.Point{
			Latitude:11,
			Longitude:12,
		}
		fmt.Println("do something")
		stream.Send(&router.RouteNote{
			Loate:p1,
			Message:"success route chat",
		})
	}
}

func main()  {
	fmt.Println("start grpc server")
	flag.Parse()
	lis, _ := net.Listen("tcp", "localhost:8899")
	grpcServer := grpc.NewServer()
	router.RegisterRouteGuideServer(grpcServer, &RouteServer{})

	grpcServer.Serve(lis)

}
