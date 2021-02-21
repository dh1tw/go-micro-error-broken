module github.com/dh1tw/natsgreeter/cli

go 1.16

require (
	github.com/asim/go-micro/plugins/broker/nats/v3 v3.0.0-20210217182006-0f0ace1a44a9
	github.com/asim/go-micro/plugins/registry/nats/v3 v3.0.0-20210217182006-0f0ace1a44a9
	github.com/asim/go-micro/plugins/transport/nats/v3 v3.0.0-20210217182006-0f0ace1a44a9
	github.com/asim/go-micro/v3 v3.5.0
	github.com/dh1tw/natsgreeter/srv v0.0.0-20210126005402-daa464050ae5
	github.com/nats-io/nats.go v1.10.0
	google.golang.org/protobuf v1.25.0 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
