package data

import (
	"context"
	"encoding/json"
	"fmt"
	protos "github.com/Buzzology/go-intro-to-microservices/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

type ProductsDB struct {
	currency protos.CurrencyClient
	log      hclog.Logger
}

func NewProductsDB(c protos.CurrencyClient, l hclog.Logger) *ProductsDB {
	return &ProductsDB{c, l}
}

func (p *ProductsDB) GetProductByID(id int, currency string) (Product, error) {

	// Retrieve the product
	product, _, err := findProduct(id)
	if err != nil {
		return *product, err
	}

	// No need to manipulate prices
	if currency == "" {
		return *product, nil
	}

	// Retrieve the exchange rate
	rate, err := p.getRate(currency)
	if err != nil {
		return *product, err
	}

	// De-reference the product so that we're not manipulating the data directly
	newProduct := *product
	newProduct.Price = newProduct.Price * rate

	return newProduct, nil
}

func (p *ProductsDB) GetProducts(currency string) (Products, error) {

	if currency == "" {
		return productList, nil
	}

	// Retrieve the exchange rate
	rate, err := p.getRate(currency)
	if err != nil {
		return nil, err
	}

	// Update the price of products based on exchange rate
	pr := Products{}
	for _, prod := range productList {
		np := *prod
		np.Price = np.Price * rate
		pr = append(pr, &np)
	}

	return productList, nil
}

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for this product
	//
	// required: true
	// min: 1
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)

	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {

	// SKU is of format xxx-xxxx-xxxxx
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	// We should find exactly one match for it to be valid
	if len(matches) != 1 {
		return false
	}

	return true
}

type Products []*Product

func GetProducts() Products {
	return productList
}

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	_, index, err := findProduct(id)
	if err != nil {
		return err
	}

	p.ID = id
	productList[index] = p
	return nil
}

func DeleteProduct(id int) error {
	_, index, err := findProduct(id)
	if err != nil {
		return err
	}

	productList = append(productList[:index], productList[:index+1]...)
	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, 0, ErrProductNotFound
}

func getNextID() int {
	lastProduct := productList[len(productList)-1]
	return lastProduct.ID + 1
}

var productList = Products{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

func (p *ProductsDB) getRate(destination string) (float64, error) {

	// Get exchange rate via gRPC
	rr := protos.RateRequest{
		Base:        protos.Currencies_EUR,
		Destination: protos.Currencies_GBP,
	}

	resp, err := p.currency.GetRate(context.Background(), &rr)
	if err != nil {
		p.log.Error("unable to get rate", "destination", destination, "error", err)
		return 0, err
	}

	return resp.Rate, nil
}
