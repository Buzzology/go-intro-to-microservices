package server

import (
	"context"
	"fmt"
	"github.com/Buzzology/go-intro-to-microservices/currency/data"
	protos "github.com/Buzzology/go-intro-to-microservices/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
	"io"
	"time"
)

type Currency struct {
	log hclog.Logger
	protos.UnimplementedCurrencyServer
	rates         *data.ExchangeRates
	subscriptions map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
}

func NewCurrency(rates *data.ExchangeRates, l hclog.Logger) *Currency {

	// Create the new instance of currency
	c := &Currency{
		log:                         l,
		UnimplementedCurrencyServer: protos.UnimplementedCurrencyServer{},
		rates:                       rates,
		subscriptions:               make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest),
	}

	// Start a new go routine to handle currency updates
	go c.handleUpdates()

	return c
}

// handleUpdates sends updated exchange rates to subscribed clients
func (c *Currency) handleUpdates() {
	rateUpdates := c.rates.MonitorRates(5 * time.Second)

	for range rateUpdates {

		c.log.Info("Got updated rates")

		// Loop over any subscribed clients
		for client, rateRequests := range c.subscriptions {

			// Loop over rates they're subscribed to
			for _, rateRequest := range rateRequests {

				// Retrieve the update rate
				rate, err := c.rates.GetRate(rateRequest.GetBase().String(), rateRequest.GetDestination().String())
				if err != nil {
					c.log.Error("Failed to get rate", "error", err)
				}

				// Stream the updated rate to the client
				err = client.Send(&protos.RateResponse{Rate: rate, Base: rateRequest.Base, Destination: rateRequest.Destination})
				if err != nil {
					c.log.Error("Failed to stream rate to client", "error", err)
				}
			}
		}
	}
}

// GetRate implements the CurrencyServer GetRate method and returns the currency excahnge rate for the two given currencies.
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		c.log.Error("failed to get rate", err)
		return nil, fmt.Errorf("failed to get rate %v", err)
	}

	return &protos.RateResponse{Rate: rate}, nil
}

// SubscribeRates streams exchange rates.
func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {

	for {
		rateRequest, err := src.Recv() // Blocks
		if err == io.EOF {
			c.log.Info("Client has closed connection")
			break
		}

		if err != nil {
			c.log.Error("Unable to read from client", "error", err)
			break
		}

		c.log.Info("Handle client request", "request", rateRequest)

		// Check if we already a map of requests
		rateRequests, ok := c.subscriptions[src]
		if !ok {

			// If not, initialise it
			rateRequests = []*protos.RateRequest{}
		}

		// Add the request to the map
		rateRequests = append(rateRequests, rateRequest)

		// Update the client's list of subscriptions
		c.subscriptions[src] = rateRequests
	}

	return nil
}
