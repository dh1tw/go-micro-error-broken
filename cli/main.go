package main

import (
	"context"
	"fmt"

	hello "github.com/dh1tw/natsgreeter/srv/proto/hello"

	micro "github.com/asim/go-micro/v3"
)

func main() {

	service := micro.NewService()

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
