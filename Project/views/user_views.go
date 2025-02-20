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

var Users []models.User

// UserServiceImplementation implements the UserService interface.
type UserServiceImplementation struct {
}

// UserService defines the methods required for user management.
type UserService interface {
	// RegisterUser handles user registration.
	RegisterUser(w http.ResponseWriter, r *http.Request)

	// LoginUser handles user login.
	LoginUser(w http.ResponseWriter, r *http.Request)

	// UpdateUser allows updating user information.
	UpdateUser(w http.ResponseWriter, r *http.Request)

	// DeleteUser deletes a user from the system.
	DeleteUser(w http.ResponseWriter, r *http.Request)

	// ListUsers returns a list of all registered users.
	ListUsers(w http.ResponseWriter, r *http.Request)
}

// LoginUser allows a user to log in using their username and password.
// It decodes the request body, checks for any decoding errors
func (s *UserServiceImplementation) LoginUser(w http.ResponseWriter, r *http.Request) {
	var User models.User
	err := json.NewDecoder(r.Body).Decode(&User)
	if err != nil {
		log.Println("Error Decoding the Request Body")
		fmt.Fprintln(w, CustomError("Error Decoding the Request Body", 400))
		return
	}

	// Logic to check username and password against stored data can be added here
}

// RegisterUser handles the registration of a new user.
func (s *UserServiceImplementation) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var User models.User
	Id := len(Users) + 1
	err := json.NewDecoder(r.Body).Decode(&User)
	if err != nil {
		log.Println("Error Decoding the Request Body")
		fmt.Fprintln(w, CustomError("Error Decoding the Request Body", 400))
		return
	}
	User.Id = uint(Id)
	Users = append(Users, User)
	// Return a success response
	fmt.Fprintf(w, "User %v registered successfully.", User.Username)
	log.Printf("User %v registered successfully.", User.Username)
}

// UpdateUser handles the updating of user information.
// It decodes the request body and updates the user details provided in the request.
func (s *UserServiceImplementation) UpdateUser(w http.ResponseWriter, r *http.Request) {

	var User models.User
	id, iderr := strconv.Atoi(mux.Vars(r)["id"])
	if iderr != nil {
		log.Println("Provide Valid UserID")
		fmt.Fprintln(w, CustomError("Provide Valid UserID", 400))
		return
	}

	err := json.NewDecoder(r.Body).Decode(&User)
	if err != nil {
		log.Println("Error Decoding the Request Body")
		fmt.Fprintln(w, CustomError("Error Decoding the Request Body", 400))
		return
	}
	// Logic to update user details

	if id != 0 {
		for i, user := range Users {
			if user.Id == User.Id {
				if User.Username != "" {
					Users[i].Username = User.Username
				}
				if User.Role != "" {
					Users[i].Role = User.Role
				}
			}
		}
		fmt.Fprintf(w, "User %v updated successfully.", User.Username)
		log.Printf("User %v updated successfully.", User.Username)
	} else {
		fmt.Fprintf(w, "Provide Valid UserID")
		log.Printf("Provide Valid UserID")
	}

}

// DeleteUser handles the deletion of a user from the system.
func (s *UserServiceImplementation) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var User models.User
	err := json.NewDecoder(r.Body).Decode(&User)
	if err != nil {
		log.Println("Error Decoding the Request Body")
		fmt.Fprintln(w, CustomError("Error Decoding the Request Body", 400))
		return
	}

	for i, user := range Users {
		// Check if request body has Username
		if User.Username != "" {
			if user.Username == User.Username {
				Users = append(Users[:i], Users[i+1:]...)
				fmt.Fprintf(w, "User %v Deleted Successfully", user.Username)
				log.Printf("User %v Deleted Successfully", user.Username)
				return
			}
			// Check if request body has User ID
		} else if User.Id != 0 {
			if user.Id == User.Id {
				Users = append(Users[:i], Users[i+1:]...)
				fmt.Fprintf(w, "User %v Deleted Successfully", user.Username)
				log.Printf("User %v Deleted Successfully", user.Username)
				return
			}
		}
	}

	// If user not found
	fmt.Fprintf(w, "User %v not found.", User.Username)
	log.Printf("User %v not found.", User.Username)
}

// ListUsers returns the list of all registered users.
func (s *UserServiceImplementation) ListUsers(w http.ResponseWriter, r *http.Request) {
	if len(Users) > 0 {
		fmt.Fprintln(w, "User Details are :")
		log.Println("User Details are :")
		w.WriteHeader(http.StatusAccepted)
		for _, user := range Users {
			fmt.Fprintf(w, " Id: %v, Username: %v, Role: %v \n", user.Id, user.Username, user.Role)
			log.Printf(" Id: %v, Username: %v, Role: %v \n", user.Id, user.Username, user.Role)
		}
	}
}
