package main

import (
	"Project/views"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	// Creates a logFile.txt if not created already and Opens it to write logs
	file := views.CreateLogFile()

	// Intialize the tables for the day
	views.IntializeTables()
	go views.QueueProcessor()

	// End points Responsbile for booking tables
	var tableserice = views.TableServiceImplementation{}
	table := router.PathPrefix("/table").Subrouter()
	table.HandleFunc("/get", tableserice.GetTables).Methods("GET")
	table.HandleFunc("/book", tableserice.BookTable).Methods("POST")
	table.HandleFunc("/cancel", tableserice.CancelTable).Methods("PATCH")
	table.HandleFunc("/info", tableserice.GetInfo).Methods("GET")
	table.HandleFunc("/free/{weekday}", tableserice.GetFreeTables).Methods("GET")
	table.HandleFunc("/add", tableserice.AddTable).Methods("POST")
	table.HandleFunc("/remove", tableserice.RemoveTable).Methods("DELETE")

	// End points Responsible for Customer handling

	customer := router.PathPrefix("/customer").Subrouter()
	customer.HandleFunc("/", views.GetAllCustomer).Methods("GET")
	customer.HandleFunc("/{id}", views.HandleCustomer).Methods("GET", "PUT", "DELETE")
	customer.HandleFunc("/add", views.AddCustomer).Methods("POST")

	// End points Responsible for Orders
	orderservice := views.OrderServiceImplementation{}
	orders := router.PathPrefix("/orders").Subrouter()
	orders.HandleFunc("/", orderservice.ListCreateOrder).Methods("GET", "POST")
	orders.HandleFunc("/{id}", orderservice.HandleOrderID).Methods("GET", "PUT", "DELETE")

	//Endpoins Responsible for Menu
	menuservice := views.MenuServiceImplementation{}
	menu := router.PathPrefix("/menu").Subrouter()
	menu.HandleFunc("/", menuservice.GetMenu).Methods("GET", "POST")
	menu.HandleFunc("/{id}", menuservice.HandleMenuItem).Methods("GET", "PUT", "DELETE")

	// Endpoints for Login and Register
	userservice := views.UserServiceImplementation{}
	user := router.PathPrefix("/user").Subrouter()
	user.HandleFunc("/get", userservice.ListUsers).Methods("GET")
	user.HandleFunc("/login", userservice.LoginUser).Methods("POST")
	user.HandleFunc("/register", userservice.RegisterUser).Methods("POST")
	user.HandleFunc("/update/{id}", userservice.UpdateUser).Methods("PUT")
	user.HandleFunc("/delete/{id}", userservice.DeleteUser).Methods("DELETE")

	// Endpoints Responsible for Bills
	bills := router.PathPrefix("/bills").Subrouter()
	bills.HandleFunc("/", views.GetCreateBill).Methods("GET", "POST")
	bills.HandleFunc("/{id}", views.HandleBillID).Methods("GET", "PUT", "DELETE")

	log.Println("Intializing the server at http://localhost:5000/:")
	if err := http.ListenAndServe(":5000", router); err != nil {
		log.Panicln("Error Starting server")
		fmt.Println("Error starting server:", err)
	}
	// Closes the logFile.txt
	defer file.Close()
}
