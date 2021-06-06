package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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

	if req.Method == http.MethodPut {
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(req.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("Invalid uri: ", g)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("Invalid uri: more than one id")
			http.Error(rw, "Invalid uri: more than one id", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Invalid id: unable to convert to number", http.StatusBadRequest)
			return
		}

		p.updateProduct(id, rw, req)
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
		http.Error(rw, "Failed to unmarshal request body into product.", http.StatusBadRequest)
	}

	data.AddProduct(product)
	p.l.Printf("Product added: %#v", product)
}


func (p *Products) updateProduct(id int, rw http.ResponseWriter, req *http.Request) {
	product := &data.Product{}

	err := product.FromJSON(req.Body)
	if err != nil {
		http.Error(rw, "Failed to unmarshal request body into product.", http.StatusBadRequest)
	}

	if err = data.UpdateProduct(id, product); err != nil {
		if err == data.ErrProductNotFound {
			http.Error(rw, "Product not found.", http.StatusNotFound)
			return
		}

		http.Error(rw, "An error occurred while updating product.", http.StatusInternalServerError)
		return
	}

	p.l.Printf("Product updated: %#v", product)
}
