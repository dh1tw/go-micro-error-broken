module github.com/dh1tw/natsgreeter/srv

go 1.16

// replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/golang/protobuf v1.4.2
	github.com/micro/go-micro/plugins/broker/nats/v2 v2.0.0-20210120135431-d94936f6c97c
	github.com/micro/go-micro/plugins/registry/nats/v2 v2.0.0-20210120135431-d94936f6c97c
	github.com/micro/go-micro/plugins/transport/nats/v2 v2.0.0-20210105173217-bf4ab679e18b
	github.com/micro/go-micro/v2 v2.9.2-0.20201226154210-35d72660c801
	github.com/nats-io/nats.go v1.10.0
	google.golang.org/protobuf v1.23.0

)
