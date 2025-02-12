package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"Description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)

}

// Custom function to validate sku
func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[a-z]+-[a-z]+-[a-z]+$`)
	return re.MatchString(fl.Field().String())
}

func (p *Product) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(p)
}

func (p *Products) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(p)
}

func GetProducts() Products {
	return productList
}

func AddProduct(prod *Product) {
	prod.ID = getId()
	productList = append(productList, prod)
	fmt.Println(productList)

}

func getId() int {
	lp := productList[len(productList)-1]
	cur := lp.ID
	return cur + 1
}

func UpdateProduct(id int, prod *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return ErrorProductNotFound
	}

	prod.ID = id
	productList[pos] = prod
	displayProducts()
	return nil
}

func displayProducts() {
	for _, val := range productList {
		fmt.Println(val)
	}
}

var ErrorProductNotFound = fmt.Errorf("could not find the required product in the database")

func findProduct(id int) (*Product, int, error) {
	for i, val := range productList {
		if val.ID == id {
			return val, i, nil
		}
	}
	return nil, -1, ErrorProductNotFound
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy Milk Coffe",
		Price:       2.45,
		SKU:         "abc123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "short and strong coffe without milk",
		Price:       1.99,
		SKU:         "gri143",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
