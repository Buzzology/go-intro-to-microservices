package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/Buzzology/go-intro-to-microservices/product-api/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, h *http.Request) {
	listProducts := data.GetProducts()
	err := listProducts.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal products", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, req *http.Request) {
	product := req.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&product)
	p.l.Printf("Product added: %#v", product)
}


func (p *Products) UpdateProduct(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
		return
	}	
	
	product := req.Context().Value(KeyProduct{}).(data.Product)
	if err = data.UpdateProduct(id, &product); err != nil {
		if err == data.ErrProductNotFound {
			http.Error(rw, "Product not found.", http.StatusNotFound)
			return
		}

		http.Error(rw, "An error occurred while updating product.", http.StatusInternalServerError)
		return
	}

	p.l.Printf("Product updated: %#v", product)
}


type KeyProduct struct {}


func (p Products) MiddlewareProductionValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func (rw http.ResponseWriter, req *http.Request) {
		product := &data.Product{}
		err := product.FromJSON(req.Body)
		if err != nil {
			http.Error(rw, "Failed to unmarshal request body into product.", http.StatusBadRequest)
			return // If it fails to validate we terminate the handler chain
		}

		// We assign product to the context using a struct. You can use a string but I think
		// this is done to prevent it being overwritten (or overwriting) entries made elsewhere.
		// Will need to confirm (i think this was said in the golang podcast).
		ctx := context.WithValue(req.Context(), KeyProduct{}, *product)
		req = req.WithContext(ctx)
		
		next.ServeHTTP(rw, req)
	})
}