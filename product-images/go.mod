module github.com/Buzzology/go-intro-to-microservices/product-images

go 1.16

replace github.com/Buzzology/go-intro-to-microservices/product-images => ./

require (
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-hclog v0.16.1
	github.com/nicholasjackson/env v0.6.0
	github.com/stretchr/testify v1.2.2
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
)
