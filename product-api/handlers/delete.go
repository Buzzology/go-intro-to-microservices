package handlers

import (
	"github.com/Buzzology/go-intro-to-microservices/product-api/data"
	"net/http"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Returns a list of products
// responses:
// 201: NoContent

// DeleteProduct deletes a product from the database
func (p *Products) Delete(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	id := getProductID(req)

	p.l.Debug("Deleting record", "id", id)

	err := p.productDB.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Error("Unable to delete record id does not exist")

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		p.l.Error("Unable to delete record", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
