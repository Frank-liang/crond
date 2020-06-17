package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	var (
		config clientv3.Config
		err    error
		client *clientv3.Client
		kv     clientv3.KV
		getRes *clientv3.GetResponse
	)

	config = clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	kv = clientv3.NewKV(client)

	if getRes, err = kv.Get(context.TODO(), "/cron/job/", clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(getRes.Kvs)
	}
}
