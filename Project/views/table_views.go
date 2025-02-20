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

var BookingQueue = make(map[string]map[int][]models.QueueEntry)

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

// IntializeTables initializes tables for a week.
func IntializeTables() {
	var new_table []models.Table
	for t := 1; t <= 10; t++ {
		// Initialize 10 tables for each day
		new_table = append(new_table, models.Table{ID: strconv.Itoa(t), Capacity: 4, IsEmpty: true})
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

	// Set up tables for today starting from the current hour
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

// ProcessQueue processes the queue of  requests.
func ProcessQueue(weekday models.WeekDay, hour int) {
	// Convert WeekDay enum to string for accessing Tables_Data and BookingQueue
	weekdayString := models.WeekDayToString[weekday]

	// If there are any entries in the queue for the given weekday and hour
	if len(BookingQueue[weekdayString][hour]) > 0 {
		// Find the first customer in the queue
		customer := BookingQueue[weekdayString][hour][0]

		// Attempt to book a table for this customer
		for _, v := range Tables_Data[weekdayString] {
			if v.Hour == hour {
				for j, tab := range v.Table {
					if tab.IsEmpty {
						// Mark the table as booked
						tab.IsEmpty = false

						// Update the table in the data structure
						Tables_Data[weekdayString][hour-9].Table[j] = tab // Adjust for index (assuming hours start at 9)

						log.Printf("Booking table for customer %v at %v '%v", customer.CustomerName, weekdayString, hour)

						// Remove the processed entry from the queue
						BookingQueue[weekdayString][hour] = BookingQueue[weekdayString][hour][1:]

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
		for weekdayString, hourMap := range BookingQueue {
			// Iterate through each hour map for the current weekday
			for hour, queue := range hourMap {
				// Process the queue for each day and hour if there are any customers
				if len(queue) > 0 {
					// Convert weekdayString back to models.WeekDay for processing if needed
					weekday := models.StringToWeekDay[weekdayString]

					// Process the queue for this weekday and hour
					ProcessQueue(weekday, hour)
				}
			}
		}
	}
}

// BookTable checks if a table is free and books it
// If no table is free for the requested time, it will add it to the Queue
func (s *TableServiceImplementation) BookTable(w http.ResponseWriter, r *http.Request) {
	var b models.Booking
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		log.Println("Error Decoding")
		http.Error(w, "Error Decoding", http.StatusBadRequest)
		return
	}

	log.Printf("Booking Table for %v at %v'O clock with BookNow: %v\n", b.Customer.FirstName+" "+b.Customer.LastName, b.Time, b.BookNow)

	// If it's a Book Now request
	if b.BookNow {
		isTableBooked := false

		// Convert WeekDay enum to string for accessing Tables_Data
		weekdayString := models.WeekDayToString[b.WeekDay]

		// Check if Tables_Data for the requested weekday is initialized
		if Tables_Data[weekdayString] == nil || len(Tables_Data[weekdayString]) == 0 {
			http.Error(w, "No available tables", http.StatusInternalServerError)
			return
		}

		for _, v := range Tables_Data[weekdayString] {
			if v.Hour == b.Time {
				for j, tab := range v.Table {
					// If a table is available, book it
					if tab.IsEmpty {
						// Create a new instance of Table and copy the values
						bookedTable := models.Table{
							ID:       tab.ID,
							Capacity: tab.Capacity,
							IsEmpty:  false, // Mark as booked
						}

						// Update the table in the data structure with the modified (booked) table
						Tables_Data[weekdayString][b.Time-9].Table[j] = bookedTable

						// Assign the booked table to the booking
						b.Table = bookedTable
						b.BookingStatus = "Confirmed"

						// Generate a Booking ID (for simplicity, using a timestamp + a random string)
						b.ID = fmt.Sprintf("BOOK-%v-%v", b.WeekDay, time.Now().UnixNano())

						// Debugging: log the successful booking
						log.Printf("Booked Table %v on %v at %v'O clock for %v\n", b.Table.ID, weekdayString, b.Time, b.Customer.FirstName+" "+b.Customer.LastName)

						w.WriteHeader(http.StatusAccepted)
						fmt.Fprintf(w, "Booking Confirmed! Your Booking ID is %v. Table %v on %v at %v'O clock for %v\n", b.ID, b.Table.ID, weekdayString, b.Time, b.Customer.FirstName+" "+b.Customer.LastName)

						isTableBooked = true
						// Break out of the loop to ensure only one booking is made
						return
					}
				}
			}
		}

		// If no table was booked, handle that case (optional)
		if !isTableBooked {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "No available tables for booking.\n")
		}

		// If no table is available for "Book Now", add to queue
		if !isTableBooked {
			// Ensure the BookingQueue map is initialized
			if BookingQueue[weekdayString] == nil {
				BookingQueue[weekdayString] = make(map[int][]models.QueueEntry)
			}

			// Ensure the slice for the requested time is initialized
			if BookingQueue[weekdayString][b.Time] == nil {
				BookingQueue[weekdayString][b.Time] = make([]models.QueueEntry, 0)
			}

			// Add the customer to the queue and provide a queue number
			queueNum := len(BookingQueue[weekdayString][b.Time]) + 1
			BookingQueue[weekdayString][b.Time] = append(BookingQueue[weekdayString][b.Time], models.QueueEntry{
				CustomerName: b.Customer.FirstName + " " + b.Customer.LastName,
				WeekDay:      b.WeekDay,
				Time:         b.Time,
			})

			// Return the queue number to the customer
			fmt.Fprintf(w, "No table is available for immediate booking. You have been added to the queue for %v at %v'O clock. Your queue number is %v.\n", weekdayString, b.Time, queueNum)
		}
	} else {
		// If it's a "Book Later" request, add it to the queue
		log.Printf("No immediate availability, adding to queue for %v at %v'O clock\n", b.WeekDay, b.Time)

		// Convert WeekDay enum to string for accessing the BookingQueue
		weekdayString := models.WeekDayToString[b.WeekDay]

		// Initialize the booking queue for the weekday and time if necessary
		if BookingQueue[weekdayString] == nil {
			BookingQueue[weekdayString] = make(map[int][]models.QueueEntry)
		}

		if BookingQueue[weekdayString][b.Time] == nil {
			BookingQueue[weekdayString][b.Time] = make([]models.QueueEntry, 0)
		}

		// Add the customer to the queue
		queueNum := len(BookingQueue[weekdayString][b.Time]) + 1
		BookingQueue[weekdayString][b.Time] = append(BookingQueue[weekdayString][b.Time], models.QueueEntry{
			CustomerName: b.Customer.FirstName + " " + b.Customer.LastName,
			WeekDay:      b.WeekDay,
			Time:         b.Time,
		})

		// Return the queue number to the customer
		fmt.Fprintf(w, "No table is available for immediate booking. You have been added to the queue for %v at %v'O clock. Your queue number is %v.\n", weekdayString, b.Time, queueNum)
	}
}

// CancelTable cancels a previously booked table and Process the Queue
func (s *TableServiceImplementation) CancelTable(w http.ResponseWriter, r *http.Request) {
	// Function to cancel the table, takes time and table id as query parameters
	weekdayString := r.URL.Query().Get("weekday")
	time := r.URL.Query().Get("time")
	tableid := r.URL.Query().Get("tableid")

	// Convert time to integer
	timeInt, err := strconv.Atoi(time)
	if err != nil {
		log.Println("Error parsing time")
		fmt.Println(CustomError("Error parsing time", 400))
		return
	}
	go func() {
		// Find the table and set IsEmpty to true
		for _, v := range Tables_Data[weekdayString] {
			if v.Hour == timeInt {
				for _, tab := range v.Table {
					if tab.ID == tableid {
						tab.IsEmpty = true
						// Log with the weekday as a string
						fmt.Fprintf(w, "Table %v at %v '%v has been canceled successfully.\n", tab.ID, weekdayString, v.Hour)
						log.Printf("Table %v at %v '%v canceled successfully.\n", tab.ID, weekdayString, v.Hour)
					}
				}
			}
			// After canceling a table, process the queue
			ProcessQueue(models.StringToWeekDay[weekdayString], timeInt)
			return
		}
	}()

	// If no table was found
	fmt.Fprintf(w, "Table not found.\n")
	log.Println("Table not found.")
}

// AddTable adds new tables to a specific weekday and hour.
func (s *TableServiceImplementation) AddTable(w http.ResponseWriter, r *http.Request) {
	// Parse the query parameters for weekday, hour, and the number of tables to add
	weekdayString := r.URL.Query().Get("weekday")
	hour := r.URL.Query().Get("hour")
	numTables := r.URL.Query().Get("numTables")

	// Convert hour and number of tables to integers
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
	go func() {
		// Adding the specified number of tables for the given weekday and hour
		for i := 0; i < numTablesInt; i++ {
			// Check if the hour exists for the given weekday
			for idx, v := range Tables_Data[weekdayString] {
				if v.Hour == hourInt {
					// Add a new table to the existing list of tables for this hour
					newTableId := strconv.Itoa(len(v.Table) + 1)
					v.Table = append(v.Table, models.Table{ID: newTableId, IsEmpty: true})
					Tables_Data[weekdayString][idx].Table = v.Table
					log.Printf("Added Table %v at %v '%v\n", newTableId, weekdayString, hourInt)
					break
				}
			}
		}

		fmt.Fprintf(w, "%d tables added at %v on %v", numTablesInt, hourInt, weekdayString)
		log.Printf("%d tables added at %v on %v", numTablesInt, hourInt, weekdayString)
	}()
}

// RemoveTable removes a specific table from a weekday and hour.
func (s *TableServiceImplementation) RemoveTable(w http.ResponseWriter, r *http.Request) {
	// Parse the query parameters for weekday, hour, and table ID to remove
	weekdayString := r.URL.Query().Get("weekday")
	hour := r.URL.Query().Get("hour")
	tableid := r.URL.Query().Get("tableid")

	hourInt, err := strconv.Atoi(hour)
	if err != nil {
		log.Println("Error parsing hour")
		fmt.Println(CustomError("Error parsing hour", 400))
		return
	}

	// Remove the specified table from the given weekday and hour
	for dayIdx, v := range Tables_Data[weekdayString] {
		if v.Hour == hourInt {
			// Loop through the tables and remove the specified table by ID
			for tableIdx, tab := range v.Table {
				if tab.ID == tableid {
					// Remove the table from the slice
					Tables_Data[weekdayString][dayIdx].Table = append(v.Table[:tableIdx], v.Table[tableIdx+1:]...)
					fmt.Fprintf(w, "Table %v removed at %v '%v successfully.\n", tableid, weekdayString, hourInt)
					log.Printf("Table %v removed at %v '%v successfully.\n", tableid, weekdayString, hourInt)
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
	// Structured logging to provide context
	log.Println("Fetching tables data...")

	// Pretty-print Tables_Data in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Marshal Tables_Data into a JSON format for better readability
	tablesDataJSON, err := json.MarshalIndent(Tables_Data, "", "  ")
	if err != nil {
		log.Printf("Error marshalling tables data: %v", err)
		http.Error(w, "Failed to retrieve tables data", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "{\n  \"tablesData\": %s\n}", string(tablesDataJSON))

	log.Printf("Successfully fetched tables data.")
}

// GetFreeTables returns the available timings to book a table for a specific weekday.
func (s *TableServiceImplementation) GetFreeTables(w http.ResponseWriter, r *http.Request) {
	weekdayString := mux.Vars(r)["weekday"]

	empty_tables := []int{}
	for _, tables := range Tables_Data[weekdayString] {
		for _, table := range tables.Table {
			if table.IsEmpty {
				empty_tables = append(empty_tables, tables.Hour)
				break
			}
		}
	}
	fmt.Fprintf(w, "Available Timings on %v to book the table are %v", weekdayString, empty_tables)
	log.Printf("Available Timings on %v to book the table are %v", weekdayString, empty_tables)
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
		empty_info[models.StringToWeekDay[days]] = empty_tables
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
}

// TableServiceImplementation implements the TableService interface.
type TableServiceImplementation struct{}
