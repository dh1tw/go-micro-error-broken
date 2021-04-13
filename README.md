## THIS DEMO CODE USES NATS
# Update 

The problem has been solved. You have to instanciate your client with the Option ContentType "application/proto-rpc".
This will propagate the error from the server correctly to the client. 


```golang

    c := client.NewClient(
		client.PoolSize(1), // some custom options
		client.ContentType("application/proto-rpc"), // <<-----
	)

```

See the full explanation in issue [go-micro/2110](https://github.com/asim/go-micro/issues/2110)

## Description

I did some further digging in the git history and the Github issues. I found a few more clues that explain IMHO the current behavior: 
- As part of PR go-micro/#362 @asim introduced the [proto codec](https://github.com/asim/go-micro/tree/master/codec/proto). Apparently with the aim to allow go-micro to process standard json, protobuf and grpc requests. In this PR he changed the codec for ContentType `application/protobuf` from `protorpc.NewCodec` to `proto.NewCodec`. This didn't cause any problems, since the defaultContentType was `application/octet-stream` and still set to `protorpc.NewCodec`. 
- In PR go-micro/#372 the `DefaultContentType` was changed from `"application/octet-stream"` to `"application/protobuf"`. And from then onwards, `proto.Codec` was the default Codec.

Git blame revealed that the `func (c *Codec) Write (m *Codec.Message...)` in the `proto.Codec` actually always ignored the `m *codec.Message`. So I suppose that @asim did this on purpose and that the `proto.Codec` was never supposed to return an error to the client. A short confirmation from @asim would be very much appreciated. 

So if you want to propagate error messages from the server to the client, you have to use the `protorpc` codec.

One way to do so is by setting the contentType Option when calling `client.NewClient()`.  


# Greeter

An example Greeter application to demonstrate that [go-micro pull request #394](https://github.com/asim/go-micro/pull/396)
likely broke the error behaviour of go-micro rpc calls.

This demo application is based on the original [Greeter example](https://github.com/asim/go-micro/tree/master/examples/greeter)
supplied with go-micro.

The major difference are:
- extended the `Say` Service by one method called `Broken` (which demonstrates the bug)
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

The `Broken` handler just throws an error. The error message should be propagated to the client. However, it never
gets sent on the socket. Therefore the client will always just time out.

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

in [PR #396](https://github.com/asim/go-micro/pull/396/files), the [func (s *service) call(...) method](https://github.com/asim/go-micro/blob/bba3107ae13fb9ce9e273106c4543c5c50a460bc/server/rpc_router.go#L202) in `rpc_router.go` was modified. Basically if the
executed handler function results in an error, the code never has the change to reach the `router.sendResponse()` method.
I presume that therefore the error message never gets returned, and hence the client runs into the time-out.

Click [here for the lines in question](https://github.com/asim/go-micro/blob/bba3107ae13fb9ce9e273106c4543c5c50a460bc/server/rpc_router.go#L239-L245).