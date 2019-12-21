package etcd

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
	"time"
)

type ServiceRegister struct {
	client        			*clientv3.Client
	lease         			clientv3.Lease
	leaseResp     			*clientv3.LeaseGrantResponse
	cancel    				func()
	keepAliveChan 			<-chan *clientv3.LeaseKeepAliveResponse
	key           			string
}

func NewServiceRegister(endpoints []string, timeout int64) (*ServiceRegister, error) {
	conf := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}
	cli, err := clientv3.New(conf)
	if err != nil {
		return nil, fmt.Errorf("create etcd client failed")
	}

	register := &ServiceRegister{
		client: cli,
	}

	err = register.setLease(timeout)
	if err != nil {
		return nil, fmt.Errorf("set lease failed")
	}

	go register.ListenLeaseRespChan()
	return register, nil
}

// set lease
func (r *ServiceRegister) setLease(timeout int64) error {
	lease := clientv3.NewLease(r.client)

	resp, err := lease.Grant(context.TODO(), timeout)
	if err != nil {
		return fmt.Errorf("set lease timeout failed")
	}

	ctx, cancel := context.WithCancel(context.TODO())
	keepAliveChan, err := lease.KeepAlive(ctx, resp.ID)
	if err != nil {
		return fmt.Errorf("get lease keep alive response failed")
	}

	r.lease = lease
	r.cancel = cancel
	r.leaseResp = resp
	r.keepAliveChan = keepAliveChan
	return nil
}

// renew lease
func (r *ServiceRegister) ListenLeaseRespChan()  {
	for {
		select {
		case keepAlive := <- r.keepAliveChan:
			if keepAlive == nil {
				fmt.Println("renew lease is closed")
				return
			} else {
				fmt.Println("renew lease success")
			}
		}
	}
}

func (r *ServiceRegister) RegisteService() error {
	
}
