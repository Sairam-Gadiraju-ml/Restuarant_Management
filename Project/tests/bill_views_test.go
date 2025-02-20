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
var Bills = []models.Bill{
	{Id: 1, Customer: models.Customer{FirstName: "Ram"}, Amount: 100.0},
	{Id: 2, Customer: models.Customer{FirstName: "Varma"}, Amount: 200.0},
}

func TestGetCreateBill(t *testing.T) {

	t.Run("GET all bills", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/bills", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(views.GetCreateBill)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Bill Id: 1 , Customer Name: Ram , Amount: 100Bill Id: 2 , Customer Name: Varma , Amount: 200"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("POST create new bill", func(t *testing.T) {
		newBill := models.Bill{Customer: models.Customer{FirstName: "Sai"}, Amount: 300.0}
		body, _ := json.Marshal(newBill)
		req, err := http.NewRequest(http.MethodPost, "/bills", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(views.GetCreateBill)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		expected := "Bill Created: Id: 3, Customer Name: Sai, Amount: 300"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})
}

func TestHandleBillID(t *testing.T) {

	t.Run("GET bill by ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/bills/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(views.HandleBillID)
		req = mux.SetURLVars(req, map[string]string{"Id": "1"})
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Bill Id: 1, Customer Name: Ram, Amount: 100"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("PUT update bill by ID", func(t *testing.T) {
		updatedBill := models.Bill{Customer: models.Customer{FirstName: "Ram Updated"}, Amount: 150.0}
		body, _ := json.Marshal(updatedBill)
		req, err := http.NewRequest(http.MethodPut, "/bills/1", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(views.HandleBillID)
		req = mux.SetURLVars(req, map[string]string{"Id": "1"})
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Bill Updated: Id: 1, Customer Name: Ram Updated, Amount: 150"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("DELETE bill by ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "/bills/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(views.HandleBillID)
		req = mux.SetURLVars(req, map[string]string{"Id": "1"})
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Bill Deleted: Id: 1"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})
}
