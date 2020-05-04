# Object Storage
this is a distributed object storage system in go


## api-server
对用户暴露接口，主要功能包括
* 文件上传、下载
* 

## scheduler-server
中央调度模块，主要功能包括
* 文件元数据存储
* 分配文件存储节点（不同的策略）
* 数据校验和去重
* 数据冗余和即时修复
* 数据维护

## data-server
负责存储节点，主要功能包括
* 文件存储
* 文件压缩存储

## feature
1、模块内部调用使用grpc
2、模块使用etcd用作元数据存储以及服务注册与发现
3、


## find object
![](https://github.com/yamanda1996/object-storage-go/blob/master/images/object_storage_struct_pic1.png)

## ETCD
ETCD用来作为心跳检测的工具，data-server定时（60s）向etcd的指定路径中汇报心跳，scheduler收集这些心跳信息并维护data-server的状态信息，如果超过一个三个周期未收到心跳，则认为data-server处于error状态。


