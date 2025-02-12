package data

import "testing"

func TestValidation(t *testing.T) {
	p := &Product{
		Name:  "kekw",
		Price: 0.01,
		SKU:   "abc-def-ghi",
	}

	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}

}
