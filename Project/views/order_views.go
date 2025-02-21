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

var Orders []models.Order

// ListCreateOrder Lists all orders if it's a GET request or places a new order if it's a POST method.
func (s *OrderServiceImplementation) ListCreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// List all orders
		if len(Orders) == 0 {
			fmt.Fprintf(w, "No orders found.")
			return
		}
		for _, v := range Orders {
			log.Printf("Order details are , OrderID: %v , CustomerId: %v , Amount: %v, Status: %v",
				v.Id, v.CustomerId, v.Amount, v.Status)
			fmt.Fprintf(w, "Order details are , OrderID: %v , CustomerId: %v , Amount: %v, Status: %v\n",
				v.Id, v.CustomerId, v.Amount, v.Status)
		}
	} else if r.Method == http.MethodPost {
		// Place a new order
		var order models.Order
		err := json.NewDecoder(r.Body).Decode(&order)
		if err != nil {
			log.Println(CustomError("Error Decoding the Response Body", 400))
			fmt.Fprintf(w, "%v", CustomError("Error Decoding the Response Body", 400))
			return
		}
		order.Id = uint(len(Orders) + 1) // Create a new Order ID based on the current length
		Orders = append(Orders, models.Order{order.Id, order.CustomerId, order.Amount, order.Status})
		log.Printf("Order Placed Successfully, Order ID: %v \n", order.Id)
		fmt.Fprintf(w, "Order Placed Successfully, Order ID: %v \n", order.Id)
	}
}

// HandleOrderID handles operations on a specific order by its ID (GET, PUT, DELETE)
func (s *OrderServiceImplementation) HandleOrderID(w http.ResponseWriter, r *http.Request) {
	// Extract the OrderID from the URL path using Gorilla Mux path variable
	vars := mux.Vars(r)
	orderIDStr := vars["orderid"] // Extract the OrderID from the path parameter
	if orderIDStr == "" {
		fmt.Fprintf(w, "OrderID is required in the URL path.")
		return
	}
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		log.Println("Error parsing Order ID", err)
		fmt.Fprintf(w, "Invalid Order ID format.")
		return
	}

	// Look for the order with the given Order ID
	var orderIndex int
	var orderFound bool
	for i, order := range Orders {
		if order.Id == uint(orderID) {
			orderIndex = i
			orderFound = true
			break
		}
	}

	if !orderFound {
		fmt.Fprintf(w, "Order with ID %v not found.", orderID)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Fetch the order by ID
		order := Orders[orderIndex]
		log.Printf("Fetching Order details for OrderID: %v", order.Id)
		fmt.Fprintf(w, "Order details: OrderID: %v, CustomerId: %v, Amount: %v, Status: %v\n",
			order.Id, order.CustomerId, order.Amount, order.Status)

	case http.MethodPut:
		// Update the order by ID
		var updatedOrder models.Order
		err := json.NewDecoder(r.Body).Decode(&updatedOrder)
		if err != nil {
			log.Println("Error decoding the update data:", err)
			fmt.Fprintf(w, "Error decoding the update data.")
			return
		}

		// Update the order details (example: updating the status and amount)
		Orders[orderIndex].Status = updatedOrder.Status
		Orders[orderIndex].Amount = updatedOrder.Amount
		log.Printf("Updated Order details for OrderID: %v", Orders[orderIndex].Id)
		fmt.Fprintf(w, "Order ID %v updated successfully.\n", Orders[orderIndex].Id)

	case http.MethodDelete:
		// Delete the order by ID
		Orders = append(Orders[:orderIndex], Orders[orderIndex+1:]...) // Remove the order from the slice
		log.Printf("Deleted OrderID: %v", orderID)
		fmt.Fprintf(w, "Order ID %v deleted successfully.\n", orderID)

	default:
		fmt.Fprintf(w, "Method not allowed.")
	}
}


