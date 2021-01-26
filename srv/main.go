// Package main
package main

import (
	"context"
	"errors"

	hello "greeter/srv/proto/hello"

	natsBroker "github.com/micro/go-micro/plugins/broker/nats/v2"
	natsReg "github.com/micro/go-micro/plugins/registry/nats/v2"
	natsTr "github.com/micro/go-micro/plugins/transport/nats/v2"

	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/util/log"
	nats "github.com/nats-io/nats.go"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Log("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

func (s *Say) Broken(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	return errors.New("simulating an error")
}

func main() {

	nopts := nats.GetDefaultOptions()
	nopts.Servers = []string{"127.0.0.1"}

	reg := natsReg.NewRegistry(natsReg.Options(nopts))
	br := natsBroker.NewBroker(natsBroker.Options(nopts))
	tr := natsTr.NewTransport(natsTr.Options(nopts))

	svr := server.NewServer(
		server.Registry(reg),
		server.Broker(br),
		server.Transport(tr),
	)

	service := micro.NewService(
		micro.Server(svr),
		micro.Registry(reg),
		micro.Broker(br),
		micro.Transport(tr),
		micro.Name("go.micro.srv.greeter"),
	)

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	hello.RegisterSayHandler(service.Server(), new(Say))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
