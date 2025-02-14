package views

import (
	"Project/models"
	"fmt"
	"log"
	"net/http"
)

var Customers []models.Customer

func GetAllCustomer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "The Customers Registered are: %v", Customers)
	w.WriteHeader(http.StatusOK)
	log.Printf("The Customers Registered are: %v", Customers)
}

func HandleCustomer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Println("GET")
	case http.MethodPut:
		fmt.Println("PUT")
	case http.MethodDelete:
		fmt.Println("Delete")
	default:
		log.Println(CustomError("Method Not Allowed", 404))
		fmt.Fprintln(w, CustomError("Method Not Allowed", 404))
	}
}

func AddCustomer(w http.ResponseWriter, r *http.Request) {

}
