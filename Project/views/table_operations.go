package views

import (
	"Project/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

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
				BookingQueue[weekdayString] = make(map[int]Queue)
			}

			// Ensure the queue for the requested time is initialized
			if BookingQueue[weekdayString][b.Time].Channel == nil {
				queue := NewQueue(100) // Assuming a capacity of 100 for the queue
				BookingQueue[weekdayString][b.Time] = queue
			}

			// Add the customer to the queue
			err := BookingQueue[weekdayString][b.Time].Enqueue(models.QueueEntry{
				CustomerName: b.Customer.FirstName + " " + b.Customer.LastName,
				WeekDay:      b.WeekDay,
				Time:         b.Time,
			})
			if err != nil {
				http.Error(w, "Failed to add to queue", http.StatusInternalServerError)
				return
			}

			// Return the queue number to the customer
			fmt.Fprintf(w, "No table is available for immediate booking. You have been added to the queue for %v at %v'O clock.\n", weekdayString, b.Time)
		}
	} else {
		// If it's a "Book Later" request, add it to the queue
		log.Printf("No immediate availability, adding to queue for %v at %v'O clock\n", b.WeekDay, b.Time)

		// Convert WeekDay enum to string for accessing the BookingQueue
		weekdayString := models.WeekDayToString[b.WeekDay]

		// Initialize the booking queue for the weekday and time if necessary
		if BookingQueue[weekdayString] == nil {
			BookingQueue[weekdayString] = make(map[int]Queue)
		}

		if BookingQueue[weekdayString][b.Time].Channel == nil {
			queue := NewQueue(100) // Assuming a capacity of 100 for the queue
			BookingQueue[weekdayString][b.Time] = queue
		}

		// Add the customer to the queue
		err := BookingQueue[weekdayString][b.Time].Enqueue(models.QueueEntry{
			CustomerName: b.Customer.FirstName + " " + b.Customer.LastName,
			WeekDay:      b.WeekDay,
			Time:         b.Time,
		})
		if err != nil {
			http.Error(w, "Failed to add to queue", http.StatusInternalServerError)
			return
		}

		// Return the queue number to the customer
		fmt.Fprintf(w, "No table is available for immediate booking. You have been added to the queue for %v at %v'O clock.\n", weekdayString, b.Time)
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

	// Find the table and set IsEmpty to true
	for _, v := range Tables_Data[weekdayString] {
		if v.Hour == timeInt {
			for _, tab := range v.Table {
				if tab.ID == tableid {
					tab.IsEmpty = true
					// Log with the weekday as a string
					fmt.Fprintf(w, "Table %v at %v '%v has been canceled successfully.\n", tab.ID, weekdayString, v.Hour)
					log.Printf("Table %v at %v '%v canceled successfully.\n", tab.ID, weekdayString, v.Hour)

					// After canceling a table, process the queue
					ProcessQueue(models.StringToWeekDay[weekdayString], timeInt)
					return
				}
			}
		}
	}

	// If no table was found
	fmt.Fprintf(w, "Table not found.\n")
	log.Println("Table not found.")
}
