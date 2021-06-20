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

	// Needs to be more than one as we always populate EUR as 1
	if len(tr.rates) > 1 {
		t.Fatal("Expected at least one rate to be assigned.")
	}

	fmt.Sprintf("Rates: %#v", tr.rates)
}
