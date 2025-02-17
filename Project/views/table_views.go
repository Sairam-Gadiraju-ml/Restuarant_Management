package views

import (
	"Project/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// CustomErrorStruct represents an error structure with a message and a code.
type CustomErrorStruct struct {
	Message string
	Code    int
}

var BookingQueue = make(map[models.WeekDay]map[int][]models.QueueEntry)

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
var Tables_Data = make(map[models.WeekDay][]models.HourDetails)

// IntializeTables initializes tables for a week.
func IntializeTables() {
	var new_table []models.Table
	for t := 1; t <= 10; t++ {
		// Initialize 10 tables for each day
		new_table = append(new_table, models.Table{ID: string(t), IsEmpty: true})
	}

	// Initializing Empty Tables for Hours from 9-22
	var new_hours []models.HourDetails
	for h := 9; h <= 22; h++ {
		// For each hour from 9 to 22, initialize tables
		new_hours = append(new_hours, models.HourDetails{Hour: h, Table: new_table})
	}

	// Get current day and time
	currentTime := time.Now()
	currentDay := currentTime.Weekday()
	currentHour := currentTime.Hour()

	// Set up tables for today starting from the current hour
	var todayHours []models.HourDetails
	for h := currentHour; h <= 22; h++ { // Loop from the current hour till 22 (10 PM)
		todayHours = append(todayHours, models.HourDetails{Hour: h, Table: new_table})
	}
	Tables_Data[models.WeekDay(currentDay)] = todayHours

	// Add tables for the next 6 days after today
	for i := 1; i <= 6; i++ {
		// Calculate the next day using modulus operator for wrapping around the week
		nextDay := (models.WeekDay(currentDay) + models.WeekDay(i)) % 7
		Tables_Data[nextDay] = new_hours
	}

	log.Println("Initializing Tables")
}

// ProcessQueue processes the queue of  requests.
func ProcessQueue(weekday models.WeekDay, hour int) {
	// If there are any entries in the queue for the given weekday and hour
	if len(BookingQueue[weekday][hour]) > 0 {
		// Find the first customer in the queue
		customer := BookingQueue[weekday][hour][0]
		// Attempt to book a table for this customer
		for _, v := range Tables_Data[weekday] {
			if v.Hour == hour {
				for _, tab := range v.Table {
					if tab.IsEmpty {
						tab.IsEmpty = false
						log.Printf("Booking table for customer %v at %v '%v", customer.CustomerName, weekday, hour)
						BookingQueue[weekday][hour] = BookingQueue[weekday][hour][1:] // Remove the processed entry from the queue
						// Notify the customer that their booking was successful
						fmt.Println("Table booked for customer:", customer.CustomerName)
						return
					}
				}
			}
		}
	}
}

// Periodically check the queue
func QueueProcessor() {
	for {
		time.Sleep(10 * time.Second) // Delay before checking the queue again
		for weekday, hourMap := range BookingQueue {
			for hour, queue := range hourMap {
				// Process the queue for each day and hour
				if len(queue) > 0 {
					ProcessQueue(weekday, hour)
				}
			}
		}
	}
}

// Book Table checks if a table is free and books it
// If no table is free for the requested time , will add it to the Queue
func (s *TableServiceImplementation) BookTable(w http.ResponseWriter, r *http.Request) {
	var b models.Booking
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		log.Panicln("Error Decoding")
		fmt.Println(CustomError("Error Decoding", 400))
	}

	log.Printf("Booking Table for %v at %v'O clock with BookNow: %v\n", b.Customer.FirstName+" "+b.Customer.LastName, b.Time, b.BookNow)

	// If it's a "Book Now" request
	if b.BookNow {
		isTableBooked := false
		// Check for an available table at the requested time
		for _, v := range Tables_Data[b.WeekDay] {
			if v.Hour == b.Time {
				for _, tab := range v.Table {
					if tab.IsEmpty {
						// Book the table
						tab.IsEmpty = false
						log.Printf("Booked Table %v on %v at %v'O clock for %v\n", tab.ID, b.WeekDay, b.Time, b.Customer.FirstName+" "+b.Customer.LastName)
						w.WriteHeader(http.StatusAccepted)
						fmt.Fprintf(w, "Booked Table %v on %v at %v'O clock for %v\n", tab.ID, b.WeekDay, b.Time, b.Customer.FirstName+" "+b.Customer.LastName)
						isTableBooked = true
						return
					}
				}
			}
		}

		// If no table is available for "Book Now", return an error
		if !isTableBooked {
			fmt.Fprintf(w, "No table is available for immediate booking at %v '%v\n", b.WeekDay, b.Time)
			log.Printf("No table available for immediate booking at %v '%v\n", b.WeekDay, b.Time)
		}
	} else {
		// If it's a "Book Later" request, add it to the queue
		log.Printf("No immediate availability, adding to queue for %v at %v'O clock\n", b.WeekDay, b.Time)

		if BookingQueue[b.WeekDay] == nil {
			BookingQueue[b.WeekDay] = make(map[int][]models.QueueEntry)
		}

		BookingQueue[b.WeekDay][b.Time] = append(BookingQueue[b.WeekDay][b.Time], models.QueueEntry{
			CustomerName: b.Customer.FirstName + " " + b.Customer.LastName,
			WeekDay:      b.WeekDay,
			Time:         b.Time,
		})

		fmt.Fprintf(w, "No table is available for immediate booking, you have been added to the queue for %v at %v'O clock.\n", b.WeekDay, b.Time)
	}
}

