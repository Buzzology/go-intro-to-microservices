module github.com/Buzzology/go-intro-to-microservices/currency

go 1.16

replace github.com/Buzzology/go-intro-to-microservices/currency => ./

require (
	github.com/hashicorp/go-hclog v0.16.1
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.25.0 // indirect
)
