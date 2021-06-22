package handlers

import (
	"github.com/Buzzology/go-intro-to-microservices/product-api/data"
	"net/http"
)

func (p *Products) AddProduct(rw http.ResponseWriter, req *http.Request) {
	product := req.Context().Value(KeyProduct{}).(data.Product)
	p.productDB.AddProduct(&product)
	p.l.Debug("Product added: %#v", product)
}
