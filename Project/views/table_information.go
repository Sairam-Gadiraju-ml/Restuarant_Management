package views

import (
	"Project/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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
	vars := mux.Vars(r)
	weekday := vars["weekday"]

	empty_tables := []int{}
	empty_info := make(map[models.WeekDay][]int)

	for days, hours := range Tables_Data {
		if models.StringToWeekDay[days] == models.StringToWeekDay[weekday] {
			for _, tables := range hours {
				for _, table := range tables.Table {
					if table.IsEmpty {
						empty_tables = append(empty_tables, tables.Hour)
						break
					}
				}
			}
			empty_info[models.StringToWeekDay[days]] = empty_tables
			break
		}
	}

	log.Printf("Tables available at following hours: %v \n", empty_info)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(empty_info)
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
