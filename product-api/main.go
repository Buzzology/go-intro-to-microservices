package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/Buzzology/go-intro-to-microservices/product-api/handlers"
)


func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	helloHandler := handlers.NewHello(l)
	goodbyeHandler := handlers.NewGoodbye(l)
	productHandler := handlers.NewProducts(l)

	sm := http.NewServeMux()
	sm.Handle("/products", productHandler)
	sm.Handle("/goodbye", goodbyeHandler)
	sm.Handle("/", helloHandler)

	s := http.Server{
		Addr: "127.0.0.1:9090",
		Handler: sm,
		IdleTimeout: 120*time.Second,
		ReadTimeout: 1 *time.Second,
		WriteTimeout: 1 *time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// Make a channel of os signal type
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <- sigChan // Read from signal channel (this appears to block until a message is received on the channel. Message is received when kill or cancel are invoked)
	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// Doesn't accept any new connections and waits for current ones to be handled before shutting down
	s.Shutdown(tc)
}