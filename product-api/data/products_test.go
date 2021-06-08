package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name: "test",
		Price: 0.01,
		SKU: "xxx-xxxx-xxxxx",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}