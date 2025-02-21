package views

import (
	"Project/models"
	"fmt"
	"log"
	"time"
)

// ProcessQueue processes the queue of requests.
func ProcessQueue(weekday models.WeekDay, hour int) {
	// Convert WeekDay enum to string for accessing Tables_Data and BookingQueue
	weekdayString := models.WeekDayToString[weekday]

	// Check if the weekday queue exists in BookingQueue
	if _, ok := BookingQueue[weekdayString]; !ok {
		log.Println("No queue for this weekday:", weekdayString)
		return
	}

	// If there are any entries in the queue for the given weekday and hour
	if len(BookingQueue[weekdayString][hour].Channel) > 0 {
		// Dequeue the first customer in the queue
		customer := BookingQueue[weekdayString][hour].Dequeue()

		// Attempt to book a table for this customer
		booked := false
		for _, v := range Tables_Data[weekdayString] {
			if v.Hour == hour {
				for j, tab := range v.Table {
					if tab.IsEmpty {
						// Mark the table as booked
						tab.IsEmpty = false

						// Update the table in the data structure
						Tables_Data[weekdayString][hour-9].Table[j] = tab // Adjust for index (assuming hours start at 9)

						log.Printf("Booking table for customer %v at %v %v", customer.CustomerName, weekdayString, hour)

						// Notify the customer that their booking was successful
						fmt.Println("Table booked for customer:", customer.CustomerName)
						booked = true
						break
					}
				}
				if booked {
					break
				}
			}
		}

		if !booked {
			log.Printf("No available table for customer %v at %v %v", customer.CustomerName, weekdayString, hour)
		}
	} else {
		log.Println("No customers in queue for", weekdayString, hour)
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
				if len(queue.Channel) > 0 {
					// Convert weekdayString back to models.WeekDay for processing if needed
					weekday := models.StringToWeekDay[weekdayString]

					// Process the queue for this weekday and hour
					ProcessQueue(weekday, hour)
				}
			}
		}
	}
}
