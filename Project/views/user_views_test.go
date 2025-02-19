package views

import (
	"Project/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteUser(t *testing.T) {
	// Mock data
	Users = []models.User{
		{Id: 1, Username: "sai_ram", Role: "admin"},
		{Id: 2, Username: "ram_varma", Role: "user"},
	}

	t.Run("DELETE existing user by username", func(t *testing.T) {
		userToDelete := models.User{Username: "sai_ram"}
		body, _ := json.Marshal(userToDelete)
		req, err := http.NewRequest(http.MethodDelete, "/user", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&UserServiceImplementation{}).DeleteUser)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "User sai_ram Deleted Successfully"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("DELETE non-existing user by username", func(t *testing.T) {
		userToDelete := models.User{Username: "non_existent_user"}
		body, _ := json.Marshal(userToDelete)
		req, err := http.NewRequest(http.MethodDelete, "/user", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&UserServiceImplementation{}).DeleteUser)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "User non_existent_user not found."
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("DELETE existing user by ID", func(t *testing.T) {
		userToDelete := models.User{Id: 2}
		body, _ := json.Marshal(userToDelete)
		req, err := http.NewRequest(http.MethodDelete, "/user", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc((&UserServiceImplementation{}).DeleteUser)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "User ram_varma Deleted Successfully"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})
}
