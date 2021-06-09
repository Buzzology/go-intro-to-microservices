// Package classification Product API
//
//     Documentation for Product API
//     Schemes: http
//     BasePath: /
//     Version 1.0.0
//     
//     Consumes: 
//     - application/json
// swagger:meta
package handlers

import (
	"context"
	"fmt"
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


type KeyProduct struct {}


func (p Products) MiddlewareProductionValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func (rw http.ResponseWriter, req *http.Request) {
		product := &data.Product{}
		err := product.FromJSON(req.Body)
		if err != nil {
			p.l.Println("[ERROR] unmarshalling product", err)
			http.Error(rw, "Failed to unmarshal request body into product.", http.StatusBadRequest)
			return // If it fails to validate we terminate the handler chain
		}

		err = product.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(rw, fmt.Sprintf("Failed to validate product. %s", err), http.StatusBadRequest)
			return
		}

		// We assign product to the context using a struct. You can use a string but I think
		// this is done to prevent it being overwritten (or overwriting) entries made elsewhere.
		// Will need to confirm (i think this was said in the golang podcast).
		ctx := context.WithValue(req.Context(), KeyProduct{}, *product)
		req = req.WithContext(ctx)
		
		next.ServeHTTP(rw, req)
	})
}