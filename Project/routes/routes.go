package routes

import (
	"Project/views"

	"github.com/gorilla/mux"
)

// InitializeRoutes sets up the routes for the application.
func InitializeRoutes(router *mux.Router) {
	// End points Responsible for booking tables
	var tableservice = views.TableServiceImplementation{}
	table := router.PathPrefix("/table").Subrouter()
	table.HandleFunc("/get", tableservice.GetTables).Methods("GET")
	table.HandleFunc("/book", tableservice.BookTable).Methods("POST")
	table.HandleFunc("/cancel", tableservice.CancelTable).Methods("PATCH")
	table.HandleFunc("/info", tableservice.GetInfo).Methods("GET")
	table.HandleFunc("/free/{weekday}", tableservice.GetFreeTables).Methods("GET")
	table.HandleFunc("/add", tableservice.AddTable).Methods("POST")
	table.HandleFunc("/remove", tableservice.RemoveTable).Methods("DELETE")

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

	// Endpoints Responsible for Menu
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
}
