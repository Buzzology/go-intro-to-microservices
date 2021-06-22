// Package classification of Product API
//
// Documentation for Product API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"github.com/Buzzology/go-intro-to-microservices/product-api/data"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"strconv"
)

// A list of products returned in the response (note this is only used for doco)
// swagger:response productsResponse
type productsResponseWrapper struct {
	// All products in the system
	// in: body
	Body []data.Product
}

// swagger:parameters deleteProduct
type productIdParameterWrapper struct {
	// The id of the product to delete from the database
	// in: path
	// required: true
	ID int `json:"id"`
}

// swagger:response NoContent
type productsNoContent struct{}

type Products struct {
	l         hclog.Logger
	productDB *data.ProductsDB
}

func NewProducts(l hclog.Logger, productDB *data.ProductsDB) *Products {
	return &Products{l, productDB}
}

type KeyProduct struct{}

func (p Products) MiddlewareProductionValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		product := &data.Product{}
		err := product.FromJSON(req.Body)
		if err != nil {
			p.l.Debug("[ERROR] unmarshalling product", err)
			http.Error(rw, "Failed to unmarshal request body into product.", http.StatusBadRequest)
			return // If it fails to validate we terminate the handler chain
		}

		err = product.Validate()
		if err != nil {
			p.l.Debug("[ERROR] validating product", err)
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

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// getProductID returns the product ID from the URL
// Panics if cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func getProductID(r *http.Request) int {
	// parse the product id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}
