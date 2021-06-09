package handlers

import (
	"net/http"
	"strconv"

	"github.com/Buzzology/go-intro-to-microservices/product-api/data"
	"github.com/gorilla/mux"
)


func (p *Products) DeleteProduct(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
		return
	}
	
	if err = data.DeleteProduct(id); err != nil {
		if err == data.ErrProductNotFound {
			http.Error(rw, "Product not found.", http.StatusNotFound)
			return
		}

		http.Error(rw, "An error occurred while updating product.", http.StatusInternalServerError)
		return
	}

	p.l.Printf("Product deleted: %s", id)
}