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
	r *data.ExchangeRates
}

// GetRate implements the CurrencyServer GetRate method and returns the currency excahnge rate for the two given currencies.
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

	rate, err := c.r.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		c.log.Error("failed to get rate", err)
		return nil, fmt.Errorf("failed to get rate %v", err)
	}

	return &protos.RateResponse{Rate: rate}, nil
}

// SubscribeRates streams exchange rates.
func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {

	go func() {
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
		}
	}()

	for {
		err := src.Send(&protos.RateResponse{Rate: 12.1})
		if err != nil {
			return err
		}

		time.Sleep(5 * time.Second)
	}
}

func NewCurrency(rates *data.ExchangeRates, l hclog.Logger) *Currency {
	return &Currency{l, protos.UnimplementedCurrencyServer{}, rates}
}
