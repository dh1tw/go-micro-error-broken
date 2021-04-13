package main

import (
	"context"
	"fmt"
	"time"

	"github.com/asim/go-micro/v3"
	hello "github.com/dh1tw/natsgreeter/srv/proto/hello"

	natsBroker "github.com/asim/go-micro/plugins/broker/nats/v3"
	natsReg "github.com/asim/go-micro/plugins/registry/nats/v3"
	natsTr "github.com/asim/go-micro/plugins/transport/nats/v3"

	client "github.com/asim/go-micro/v3/client"
	nats "github.com/nats-io/nats.go"
)

func main() {

	nopts := nats.GetDefaultOptions()
	nopts.Servers = []string{"127.0.0.1"}

	reg := natsReg.NewRegistry(natsReg.Options(nopts))
	br := natsBroker.NewBroker(natsBroker.Options(nopts))
	tr := natsTr.NewTransport(natsTr.Options(nopts))

	cli := client.NewClient(
		client.Broker(br),
		client.Transport(tr),
		client.Registry(reg),
		client.PoolSize(1),
		client.PoolTTL(time.Hour*8760),              // one year - don't TTL our connection
		client.ContentType("application/proto-rpc"), // <<-- must be used to propagate errors from server to client
	)

	service := micro.NewService(
		micro.Registry(reg),
		micro.Broker(br),
		micro.Transport(tr),
		micro.Client(cli),
	)

	service.Init()

	// Use the generated client stub
	cl := hello.NewSayService("go.micro.srv.greeter", cli)

	// Make request
	rsp, err := cl.Hello(context.Background(), &hello.Request{
		Name: "John",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp.Msg)

	// Make request --> This should return the error message
	// "simulating an error"
	newRsp, err := cl.Broken(context.Background(), &hello.Request{
		Name: "John",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(newRsp.Msg)
}
