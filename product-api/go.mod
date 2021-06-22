module github.com/Buzzology/go-intro-to-microservices/product-api

go 1.16

replace github.com/Buzzology/go-intro-to-microservices/product-api => ./

replace github.com/Buzzology/go-intro-to-microservices/currency => ../currency

require (
	github.com/Buzzology/go-intro-to-microservices/currency v0.0.0-00010101000000-000000000000
	github.com/go-openapi/runtime v0.19.29
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-hclog v0.16.1
	github.com/leodido/go-urn v1.2.1 // indirect
	google.golang.org/grpc v1.38.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
)
