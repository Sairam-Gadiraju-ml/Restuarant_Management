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
var MenuItems = []models.MenuItem{
	{Id: 1, Name: "Pizza", Price: 100.0},
	{Id: 2, Name: "Burger", Price: 50.0},
}

func TestGetMenu(t *testing.T) {

	t.Run("GET all menu items", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/menu", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&views.MenuServiceImplementation{}).GetMenu)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Menu Id : 1 , Item Name: Pizza , Item Price : 100 Menu Id : 2 , Item Name: Burger , Item Price : 50"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("POST create new menu item", func(t *testing.T) {
		newMenuItem := models.MenuItem{Name: "Pasta", Price: 80.0}
		body, _ := json.Marshal(newMenuItem)
		req, err := http.NewRequest(http.MethodPost, "/menu", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&views.MenuServiceImplementation{}).GetMenu)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		expected := "Menu Item Created: \n Id: 3, Name: Pasta, Price: 80 \n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})
}

func TestHandleMenuItem(t *testing.T) {

	t.Run("GET menu item by ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/menu/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&views.MenuServiceImplementation{}).HandleMenuItem)
		req = mux.SetURLVars(req, map[string]string{"Id": "1"})
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Item Id: 1, Item Name: Pizza, Item Price: 100"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("PUT update menu item by ID", func(t *testing.T) {
		updatedItem := models.MenuItem{Name: "Pizza Updated", Price: 120.0}
		body, _ := json.Marshal(updatedItem)
		req, err := http.NewRequest(http.MethodPut, "/menu/1", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&views.MenuServiceImplementation{}).HandleMenuItem)
		req = mux.SetURLVars(req, map[string]string{"Id": "1"})
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Menu Item Updated: Id: 1, Name: Pizza Updated, Price: 120"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("DELETE menu item by ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "/menu/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&views.MenuServiceImplementation{}).HandleMenuItem)
		req = mux.SetURLVars(req, map[string]string{"Id": "1"})
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Menu Item Deleted: Id: 1"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})
}
