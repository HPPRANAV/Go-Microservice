package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"main.go/data"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

// MAIN HANDLER FUNCTION
func (p *Products) ServeHTTP(wtr http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		p.getProducts(wtr)
		return
	}
	if req.Method == http.MethodPost {
		p.addProduct(wtr, req)
		return
	}
	if req.Method == http.MethodPut {

		regex := regexp.MustCompile(`/([0-9]+)`)
		path := req.URL.Path
		group := regex.FindAllStringSubmatch(path, -1)

		if len(group) != 1 {
			http.Error(wtr, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(group[0]) != 2 {
			http.Error(wtr, "Invalid URI", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(group[0][1])
		if err != nil {
			http.Error(wtr, "Could not convert into integer", http.StatusBadRequest)
			return
		}
		p.updateProduct(id, wtr, req)
		return

	}

	wtr.WriteHeader(http.StatusMethodNotAllowed)

}

func (p *Products) getProducts(wtr http.ResponseWriter) {
	lp := data.GetProducts()
	err := lp.ToJSON(wtr)
	if err != nil {
		http.Error(wtr, "Can Not marshal Data", http.StatusInternalServerError)
	}

}

func (p *Products) addProduct(wtr http.ResponseWriter, res *http.Request) {
	prod := &data.Product{}

	err := prod.FromJSON(res.Body)
	if err != nil {
		http.Error(wtr, "Unable to unmarshal data", http.StatusBadRequest)
	}

	data.AddProduct(prod)

	p.l.Printf("Prod %v", prod)

}

func (p *Products) updateProduct(id int, wtr http.ResponseWriter, res *http.Request) {
	prod := &data.Product{}

	err := prod.FromJSON(res.Body)
	if err != nil {
		http.Error(wtr, "Unable to unmarshal data", http.StatusBadRequest)
	}
	err1 := data.UpdateProduct(id, prod)
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
