module github.com/dh1tw/natsgreeter/srv

go 1.16

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/asim/go-micro/cmd/protoc-gen-micro/v3 v3.0.0-20210217182006-0f0ace1a44a9 // indirect
	github.com/asim/go-micro/plugins/broker/nats/v3 v3.0.0-20210217182006-0f0ace1a44a9
	github.com/asim/go-micro/plugins/registry/nats/v3 v3.0.0-20210217182006-0f0ace1a44a9
	github.com/asim/go-micro/plugins/transport/nats/v3 v3.0.0-20210217182006-0f0ace1a44a9
	github.com/asim/go-micro/v3 v3.5.0
	github.com/golang/protobuf v1.4.2
	github.com/micro/go-micro/v2 v2.9.2-0.20201226154210-35d72660c801
	github.com/nats-io/nats.go v1.10.0
	google.golang.org/protobuf v1.23.0

)
