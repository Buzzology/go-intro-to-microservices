package data

import (
	"context"
	"encoding/json"
	"fmt"
	protos "github.com/Buzzology/go-intro-to-microservices/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

type ProductsDB struct {
	currency protos.CurrencyClient
	log      hclog.Logger
	rates    map[string]float64
	client   protos.Currency_SubscribeRatesClient
}

func NewProductsDB(c protos.CurrencyClient, l hclog.Logger) *ProductsDB {

	// Create a new instance of products db
	productsDb := &ProductsDB{c, l, make(map[string]float64), nil}

	// Subscribe to exchange rate changes so that our local cache doesn't go stale
	go productsDb.handleUpdates()

	return productsDb
}

// handleUpdates subscribes to the currency server and handles updating the exchange rate in our local cache
func (p *ProductsDB) handleUpdates() {

	// Setup a subscription to rate changes
	sub, err := p.currency.SubscribeRates(context.Background())
	if err != nil {
		p.log.Error("Failed to subscribe to rates", "error", err)
		return
	}

	// Add the client to the db
	p.client = sub

	for {

		// Handle changes received as rate responses. NOTE: receive is blocking
		rateResponse, err := sub.Recv()
		if err != nil {
			p.log.Error("Failed to retrieve rate", "error", err)
			return
		}

		// Check if the response is a streamed error type (the one we created using oneof)
		if grpcError := rateResponse.GetError(); grpcError != nil {
			p.log.Error("Error subscribing for rates", "error", grpcError)
			continue
		}

		// Check if it's a rate response (successful message)
		if resp := rateResponse.GetRateResponse(); resp != nil {

			p.log.Info(fmt.Sprintf("Receive updated rate from currency stream: %v:%v %v", resp.Base, resp.Destination, resp.Rate))

			// Update the exchange rate in the local cache
			p.rates[resp.GetDestination().String()] = resp.Rate
			continue
		}

		p.log.Error("Unexpected message type received")
	}
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

		// Check if we have access to detailed error
		if status, ok := status.FromError(err); ok {

			// Extract the returned rate request
			messageDetails := status.Details()[0].(*protos.RateRequest)

			if status.Code() == codes.InvalidArgument {
				return *product, fmt.Errorf("unable to get rate from currency server. Destination and base currencies cannot be the same")
			}

			return *product, fmt.Errorf("unable to get rate from currency server. Base: %s Destination %s", messageDetails.Base, messageDetails.Destination)
		}

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

	return pr, nil
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

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *ProductsDB) AddProduct(prod *Product) {
	prod.ID = getNextID()
	productList = append(productList, prod)
}

func (p *ProductsDB) UpdateProduct(id int, prod *Product) error {
	_, index, err := findProduct(id)
	if err != nil {
		return err
	}

	prod.ID = id
	productList[index] = prod
	return nil
}

func (p *ProductsDB) DeleteProduct(id int) error {
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

	// If it's in our local cache just return that
	if _, ok := p.rates[destination]; ok {
		return p.rates[destination], nil
	}

	// Prepare rate request
	rr := protos.RateRequest{
		Base:        protos.Currencies_EUR,
		Destination: protos.Currencies(protos.Currencies_value[destination]),
	}

	// Get initial exchange rate via gRPC
	resp, err := p.currency.GetRate(context.Background(), &rr)
	if err != nil {
		p.log.Error("unable to get rate", "destination", destination, "error", err)
		return 0, err
	}

	// Set the initial rate
	p.rates[destination] = resp.Rate

	// Subscribe to future rates
	err = p.client.Send(&rr)
	if err != nil {
		p.log.Error("failed to subscribe to future updates", "error", err)
	}

	return resp.Rate, nil
}

// UpdateProduct replaces a product in the database with the given
// item.
// If a product with the given id does not exist in the database
// this function returns a ProductNotFound error
//func (p *ProductsDB) UpdateProduct(pr Product) error {
//	i := findIndexByProductID(pr.ID)
//	if i == -1 {
//		return ErrProductNotFound
//	}
//
//	// update the product in the DB
//	productList[i] = &pr
//
//	return nil
//}

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}
