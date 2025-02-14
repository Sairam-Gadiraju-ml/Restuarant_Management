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
		new_table = append(new_table, models.Table{TableId: t, IsEmpty: true})
	}

	// Initializing Empty Tables for Hours from 9-22
	var new_hours []models.HourDetails
	for h := 9; h <= 22; h++ {
		new_hours = append(new_hours, models.HourDetails{Hour: h, Table: new_table})
	}

	// Initializing Empty Tables for all Days (7 days from today)
	day := time.Now().Weekday()
	Tables_Data[models.WeekDay(day)] = new_hours
	for i := 0; i < 6; i++ {
		day := (models.WeekDay(day) + models.WeekDay(i)) % 7
		Tables_Data[models.WeekDay(day)] = new_hours
	}

	log.Println("Initializing Tables")
}

// GetTables retrieves the table availability for each day and hour.
func GetTables(w http.ResponseWriter, r *http.Request) {
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

// BookTable books a table for a given weekday and time.
func BookTable(w http.ResponseWriter, r *http.Request) {
	var b models.Book
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		log.Panicln("Error Decoding")
		fmt.Println(CustomError("Error Decoding", 400))
	}
	log.Printf("Booking Table on %v at %v'O clock \n", b.WeekDay, b.Time)
	for _, v := range Tables_Data[b.WeekDay] {
		if v.Hour == b.Time {
			for _, tab := range v.Table {
				if tab.IsEmpty {
					Tables_Data[b.WeekDay][b.Time].Table[tab.TableId-1].IsEmpty = false
					log.Printf("Booked Table %v on %v at %v'O clock \n", tab.TableId, b.WeekDay, b.Time)
					w.WriteHeader(http.StatusAccepted)
					fmt.Fprintf(w, "Booked Table %v on %v at %v'O clock \n", tab.TableId, b.WeekDay, b.Time)
					return
				}
			}
		}
	}

	fmt.Fprintf(w, "No table is available for booking \n")
}

// CancelTable cancels a previously booked table.
func CancelTable(w http.ResponseWriter, r *http.Request) {
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

	tableIdInt, err := strconv.Atoi(tableid)
	if err != nil {
		log.Println("Error parsing table ID")
		fmt.Println(CustomError("Error parsing table ID", 400))
		return
	}

	// Find the table and set IsEmpty to true
	for _, v := range Tables_Data[models.WeekDay(weekDayInt)] {
		if v.Hour == timeInt {
			for _, tab := range v.Table {
				if tab.TableId == tableIdInt {
					tab.IsEmpty = true
					fmt.Fprintf(w, "Table %v at %v '%v has been canceled successfully.\n", tab.TableId, models.WeekDay(weekDayInt), v.Hour)
					log.Printf("Table %v at %v '%v canceled successfully.\n", tab.TableId, models.WeekDay(weekDayInt), v.Hour)
					return
				}
			}
		}
	}

	// If no table was found
	fmt.Fprintf(w, "Table not found.\n")
	log.Println("Table not found.")
}

// GetInfo provides the entire tables' data.
func GetInfo(w http.ResponseWriter, r *http.Request) {
	log.Printf("Printing Tables Data %v", Tables_Data)
	fmt.Fprintf(w, "Tables Data: %v", Tables_Data)
}

// GetFreeTables returns the available timings to book a table for a specific weekday.
func GetFreeTables(w http.ResponseWriter, r *http.Request) {
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

// AddTable adds new tables to a specific weekday and hour.
func AddTable(w http.ResponseWriter, r *http.Request) {
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
				newTableId := len(v.Table) + 1
				v.Table = append(v.Table, models.Table{TableId: newTableId, IsEmpty: true})
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
func RemoveTable(w http.ResponseWriter, r *http.Request) {
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

	tableIdInt, err := strconv.Atoi(tableid)
	if err != nil {
		log.Println("Error parsing table ID")
		fmt.Println(CustomError("Error parsing table ID", 400))
		return
	}

	// Remove the specified table from the given weekday and hour
	for dayIdx, v := range Tables_Data[models.WeekDay(weekDayInt)] {
		if v.Hour == hourInt {
			// Loop through the tables and remove the specified table by ID
			for tableIdx, tab := range v.Table {
				if tab.TableId == tableIdInt {
					// Remove the table from the slice
					Tables_Data[models.WeekDay(weekDayInt)][dayIdx].Table = append(v.Table[:tableIdx], v.Table[tableIdx+1:]...)
					fmt.Fprintf(w, "Table %v removed at %v '%v successfully.\n", tableIdInt, models.WeekDay(weekDayInt), hourInt)
					log.Printf("Table %v removed at %v '%v successfully.\n", tableIdInt, models.WeekDay(weekDayInt), hourInt)
					return
				}
			}
		}
	}

	// If no table was found for removal
	fmt.Fprintf(w, "Table not found.\n")
	log.Println("Table not found.")
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
