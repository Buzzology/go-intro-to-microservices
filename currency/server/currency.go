package server

import (
	"context"
	"fmt"
	"github.com/Buzzology/go-intro-to-microservices/currency/data"
	protos "github.com/Buzzology/go-intro-to-microservices/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
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

func NewCurrency(rates *data.ExchangeRates, l hclog.Logger) *Currency {
	return &Currency{l, protos.UnimplementedCurrencyServer{}, rates}
}
