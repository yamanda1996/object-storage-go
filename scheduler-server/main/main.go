package main

import "fmt"

func main()  {
	// 分布式对象系统中的调度模块，用来负责文件存储调度，模块之间使用grpc调用，元数据存储在etcd中
	// 目前模块之间使用的是restful http请求调用，元数据存储elastic search 中 TODO
	fmt.Println("this is a scheduler module for total project")

}
