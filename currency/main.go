package main

import (
	"fmt"
	protos "github.com/Buzzology/go-intro-to-microservices/currency/protos/currency"
	"github.com/Buzzology/go-intro-to-microservices/currency/server"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"strconv"
)

func main() {
	var port = 9092
	log := hclog.Default()

	// Create a new gRPC server, use WithInsecure to allow http connections
	gs := grpc.NewServer()

	// Create an instance of the Currency server
	cs := server.NewCurrency(log)

	// Register the currency server
	protos.RegisterCurrencyServer(gs, cs)

	// Used for gRPCurl etc
	reflection.Register(gs)

	// Define port for grpc to list on
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	fmt.Printf("Running currency gRPC server on %v...", port)

	gs.Serve(listener)

}
