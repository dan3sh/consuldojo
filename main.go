package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

func main() {

	c, err := newClient()
	if err != nil {
		log.Fatalf("new client: %v", err)
	}

	err = setupKV(c)
	if err != nil {
		log.Fatalf("setup: %v", err)
	}
	defer func() {
		err = teardownKV(c)
		if err != nil {
			log.Fatalf("teardown: %v", err)
		}
	}()
	kv := c.KV()
	got, meta, err := kv.Get("app/k1", nil)
	if err != nil {
		log.Fatalf("get: %v", err)
	}
	fmt.Printf("got: %#v\n", got)
	fmt.Printf("meta: %#v\n", meta)
}

func newClient() (*api.Client, error) {
	config := api.DefaultConfig()
	return api.NewClient(config)
}

const defaultKeyPrefix = "app"

func setupKV(c *api.Client) error {
	pairs := []struct {
		k string
		v string
	}{
		{"k1", "v1"},
		{"k2", "v2"},
	}
	for _, p := range pairs {
		kv := &api.KVPair{
			Key:   fmt.Sprintf("%s/%s", defaultKeyPrefix, p.k),
			Value: []byte(p.v),
		}
		_, err := c.KV().Put(kv, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func teardownKV(c *api.Client) error {
	_, err := c.KV().DeleteTree(defaultKeyPrefix, nil)
	return err
}
