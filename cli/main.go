package main

import (
	"context"
	"fmt"

	grpc "github.com/asim/go-micro/plugins/client/grpc/v3"
	micro "github.com/asim/go-micro/v3"
	hello "github.com/dh1tw/natsgreeter/srv/proto/hello"
)

func main() {

	service := micro.NewService(
		micro.Client(grpc.NewClient()),
	)

	service.Init()

	// Use the generated client stub
	cl := hello.NewSayService("go.micro.srv.greeter", service.Client())

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
