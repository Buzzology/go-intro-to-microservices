package data

import (
	"encoding/xml"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"strconv"
)

type ExchangeRates struct {
	log   hclog.Logger
	rates map[string]float64
}

// GetRate get the exchange rate between these two currencies
func (e *ExchangeRates) GetRate(base, dest string) (float64, error) {
	br, ok := e.rates[base]
	if ok != true {
		return 0, fmt.Errorf("rate no found for currency %v", base)
	}

	dr, ok := e.rates[dest]
	if ok != true {
		return 0, fmt.Errorf("rate no found for currency %v", dest)
	}

	return dr / br, nil
}

func NewRates(l hclog.Logger) (*ExchangeRates, error) {
	exchangeRates := &ExchangeRates{log: l, rates: map[string]float64{}}

	// Populate the rates map
	exchangeRates.getRates()

	return exchangeRates, nil
}

// getRates retrieves exchange rates from endpoint and assigns to rates map
func (e *ExchangeRates) getRates() error {

	// Retrieve the rate values as xml
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected error code 200 but receive %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	// Parse the xml using the structs defined below (note the xml definitions)
	md := &Cubes{}
	xml.NewDecoder(resp.Body).Decode(&md)

	// Loop through each of the retrieved rates and assign to the exchange rates map
	for _, cube := range md.CubeData {
		r, err := strconv.ParseFloat(cube.Rate, 64)
		if err != nil {
			return err
		}

		e.rates[cube.Currency] = r
	}

	e.rates["EUR"] = 1

	return nil
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
