package views

import (
	"Project/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var Customers []models.Customer

// GetAllCustomer handles the request to get all registered customers
func GetAllCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		log.Println("Listing all registered customers.")
		// If there are no customers, return a friendly message
		if len(Customers) == 0 {
			fmt.Fprintf(w, "No customers found!")
			return
		}

		// Respond with the customer details
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "The Customers Registered are: %v", Customers)
		log.Printf("The Customers Registered are: %v", Customers)
	} else {
		log.Println(CustomError("Method Not Allowed", http.StatusMethodNotAllowed))
	}
}

// HandleCustomer manages individual customer actions (GET, PUT, DELETE)
func HandleCustomer(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/customers/"):]
	Id, err := strconv.Atoi(param)
	if err != nil {
		log.Println("Error parsing customer ID")
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	var customerIndex int
	itemFound := false
	for i, customer := range Customers {
		if customer.Id == uint(Id) {
			customerIndex = i
			itemFound = true
			break
		}
	}

	if !itemFound {
		log.Printf("Customer with ID %d not found", Id)
		http.Error(w, "Customer Not Found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Fetch the customer details by ID
		customer := Customers[customerIndex]
		log.Printf("Fetching customer: %v", customer)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Customer ID: %v, Name: %v %v, Contact: %v, Type: %v", customer.Id, customer.FirstName, customer.LastName, customer.Contact, customer.Type)

	case http.MethodPut:
		// Update the customer details
		var updatedCustomer models.Customer
		err := json.NewDecoder(r.Body).Decode(&updatedCustomer)
		if err != nil {
			log.Println("Error decoding request body")
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}

		// Update customer fields
		Customers[customerIndex].FirstName = updatedCustomer.FirstName
		Customers[customerIndex].LastName = updatedCustomer.LastName
		Customers[customerIndex].Contact = updatedCustomer.Contact
		Customers[customerIndex].Type = updatedCustomer.Type

		log.Printf("Updated customer: %v", updatedCustomer)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successfully updated customer: %v %v", updatedCustomer.FirstName, updatedCustomer.LastName)

	case http.MethodDelete:
		// Delete the customer
		Customers = append(Customers[:customerIndex], Customers[customerIndex+1:]...)
		log.Printf("Deleted customer with ID: %d", Id)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successfully deleted customer with ID: %d", Id)

	default:
		log.Println(CustomError("Method Not Allowed", http.StatusMethodNotAllowed))
	}
}

// AddCustomer adds a new customer to the list
func AddCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Decode the request body into a new customer object
		var customer models.Customer
		err := json.NewDecoder(r.Body).Decode(&customer)
		if err != nil {
			log.Println("Error decoding the request")
			http.Error(w, "Error decoding the request", http.StatusBadRequest)
			return
		}

		// Assign a new ID to the customer
		customer.Id = uint(len(Customers) + 1)

		// Append the customer to the list
		Customers = append(Customers, customer)

		log.Printf("Successfully added customer: %v", customer)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Successfully added the customer %v %v", customer.FirstName, customer.LastName)
	} else {
		log.Println("Method Not Allowed")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
