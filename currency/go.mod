module github.com/Buzzology/go-intro-to-microservices/currency

go 1.16

replace github.com/Buzzology/go-intro-to-microservices/currency => ./

require (
	github.com/hashicorp/go-hclog v0.16.1
	google.golang.org/genproto v0.0.0-20210624195500-8bfb893ecb84
	google.golang.org/grpc v1.38.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0 // indirect
	google.golang.org/protobuf v1.26.0
)
