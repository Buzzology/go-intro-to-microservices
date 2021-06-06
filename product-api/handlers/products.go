package handlers

import (
	"log"
	"net/http"
	"github.com/Buzzology/go-intro-to-microservices/product-api/data"
)

type Products struct {
	l *log.Logger
}


func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		p.GetProducts(rw, req)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}


func (p *Products) GetProducts(rw http.ResponseWriter, h *http.Request) {
	listProducts := data.GetProducts()
	err := listProducts.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal products", http.StatusInternalServerError)
	}
}
