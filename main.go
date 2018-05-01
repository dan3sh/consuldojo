package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	s, err := c.Agent().Services()
	if err != nil {
		log.Fatalf("services: %v", err)
	}
	fmt.Println(s)

	kv := c.KV()
	got, _, err := kv.Get("app/k1", nil)
	if err != nil {
		log.Fatalf("get: %v", err)
	}
	fmt.Printf("got: %#s\n", got.Value)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
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
