// Package main
package main

import (
	"context"
	"errors"
	"time"

	hello "github.com/dh1tw/natsgreeter/srv/proto/hello"

	micro "github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/util/log"
	"google.golang.org/grpc"
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

	go func() {
		for {
			grpc.DialContext(context.TODO(), "127.0.0.1:9091")
			time.Sleep(time.Second)
		}
	}()

	service := micro.NewService(
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
