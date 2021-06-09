package handlers

import (
	"net/http"
	"github.com/Buzzology/go-intro-to-microservices/product-api/data"
)


func (p *Products) GetProducts(rw http.ResponseWriter, h *http.Request) {
	listProducts := data.GetProducts()
	err := listProducts.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal products", http.StatusInternalServerError)
	}
}