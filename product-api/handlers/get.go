package handlers

import (
	"github.com/Buzzology/go-intro-to-microservices/product-api/data"
	"net/http"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
// 200: productsResponse

// GetProducts returns the products from the data store
func (p *Products) GetProducts(rw http.ResponseWriter, h *http.Request) {

	listProducts := data.GetProducts()

	err := listProducts.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal products", http.StatusInternalServerError)
	}
}
