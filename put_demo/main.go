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
		putRes *clientv3.PutResponse
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

	if putRes, err = kv.Put(context.TODO(), "/cron/job/job", "dog", clientv3.WithPrevKV()); err != nil {
		fmt.Print(err)
	} else {
		fmt.Println("revition:", putRes.Header.Revision)
		if putRes.PrevKv != nil {
			fmt.Print(string(putRes.PrevKv.Value))
		}

	}

}
