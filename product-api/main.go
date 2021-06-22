package main

import (
	"context"
	protos "github.com/Buzzology/go-intro-to-microservices/currency/protos/currency"
	"github.com/Buzzology/go-intro-to-microservices/product-api/data"
	"github.com/Buzzology/go-intro-to-microservices/product-api/handlers"
	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	l := hclog.Default()

	// Create gRPC client
	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	currencyClient := protos.NewCurrencyClient(conn)

	sm := mux.NewRouter()

	// Instantiate the product db
	db := data.NewProductsDB(currencyClient, l)

	// Instantiate the product handlers
	productHandler := handlers.NewProducts(l, db)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", productHandler.ListAll).Queries("currency", "{[A-Z]{3}}")
	getRouter.HandleFunc("/products", productHandler.ListAll)
	getRouter.HandleFunc("/products/{id:[0-9]}", productHandler.ListSingle).Queries("currency", "{[A-Z]{3}}")
	getRouter.HandleFunc("/products/{id:[0-9]}", productHandler.ListSingle)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", productHandler.Update)
	putRouter.Use(productHandler.MiddlewareProductionValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", productHandler.AddProduct)
	postRouter.Use(productHandler.MiddlewareProductionValidation)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", productHandler.Delete)

	// Add swagger
	opt := middleware.RedocOpts{SpecURL: "/swagger.yaml"} // Points to generated swagger
	sh := middleware.Redoc(opt, nil)
	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	s := http.Server{
		Addr:         "127.0.0.1:9090",
		Handler:      ch(sm),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}),
	}

	go func() {

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Failed to start", "err", err)
		}
	}()

	l.Info("Starting Product API...")

	// Make a channel of os signal type
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan // Read from signal channel (this appears to block until a message is received on the channel. Message is received when kill or cancel are invoked)
	l.Info("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// Doesn't accept any new connections and waits for current ones to be handled before shutting down
	err = s.Shutdown(tc)
	if err != nil {
		l.Error("Failed to handle error", "error", err)
	}
}
