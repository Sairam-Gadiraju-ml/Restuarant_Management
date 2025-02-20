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

func TestBookTable(t *testing.T) {
	views.IntializeTables() // Initialize tables before running tests

	t.Run("Book a table immediately", func(t *testing.T) {
		booking := models.Booking{
			Customer: models.Customer{FirstName: "Sai", LastName: "Ram"},
			WeekDay:  models.Monday,
			Time:     10,
			BookNow:  true,
		}
		body, _ := json.Marshal(booking)
		req, err := http.NewRequest(http.MethodPost, "/book", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&views.TableServiceImplementation{}).BookTable)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusAccepted {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
		}

		expected := "Booking Confirmed! Your Booking ID is BOOK-Monday-"
		if !bytes.Contains(rr.Body.Bytes(), []byte(expected)) {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("Book a table for later", func(t *testing.T) {
		booking := models.Booking{
			Customer: models.Customer{FirstName: "Sai", LastName: "Ram"},
			WeekDay:  models.Tuesday,
			Time:     11,
			BookNow:  false,
		}
		body, _ := json.Marshal(booking)
		req, err := http.NewRequest(http.MethodPost, "/book", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&views.TableServiceImplementation{}).BookTable)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "No table is available for immediate booking. You have been added to the queue for Tuesday at 11'O clock."
		if !bytes.Contains(rr.Body.Bytes(), []byte(expected)) {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})
}

func TestCancelTable(t *testing.T) {
	views.IntializeTables() // Initialize tables before running tests

	t.Run("Cancel a booked table", func(t *testing.T) {
		// First, book a table to cancel
		booking := models.Booking{
			Customer: models.Customer{FirstName: "Sai", LastName: "Ram"},
			WeekDay:  models.Monday,
			Time:     10,
			BookNow:  true,
		}
		body, _ := json.Marshal(booking)
		req, err := http.NewRequest(http.MethodPost, "/book", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&views.TableServiceImplementation{}).BookTable)
		handler.ServeHTTP(rr, req)

		// Now, cancel the booked table
		req, err = http.NewRequest(http.MethodDelete, "/cancel?weekday=Monday&time=10&tableid=1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr = httptest.NewRecorder()
		handler = http.HandlerFunc((&views.TableServiceImplementation{}).CancelTable)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Table 1 at Monday '10 has been canceled successfully."
		if !bytes.Contains(rr.Body.Bytes(), []byte(expected)) {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})
}

func TestAddTable(t *testing.T) {
	views.IntializeTables() // Initialize tables before running tests

	t.Run("Add tables to a specific weekday and hour", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/add?weekday=Monday&hour=10&numTables=2", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&views.TableServiceImplementation{}).AddTable)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "2 tables added at 10 on Monday"
		if !bytes.Contains(rr.Body.Bytes(), []byte(expected)) {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})
}

func TestRemoveTable(t *testing.T) {
	views.IntializeTables() // Initialize tables before running tests

	t.Run("Remove a specific table from a weekday and hour", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "/remove?weekday=Monday&hour=10&tableid=1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&views.TableServiceImplementation{}).RemoveTable)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Table 1 removed at Monday '10 successfully."
		if !bytes.Contains(rr.Body.Bytes(), []byte(expected)) {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})
}

func TestGetInfo(t *testing.T) {
	views.IntializeTables() // Initialize tables before running tests

	t.Run("Get all tables information", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/info", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&views.TableServiceImplementation{}).GetInfo)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if !bytes.Contains(rr.Body.Bytes(), []byte("tablesData")) {
			t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
		}
	})
}

func TestGetFreeTables(t *testing.T) {
	views.IntializeTables() // Initialize tables before running tests

	t.Run("Get available timings to book a table for a specific weekday", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/freetables/Monday", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&views.TableServiceImplementation{}).GetFreeTables)
		req = mux.SetURLVars(req, map[string]string{"weekday": "Monday"})
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if !bytes.Contains(rr.Body.Bytes(), []byte("Available Timings on Monday to book the table are")) {
			t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
		}
	})
}

func TestGetTables(t *testing.T) {
	views.IntializeTables() // Initialize tables before running tests

	t.Run("Get table availability for each day and hour", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/tables", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&views.TableServiceImplementation{}).GetTables)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if !bytes.Contains(rr.Body.Bytes(), []byte("Tables available at following hours")) {
			t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
		}
	})
}
