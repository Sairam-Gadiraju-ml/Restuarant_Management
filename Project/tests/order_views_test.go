package tests

import (
	"Project/models"
	"Project/views"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// Mock data
var Orders = []models.Order{
	{Id: 1, CustomerId: 1, Amount: 100.0, Status: "Pending"},
	{Id: 2, CustomerId: 2, Amount: 200.0, Status: "Completed"},
}

func TestListCreateOrder(t *testing.T) {

	t.Run("GET all orders", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/orders", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&views.OrderServiceImplementation{}).ListCreateOrder)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Order details are , OrderID: 1 , CustomerId: 1 , Amount: 100, Status: Pending\nOrder details are , OrderID: 2 , CustomerId: 2 , Amount: 200, Status: Completed\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("POST create new order", func(t *testing.T) {
		newOrder := models.Order{CustomerId: 3, Amount: 300.0, Status: "Pending"}
		body, _ := json.Marshal(newOrder)
		req, err := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&views.OrderServiceImplementation{}).ListCreateOrder)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		expected := "Order Placed Successfully, Order ID: 3 \n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})
}

func TestHandleOrderID(t *testing.T) {

	t.Run("GET order by ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/orders/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&views.OrderServiceImplementation{}).HandleOrderID)
		req = mux.SetURLVars(req, map[string]string{"orderid": "1"})
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Order details: OrderID: 1, CustomerId: 1, Amount: 100, Status: Pending\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("PUT update order by ID", func(t *testing.T) {
		updatedOrder := models.Order{Status: "Completed", Amount: 150.0}
		body, _ := json.Marshal(updatedOrder)
		req, err := http.NewRequest(http.MethodPut, "/orders/1", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&views.OrderServiceImplementation{}).HandleOrderID)
		req = mux.SetURLVars(req, map[string]string{"orderid": "1"})
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Order ID 1 updated successfully.\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("DELETE order by ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "/orders/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&views.OrderServiceImplementation{}).HandleOrderID)
		req = mux.SetURLVars(req, map[string]string{"orderid": "1"})
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Order ID 1 deleted successfully.\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})
}