// CancelTable cancels a previously booked table and Process the Queue
func (s *TableServiceImplementation) CancelTable(w http.ResponseWriter, r *http.Request) {
	// Function to cancel the table, takes time and table id as query parameters
	weekday := r.URL.Query().Get("weekday")
	time := r.URL.Query().Get("time")
	tableid := r.URL.Query().Get("tableid")

	// Convert string params to appropriate data types
	weekDayInt, err := strconv.Atoi(weekday)
	if err != nil {
		log.Println("Error parsing weekday")
		fmt.Println(CustomError("Error parsing weekday", 400))
		return
	}

	timeInt, err := strconv.Atoi(time)
	if err != nil {
		log.Println("Error parsing time")
		fmt.Println(CustomError("Error parsing time", 400))
		return
	}

	// Find the table and set IsEmpty to true
	for _, v := range Tables_Data[models.WeekDay(weekDayInt)] {
		if v.Hour == timeInt {
			for _, tab := range v.Table {
				if tab.ID == tableid {
					tab.IsEmpty = true
					fmt.Fprintf(w, "Table %v at %v '%v has been canceled successfully.\n", tab.ID, models.WeekDay(weekDayInt), v.Hour)
					log.Printf("Table %v at %v '%v canceled successfully.\n", tab.ID, models.WeekDay(weekDayInt), v.Hour)

				}
			}
		}
		// After canceling a table, process the queue
		ProcessQueue(models.WeekDay(weekDayInt), timeInt)
		return
	}

	// If no table was found
	fmt.Fprintf(w, "Table not found.\n")
	log.Println("Table not found.")
}

// AddTable adds new tables to a specific weekday and hour.
func (s *TableServiceImplementation) AddTable(w http.ResponseWriter, r *http.Request) {
	// Parse the query parameters for weekday, hour, and the number of tables to add
	weekday := r.URL.Query().Get("weekday")
	hour := r.URL.Query().Get("hour")
	numTables := r.URL.Query().Get("numTables")

	// Convert weekday and hour to integers
	weekDayInt, err := strconv.Atoi(weekday)
	if err != nil {
		log.Println("Error parsing weekday")
		fmt.Println(CustomError("Error parsing weekday", 400))
		return
	}

	hourInt, err := strconv.Atoi(hour)
	if err != nil {
		log.Println("Error parsing hour")
		fmt.Println(CustomError("Error parsing hour", 400))
		return
	}

	numTablesInt, err := strconv.Atoi(numTables)
	if err != nil {
		log.Println("Error parsing number of tables")
		fmt.Println(CustomError("Error parsing number of tables", 400))
		return
	}

	// Adding the specified number of tables for the given weekday and hour
	for i := 0; i < numTablesInt; i++ {
		// Check if the hour exists for the given weekday
		for idx, v := range Tables_Data[models.WeekDay(weekDayInt)] {
			if v.Hour == hourInt {
				// Add a new table to the existing list of tables for this hour
				newTableId := string(len(v.Table) + 1)
				v.Table = append(v.Table, models.Table{ID: newTableId, IsEmpty: true})
				Tables_Data[models.WeekDay(weekDayInt)][idx].Table = v.Table
				log.Printf("Added Table %v at %v '%v\n", newTableId, models.WeekDay(weekDayInt), hourInt)
				break
			}
		}
	}

	fmt.Fprintf(w, "%d tables added at %v on %v", numTablesInt, hourInt, models.WeekDay(weekDayInt))
	log.Printf("%d tables added at %v on %v", numTablesInt, hourInt, models.WeekDay(weekDayInt))
}

