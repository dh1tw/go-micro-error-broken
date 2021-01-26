# Greeter

An example Greeter application to demonstrate that [go-micro pull request #394](https://github.com/asim/go-micro/pull/396)
likely broke the error behaviour of go-micro rpc calls when using transport plugins (at least the [nats transport plugin](https://github.com/asim/go-micro/tree/master/plugins/transport/nats)).

This demo application is based on the original [Greeter example](https://github.com/asim/go-micro/tree/master/examples/greeter)
supplied with go-micro.

The major difference are:
- using NATS as Transport, Broker, Plugin
- extended the `Say` Service by one method called `Broken` (which demonstrates the bug)
- Server / Client / Service are invoked explicitely instead of the MICRO_* environment variables (I have to use non-standard nats-options)
- The `svr` and `cli` are both indiviual go modules (and must therefore be executed directly from their module directory)

```protobuf
service Say {
	rpc Hello(Request) returns (Response) {}
	rpc Broken(Request) returns (Response) {}
}
```

In `srv/main.go` the two handlers are implemented:

```go
func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Log("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

func (s *Say) Broken(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	return errors.New("simulating an error")
}
```

The `Broken` handler just throws an error. The error message should be propagated to the client. However, it
never gets sent on the socket. Therefore the client will always just time out.

## Run Service & Client

assuming you cloned this repo into your GOPATH, run:

Start go.micro.srv.greeter
```shell
cd $GOPATH/src/github.com/dh1tw/natsgreeter/srv
go run main.go
```

Call go.micro.srv.greeter via client
```shell
cd $GOPATH/src/github.com/dh1tw/natsgreeter/cli
go run main.go
```

## What the client returns
```shell
Hello John
{"id":"go.micro.client","code":408,"detail":"call timeout: context deadline exceeded","status":"Request Timeout"}
```
## Expected output
```shell
Hello John
{"id":"go.micro.client","code":500,"detail":"simulating an error","status":"Error"} // or something similar
```

## Likely source of the problem

in [PR #394](https://github.com/asim/go-micro/pull/396/files), the [func (s *service) call(...) method](https://github.com/asim/go-micro/blob/bba3107ae13fb9ce9e273106c4543c5c50a460bc/server/rpc_router.go#L202) in `rpc_router.go` was modified. Basically if the
executed handler function results in an error, the code never has the change to reach the `router.sendResponse()` method.
I presume that therefore the error message never gets returned, and hence the client runs into the time-out.

Click [here for the lines in question](https://github.com/asim/go-micro/blob/bba3107ae13fb9ce9e273106c4543c5c50a460bc/server/rpc_router.go#L239-L245).