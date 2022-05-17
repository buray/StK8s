package main

import (
	"context"
	"fmt"
	etcdCt "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

var etcdServerList = []string{"http://192.168.123.161:2379", "http://192.168.123.160:2379", "http://192.168.123.162:2379"}

func main() {
	etcdClient, err := etcdCt.New(etcdCt.Config{
		Endpoints:   etcdServerList,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer etcdClient.Close()

	for _, i := range etcdServerList {
		clientStatus, err := etcdClient.Status(context.Background(), i)
		if err != nil {
			log.Fatal(err)
		}
		println(clientStatus.Version)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ttl, err := etcdClient.Grant(ctx, 20)
	if err != nil {
		log.Fatal(err)
	}
	leaseId := ttl.ID
	fmt.Println("leaseId", leaseId)
	_, err = etcdClient.Put(ctx, "sample_key2", "sample_value", etcdCt.WithLease(leaseId))

	if err != nil {
		log.Fatal(err)
	}

	result, err := etcdClient.Get(ctx, "sample_key1")
	if err != nil {
		log.Fatal(err)
	}
	cancel()

	println(string(result.Kvs[0].Value))
	println(string(result.Kvs[0].Key))

	fmt.Println("this is man function")
}
