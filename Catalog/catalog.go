package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"github.com/gorilla/mux"
)

type Product struct {
	Uuid string `json:"uuid"`
	Product string `json:"product"`
	Price float64 `json:"price,string"`
}

type Products struct {
	Products []Product `json:"products"`
}

var productURL string

func init() {
	productURL = os.Getenv("PRODUCT_URL")
}

func loadProducts() []Product {
	response, err := http.Get(productURL + "/products")

	if err != nil {
		fmt.Println("Error no HTTP")
	}

	data, _ := ioutil.ReadAll(response.Body)

	var products Products
	json.Unmarshal(data, &products)

	return products.Products
}

func main(){
	r := mux.NewRouter()
	r.HandleFunc("/", ListProducts)
	r.HandleFunc("/product/{id}", ShowProduct)
	http.ListenAndServe(":8080", r)
}

func ListProducts (w http.ResponseWriter, r *http.Request) {
	products := loadProducts()

	t := template.Must(template.ParseFiles("templates/catalog.html"))
	t.Execute(w, products)
}

func ShowProduct (w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response, err := http.Get(productURL + "/product/" + vars["id"])
	if err != nil {
		fmt.Println("Error HTTP")
	}

	data, _ := ioutil.ReadAll(response.Body)

	var product Product
	json.Unmarshal(data, &product)

	t := template.Must(template.ParseFiles("templates/view.html"))
	t.Execute(w, product)
}
