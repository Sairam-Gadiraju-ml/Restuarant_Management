package views

import (
	"Project/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var Bills []models.Bill

func GetCreateBill(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// List all bills
		log.Println("Listing all Bills")
		for _, bill := range Bills {
			w.WriteHeader(http.StatusOK)
			log.Printf("Bill Id: %v , Customer Name: %v , Amount: %v", bill.Id, bill.Customer.FirstName, bill.Amount)
			fmt.Fprintf(w, "Bill Id: %v , Customer Name: %v , Amount: %v", bill.Id, bill.Customer.FirstName, bill.Amount)
		}
	case http.MethodPost:
		// Create a new bill
		var newBill models.Bill
		err := json.NewDecoder(r.Body).Decode(&newBill)
		if err != nil {
			log.Println("Error Decoding the Request Body")
			http.Error(w, "Error Decoding the Request Body", http.StatusBadRequest)
			return
		}

		// Assign an ID to the new bill
		newBill.Id = uint(len(Bills) + 1)
		Bills = append(Bills, newBill)

		log.Printf("Bill Created: Id: %v, Customer Name: %v, Amount: %v", newBill.Id, newBill.Customer.FirstName, newBill.Amount)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Bill Created: Id: %v, Customer Name: %v, Amount: %v", newBill.Id, newBill.Customer.FirstName, newBill.Amount)
	}
}

func HandleBillID(w http.ResponseWriter, r *http.Request) {
	// Extract the bill ID from the URL
	param := mux.Vars(r)["Id"]
	Id, err := strconv.Atoi(param)
	if err != nil {
		log.Println("Error Parsing ID from Request")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Find the bill by ID
	var billIndex int
	itemFound := false
	for i, bill := range Bills {
		if bill.Id == uint(Id) {
			billIndex = i
			itemFound = true
			break
		}
	}

	if !itemFound {
		log.Printf("Bill with ID %v not found", Id)
		http.Error(w, "Bill Not Found", http.StatusNotFound)
		return
	}

	// Handle HTTP methods
	switch r.Method {
	case http.MethodGet:
		// Get the details of a bill
		bill := Bills[billIndex]
		log.Printf("Bill Id: %v, Customer Name: %v, Amount: %v", bill.Id, bill.Customer.FirstName, bill.Amount)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Bill Id: %v, Customer Name: %v, Amount: %v", bill.Id, bill.Customer.FirstName, bill.Amount)

	case http.MethodPut:
		// Update the bill
		var updatedBill models.Bill
		err := json.NewDecoder(r.Body).Decode(&updatedBill)
		if err != nil {
			log.Println("Error Decoding the Request Body for Update")
			http.Error(w, "Error Decoding the Request Body", http.StatusBadRequest)
			return
		}

		// Update the bill details
		Bills[billIndex].Customer.FirstName = updatedBill.Customer.FirstName
		Bills[billIndex].Amount = updatedBill.Amount

		log.Printf("Bill Updated: Id: %v, Customer Name: %v, Amount: %v", updatedBill.Id, updatedBill.Customer.FirstName, updatedBill.Amount)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Bill Updated: Id: %v, Customer Name: %v, Amount: %v", updatedBill.Id, updatedBill.Customer.FirstName, updatedBill.Amount)

	case http.MethodDelete:
		// Delete the bill
		Bills = append(Bills[:billIndex], Bills[billIndex+1:]...)
		log.Printf("Bill Deleted: Id: %v", Id)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Bill Deleted: Id: %v", Id)
	}
}

type BillService interface {
	GetCreateBill(w http.ResponseWriter, r *http.Request)
	HandleBillID(w http.ResponseWriter, r *http.Request)
}

type BillServiceImplementation struct{}
