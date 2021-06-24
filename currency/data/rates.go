package data

import (
	"encoding/xml"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"math/rand"
	"net/http"
	"strconv"
	"time"
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

// MonitorRates checks the rates in the ECB API and sends messages to the returned channel when there are changes.
// NOTE: The ECB API only updates once a day
func (e *ExchangeRates) MonitorRates(interval time.Duration) chan struct{} {
	ret := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				// Just add a random difference to the rate and return it. This simulates the fluctuations in currency rates.
				for currency, rate := range e.rates {

					// Change can be 10% of original value
					change := (rand.Float64() / 10)

					// Is this a positive or negative change
					direction := rand.Intn(1)

					if direction == 0 {
						// New value will be min 90% of old
						change = 1 - change
					} else {
						// New value wil lbe 110% of old
						change = 1 + change
					}

					// Modify the rate
					e.rates[currency] = rate * change
				}

				// Notify updates, this will block unless there is a listener on the other end
				// This seems to send an "empty" value via the channel to inform the receiver that a change has occurred
				// and they should retrieve an updated list of rates
				ret <- struct{}{}
			}
		}
	}()

	return ret
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