// RemoveTable removes a specific table from a weekday and hour.
func (s *TableServiceImplementation) RemoveTable(w http.ResponseWriter, r *http.Request) {
	// Parse the query parameters for weekday, hour, and table ID to remove
	weekday := r.URL.Query().Get("weekday")
	hour := r.URL.Query().Get("hour")
	tableid := r.URL.Query().Get("tableid")

	// Convert weekday, hour, and table ID to integers
	weekDayInt, err := strconv.Atoi(weekday)
	if err != nil {
		log.Println("Error parsing weekday")
		fmt.Println(CustomError("Error parsing weekday", 400))
		return
	}

	hourInt, err := strconv.Atoi(hour)
	if err != nil {
		log.Println("Error parsing hour")
		fmt.Println(CustomError("Error parsing hour", 400))
		return
	}

	// Remove the specified table from the given weekday and hour
	for dayIdx, v := range Tables_Data[models.WeekDay(weekDayInt)] {
		if v.Hour == hourInt {
			// Loop through the tables and remove the specified table by ID
			for tableIdx, tab := range v.Table {
				if tab.ID == tableid {
					// Remove the table from the slice
					Tables_Data[models.WeekDay(weekDayInt)][dayIdx].Table = append(v.Table[:tableIdx], v.Table[tableIdx+1:]...)
					fmt.Fprintf(w, "Table %v removed at %v '%v successfully.\n", tableid, models.WeekDay(weekDayInt), hourInt)
					log.Printf("Table %v removed at %v '%v successfully.\n", tableid, models.WeekDay(weekDayInt), hourInt)
					return
				}
			}
		}
	}

	// If no table was found for removal
	fmt.Fprintf(w, "Table not found.\n")
	log.Println("Table not found.")
}

// GetInfo provides the entire tables' data.
func (s *TableServiceImplementation) GetInfo(w http.ResponseWriter, r *http.Request) {
	log.Printf("Printing Tables Data %v", Tables_Data)
	fmt.Fprintf(w, "Tables Data: %v", Tables_Data)
}

// GetFreeTables returns the available timings to book a table for a specific weekday.
func (s *TableServiceImplementation) GetFreeTables(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)["weekday"]
	weekday, err := strconv.Atoi(param)
	if err != nil {
		log.Println("Error Getting the Path Param")
		fmt.Println(CustomError("Error Getting the path param", 400))
		return
	}
	empty_tables := []int{}
	for _, tables := range Tables_Data[models.WeekDay(weekday)] {
		for _, table := range tables.Table {
			if table.IsEmpty {
				empty_tables = append(empty_tables, tables.Hour)
				break
			}
		}
	}
	fmt.Fprintf(w, "Available Timings on %v to book the table are %v", models.WeekDay(weekday), empty_tables)
	log.Printf("Available Timings on %v to book the table are %v", models.WeekDay(weekday), empty_tables)
}

// GetTables retrieves the table availability for each day and hour.
func (s *TableServiceImplementation) GetTables(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting Tables Info")
	var empty_info = make(map[models.WeekDay][]int)
	for days, hours := range Tables_Data {
		empty_tables := []int{}
		for _, tables := range hours {
			for _, table := range tables.Table {
				if table.IsEmpty {
					empty_tables = append(empty_tables, tables.Hour)
					break
				}
			}
		}
		empty_info[days] = empty_tables
	}

	log.Printf("Tables available at following hours: %v \n", empty_info)
	fmt.Fprintf(w, "Tables available at following hours: %v \n", empty_info)
}

// TableService defines the interface for managing table-related operations.
type TableService interface {
	BookTable(w http.ResponseWriter, r *http.Request)
	CancelTable(w http.ResponseWriter, r *http.Request)
	GetTables(w http.ResponseWriter, r *http.Request)
	GetFreeTables(w http.ResponseWriter, r *http.Request)
	GetInfo(w http.ResponseWriter, r *http.Request)
	AddTable(w http.ResponseWriter, r *http.Request)
	RemoveTable(w http.ResponseWriter, r *http.Request)
}

// TableServiceImplementation implements the TableService interface.
type TableServiceImplementation struct{}
