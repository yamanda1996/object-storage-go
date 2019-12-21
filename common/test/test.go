package main

import (
	"fmt"
	"object-storage-go/common/etcd"
)

var (
	endpoints = []string{"127.0.0.1:2379"}
	timeout int64 = 5
)

func main()  {
	register, err := etcd.NewServiceRegister(endpoints, timeout)
	if err != nil {
		fmt.Println("test new service register failed")
		return
	}

	register.RegisterService("/hello/world1", "人之初")
	register.RegisterService("/hello/world2", "性本善")

	client, _ := etcd.NewDiscoveryClient(endpoints)

	result, _ := client.DiscoveryService("/hello")
	fmt.Printf("get from etcd server %v\n", result)
}


