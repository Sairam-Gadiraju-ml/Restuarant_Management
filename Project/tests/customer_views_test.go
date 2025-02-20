package tests

import (
	"Project/models"
	"Project/views"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var Customers []models.Customer

func TestHandleCustomer(t *testing.T) {
	// Mock data
	Customers = []models.Customer{
		{Id: 1, FirstName: "Sai", LastName: "Ram", Contact: "1234567890", Type: "Regular"},
		{Id: 2, FirstName: "Ram", LastName: "Varma", Contact: "0987654321", Type: "VIP"},
	}

	t.Run("GET existing customer", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/customer/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(views.HandleCustomer)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Customer ID: 1, Name: Sai Ram, Contact: 1234567890, Type: Regular"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("GET non-existing customer", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/customer/3", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(views.HandleCustomer)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	t.Run("PUT update customer", func(t *testing.T) {
		updatedCustomer := models.Customer{FirstName: "Sai", LastName: "Ram", Contact: "1111111111", Type: "Regular"}
		body, _ := json.Marshal(updatedCustomer)
		req, err := http.NewRequest(http.MethodPut, "/customer/1", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(views.HandleCustomer)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Successfully updated customer: Sai Ram"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("DELETE customer", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "/customer/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(views.HandleCustomer)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Successfully deleted customer with ID: 1"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("POST add customer", func(t *testing.T) {
		newCustomer := models.Customer{FirstName: "Varma", LastName: "Sai", Contact: "2222222222", Type: "Regular"}
		body, _ := json.Marshal(newCustomer)
		req, err := http.NewRequest(http.MethodPost, "/customer/add", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(views.AddCustomer)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		expected := "Successfully added the customer Varma Sai"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

}
