// Package main
package main

import (
	"context"
	"errors"

	hello "github.com/dh1tw/natsgreeter/srv/proto/hello"

	grpc "github.com/asim/go-micro/plugins/server/grpc/v3"
	micro "github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/util/log"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Log("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

func (s *Say) Broken(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Log("Received Say.Broken request")
	// return microErr.New("simulating an error", "testing", 500)
	return errors.New("simulating an error")
}

func main() {

	service := micro.NewService(
		micro.Server(grpc.NewServer()),
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
