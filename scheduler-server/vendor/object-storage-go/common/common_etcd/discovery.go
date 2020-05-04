package common_etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"sync"
	"time"
)

type DiscoveryClient struct {
	client        *clientv3.Client
	serverList    map[string]string
	lock          sync.Mutex
}

func NewDiscoveryClient(endpoints []string) (*DiscoveryClient, error) {
	conf := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}
	cli, err := clientv3.New(conf)
	if err != nil {
		return nil, fmt.Errorf("create etcd client failed")
	}

	client := &DiscoveryClient{
		client:cli,
		serverList:make(map[string]string),
	}

	return client, nil
}

func (c * DiscoveryClient) DiscoveryService(prefix string) ([]string ,error){
	resp, err := c.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("discovery service with prefix %s failed", prefix)
	}

	serverList := c.getServerList(resp)

	go c.watcher(prefix)
	return serverList ,nil
}


func (c *DiscoveryClient) watcher(prefix string) {
	watchChan := c.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for resp := range watchChan {
		for _, e := range resp.Events {
			switch e.Type {
			case mvccpb.PUT:
				c.SetServiceList(string(e.Kv.Key),string(e.Kv.Value))
			case mvccpb.DELETE:
				c.DelServiceList(string(e.Kv.Key))
			}
		}
	}
}

func (c *DiscoveryClient) getServerList(resp *clientv3.GetResponse) []string {
	serverList := make([]string,0)
	if resp == nil || resp.Kvs == nil {
		return serverList
	}
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			c.SetServiceList(string(resp.Kvs[i].Key),string(resp.Kvs[i].Value))
			serverList = append(serverList, string(v))
		}
	}
	return serverList
}

func (c *DiscoveryClient) SetServiceList(k,v string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.serverList[k] = string(v)
	fmt.Println("set key :",k,"val:",v)
}

func (c *DiscoveryClient) DelServiceList(k string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.serverList, k)
	fmt.Println("del key:", k)
}