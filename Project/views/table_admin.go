package views

import (
	"Project/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

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
