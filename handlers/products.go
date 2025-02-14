// Package classification of product API
// Documentation of Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta

package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"main.go/data"
)

// A list of products
// swagger:response productsResponse
type productResponse struct {
	body []data.Product
}

// swagger:parameters deleteProduct
type productIDParamaeterWrapper struct {
	// Id of the product to delete from database
	// in: path
	// required: true
	ID int `json:"id"`
}

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

// swagger:route GET /products products listProducts
// returns a list of products
// responses:
// 		200: productsResponse

func (p *Products) GetProducts(wtr http.ResponseWriter, res *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(wtr)
	if err != nil {
		http.Error(wtr, "Can Not marshal Data", http.StatusInternalServerError)
	}

}

func (p *Products) AddProduct(wtr http.ResponseWriter, res *http.Request) {
	prod := res.Context().Value(keyProduct{}).(*data.Product)
	data.AddProduct(prod)
	p.l.Printf("Prod %v", prod)

}

func (p *Products) UpdateProduct(wtr http.ResponseWriter, res *http.Request) {
	vars := mux.Vars(res)
	id, err1 := strconv.Atoi(vars["id"])
	if err1 != nil {
		http.Error(wtr, "ouldnt convert string", http.StatusBadRequest)
		return
	}
	prod := res.Context().Value(keyProduct{}).(*data.Product)

	err1 = data.UpdateProduct(id, prod)
	if err1 == data.ErrorProductNotFound {
		fmt.Println(err1)
		http.Error(wtr, "Could not update the resource comming from 1", http.StatusNotFound)
		return
	}
	if err1 != nil {

		http.Error(wtr, "Could not update the resource comming from 2", http.StatusInternalServerError)
		return
	}

}

// swagger:route DELETE /products/{id} products deleteProduct
// Deletes a product with the given ID
// responses:
// 	201: productResponse

// Delets a Product From the ProductList with the stipulated Id
func (p *Products) DeleteProduct(wtr http.ResponseWriter, res *http.Request) {
	vars := mux.Vars(res)
	id, _ := strconv.Atoi(vars["id"])

	err := data.DeleteProduct(id)
	if err == data.ErrorProductNotFound {
		http.Error(wtr, "Product Not found", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(wtr, "Product Not Found", http.StatusInternalServerError)
		return
	}
	wtr.Write([]byte("Product deleted successfully"))
}

type keyProduct struct{}

func (p *Products) MiddleWareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, res *http.Request) {
		prod := &data.Product{}
		err1 := prod.Validate()
		if err1 != nil {
			p.l.Println("Could not validate the product details we recieved!")
			http.Error(wtr, "the data we recieved is not in the appropriate format", http.StatusBadRequest)
		}
		err := prod.FromJSON(res.Body)
		if err != nil {
			http.Error(wtr, fmt.Sprintf("Unable to unmarshal data: %s", err), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(res.Context(), keyProduct{}, prod)
		res = res.WithContext(ctx)

		next.ServeHTTP(wtr, res)
	})
}
