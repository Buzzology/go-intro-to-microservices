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

	if req.Method == http.MethodPost {
		p.addProduct(rw, req)
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


func (p *Products) addProduct(rw http.ResponseWriter, req *http.Request) {
	product := &data.Product{}
	
	err := product.FromJSON(req.Body)
	if err != nil {
		http.Error(rw, "Failed to unmarshal request body into products.", http.StatusBadRequest)
	}

	data.AddProduct(product)
	p.l.Printf("Product added: %#v", product)
}