package handlers

import (
	"context"
	protos "github.com/Buzzology/go-intro-to-microservices/currency/protos/currency"
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

	// Get exchange rate via gRPC
	rr := protos.RateRequest{
		Base:        protos.Currencies_EUR,
		Destination: protos.Currencies_GBP,
	}
	resp, err := p.cc.GetRate(context.Background(), &rr)
	if err != nil {
		p.l.Println("[ERROR] error getting new rate", err)
		return
	}

	// Update the price based on exchange rate
	for _, prod := range listProducts {
		prod.Price = prod.Price * resp.Rate
	}
}
