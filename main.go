package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

func main() {

	config := api.DefaultConfig()
	c, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("new client: %v", err)
	}

	kv := c.KV()

	p1 := &api.KVPair{Key: "app/k1", Value: []byte("v1")}
	p2 := &api.KVPair{Key: "app/k2", Value: []byte("v2")}

	_, err = kv.Put(p1, nil)
	if err != nil {
		log.Fatalf("put %v: %v", p1, err)
	}
	_, err = kv.Put(p2, nil)
	if err != nil {
		log.Fatalf("put %v: %v", p2, err)
	}

	got, meta, err := kv.Get("app/k1", nil)
	if err != nil {
		log.Fatalf("get: %v", err)
	}
	fmt.Printf("got: %#v\n", got)
	fmt.Printf("meta: %#v\n", meta)

}
