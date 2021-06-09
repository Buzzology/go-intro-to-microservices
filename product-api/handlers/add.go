package handlers

import (
	"net/http"
	"github.com/Buzzology/go-intro-to-microservices/product-api/data"
)


func (p *Products) AddProduct(rw http.ResponseWriter, req *http.Request) {
	product := req.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&product)
	p.l.Printf("Product added: %#v", product)
}