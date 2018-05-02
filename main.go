package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
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

	const (
		serviceID   = "serviceID"
		serviceName = "serviceName"
	)

	reg := &api.AgentServiceRegistration{
		ID:   serviceID,
		Name: serviceName,
		//Port: 4242,
	}
	err = c.Agent().ServiceRegister(reg)
	if err != nil {
		log.Fatalf("register: %v", err)
	}
	defer func() {
		err = c.Agent().ServiceDeregister(serviceID)
		if err != nil {
			log.Fatalf("unregister: %v", err)
		}
	}()

	ss, err := c.Agent().Services()
	if err != nil {
		log.Fatalf("services: %v", err)
	}
	for _, s := range ss {
		fmt.Printf("%#v\n", s)
	}

	handleFunc := func(a uint64, b interface{}) {
		fmt.Printf("%#v\n", b)
	}
	params := map[string]interface{}{
		"type":   "keyprefix",
		"prefix": defaultKeyPrefix,
	}

	plan, err := watch.Parse(params)
	if err != nil {
		fmt.Printf("%#v\n", plan)
		log.Fatalf("watch parse: %v", err)
	}
	plan.Handler = handleFunc
	err = plan.Run("")
	if err != nil {
		log.Fatalf("watch plan run: %v", err)
	}
	defer plan.Stop()

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
