package data

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"testing"
)

func TestNewRates(t *testing.T) {
	tr, err := NewRates(hclog.Default())

	if err != nil {
		t.Fatal(err)
	}

	if len(tr.rates) == 0 {
		t.Fatal("Expected at least one rate to be assigned.")
	}

	fmt.Sprintf("Rates: %#v", tr.rates)
}
