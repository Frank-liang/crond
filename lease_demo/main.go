package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	var (
		config                 clientv3.Config
		err                    error
		client                 *clientv3.Client
		kv                     clientv3.KV
		getRes                 *clientv3.GetResponse
		lease                  clientv3.Lease
		id                     clientv3.LeaseID
		leaseGRes              *clientv3.LeaseGrantResponse
		putRes                 *clientv3.PutResponse
		leaaseKeepAliveRes     *clientv3.LeaseKeepAliveResponse
		leaaseKeepAliveResChan <-chan *clientv3.LeaseKeepAliveResponse
	)

	config = clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	lease = clientv3.NewLease(client)

	if leaseGRes, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
	}

	id = leaseGRes.ID
	kv = clientv3.NewKV(client)

	if leaaseKeepAliveResChan, err = lease.KeepAlive(context.TODO(), id); err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			select {
			case leaaseKeepAliveRes = <-leaaseKeepAliveResChan:
				if leaaseKeepAliveRes == nil {
					fmt.Println("租约已经失效了")
					goto END
				} else {

					fmt.Println("收到自动续租应答:", leaaseKeepAliveRes.ID)
				}

			}
		}
	END:
	}()

	if putRes, err = kv.Put(context.TODO(), "/cron/lock/job", "", clientv3.WithLease(id)); err != nil {
		fmt.Println(err)
	}
	fmt.Println("success: ", putRes.Header.Revision)

	for {
		if getRes, err = kv.Get(context.TODO(), "/cron/lock/job"); err != nil {
			fmt.Println(err)
			return
		}

		if getRes.Count == 0 {
			fmt.Println("Key has been deleted")
			break
		}
		time.Sleep(2 * time.Second)
		now := time.Now()
		fmt.Println(now)
		fmt.Println(getRes.Kvs)
	}
}
