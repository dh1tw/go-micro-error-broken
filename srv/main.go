// Package main
package main

import (
	"context"
	"errors"

	hello "github.com/dh1tw/natsgreeter/srv/proto/hello"

	natsBroker "github.com/asim/go-micro/plugins/broker/nats/v3"
	natsReg "github.com/asim/go-micro/plugins/registry/nats/v3"
	natsTr "github.com/asim/go-micro/plugins/transport/nats/v3"

	micro "github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/server"
	"github.com/asim/go-micro/v3/util/log"
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
