package views

import (
	"Project/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// CustomErrorStruct represents an error structure with a message and a code.
type CustomErrorStruct struct {
	Message string
	Code    int
}

// Error method implements the error interface.
func (e *CustomErrorStruct) Error() string {
	return fmt.Sprintf("Message: %s (code: %d)", e.Message, e.Code)
}

// CustomError returns a new CustomErrorStruct instance.
func CustomError(message string, code int) *CustomErrorStruct {
	return &CustomErrorStruct{
		Message: message,
		Code:    code,
	}
}

// Intializing the Tables
var Tables_Data = make(map[string][]models.HourDetails)
var BookingQueue = make(map[string]map[int]Queue)

// IntializeTables initializes tables for a week.
func IntializeTables() {
	var new_table []models.Table
	for t := 1; t <= 10; t++ {
		// Initialize 10 tables for each daych day
		new_table = append(new_table, models.Table{ID: strconv.Itoa(t), Capacity: 4, IsEmpty: true})
	}

	// Initializing Empty Tables for Hours from 9-22-22
	var new_hours []models.HourDetails
	for h := 9; h <= 22; h++ {
		// For each hour from 9 to 22, initialize tablestables
		new_hours = append(new_hours, models.HourDetails{Hour: h, Table: new_table})
	}

	// Get current day and timeime
	currentTime := time.Now()
	currentDay := currentTime.Weekday()
	currentHour := currentTime.Hour()

	// Convert currentDay to our WeekDay enum (adjusting since Go Weekday is Sunday=0)
	var weekdayEnum models.WeekDay
	switch currentDay {
	case time.Monday:
		weekdayEnum = models.Monday
	case time.Tuesday:
		weekdayEnum = models.Tuesday
	case time.Wednesday:
		weekdayEnum = models.Wednesday
	case time.Thursday:
		weekdayEnum = models.Thursday
	case time.Friday:
		weekdayEnum = models.Friday
	case time.Saturday:
		weekdayEnum = models.Saturday
	case time.Sunday:
		weekdayEnum = models.Sunday
	}

	// Set up tables for today starting from the current hourour
	var todayHours []models.HourDetails
	for h := currentHour; h <= 22; h++ { // Loop from the current hour till 22 (10 PM)
		todayHours = append(todayHours, models.HourDetails{Hour: h, Table: new_table})
	}
	Tables_Data[models.WeekDayToString[weekdayEnum]] = todayHours

	// Add tables for the next 6 days after today
	for i := 1; i <= 6; i++ {
		// Calculate the next day using modulus operator for wrapping around the week
		nextDay := (weekdayEnum + models.WeekDay(i)) % 7
		Tables_Data[models.WeekDayToString[nextDay]] = new_hours
	}

	log.Println("Initializing Tables")
}

// TableService defines the interface for managing table-related operations.
// TableService defines the interface for managing table-related operations.
type TableService interface {
	// BookTable handles table booking requests.
	BookTable(w http.ResponseWriter, r *http.Request)

	// CancelTable handles table cancellation requests.
	CancelTable(w http.ResponseWriter, r *http.Request)

	// GetTables retrieves all tables.
	GetTables(w http.ResponseWriter, r *http.Request)

	// GetFreeTables retrieves all free tables.
	GetFreeTables(w http.ResponseWriter, r *http.Request)

	// GetInfo retrieves information about a specific table.
	GetInfo(w http.ResponseWriter, r *http.Request)

	// AddTable adds a new table to the system.
	AddTable(w http.ResponseWriter, r *http.Request)

	// RemoveTable removes a table from the system.
	RemoveTable(w http.ResponseWriter, r *http.Request)
}
type TableServiceImplementation struct{}

// BillService defines the interface for handling bill-related operations in the restaurant management system.
type BillService interface {
	// GetCreateBill handles the creation of a new bill.
	GetCreateBill(w http.ResponseWriter, r *http.Request)

	// HandleBillID handles operations related to a specific bill identified by its ID.
	HandleBillID(w http.ResponseWriter, r *http.Request)
}

type BillServiceImplementation struct{}

// MenuService defines the interface for menu-related operations in the restaurant management system.
type MenuService interface {
	// GetMenu retrieves the menu and writes it to the response.
	GetMenu(w http.ResponseWriter, r *http.Request)

	// HandleMenuItem handles operations related to individual menu items.
	HandleMenuItem(w http.ResponseWriter, r *http.Request)
}

type MenuServiceImplementation struct{}

// TableService defines the interface for managing table-related operations.
type OrderService interface {
	ListCreateOrder(w http.ResponseWriter, r *http.Request)
	HandleOrderID(w http.ResponseWriter, r *http.Request)
}

// TableServiceImplementation implements the TableService interface.
type OrderServiceImplementation struct{}

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
