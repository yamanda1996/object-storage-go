package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	endpoints = []string{"127.0.0.1:2379"}
	timeout int64 = 5
)

func main()  {

	config := clientv3.Config{
		Endpoints: endpoints,
		DialTimeout: time.Duration(timeout) * time.Second,
	}

	client, _ := clientv3.New(config)
	defer client.Close()

	kv := clientv3.NewKV(client)

	ctx, _ := context.WithTimeout(context.TODO(),
		time.Duration(timeout) * time.Second)

	putResp, _ := kv.Put(ctx, "/hello", "world666", clientv3.WithPrevKV())

	getResp, _ := kv.Get(ctx, "/hello", clientv3.WithPrevKV())

	fmt.Println(putResp.PrevKv)
	fmt.Println(getResp.Kvs)
	//register, err := common_etcd.NewServiceRegister(endpoints, timeout)
	//if err != nil {
	//	fmt.Println("test new service register failed")
	//	return
	//}
	//
	//register.RegisterService("/hello/world1", "人之初")
	//register.RegisterService("/hello/world2", "性本善")

	//client, _ := common_etcd.NewDiscoveryClient(endpoints)
	//
	//result, _ := client.DiscoveryService("/hello")
	//fmt.Printf("get from etcd server %v\n", result)
}


