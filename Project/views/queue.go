package views

import (
	"Project/models"
	"fmt"
)

// Representing the Queue
type Queue struct {
	Channel chan models.QueueEntry
}

// Enqueue adds a new value to the queue.
func (q Queue) Enqueue(val models.QueueEntry) error {
	select {
	case q.Channel <- val:
		return nil
	default:
		return fmt.Errorf("Queue is full, cannot enqueue item")
	}
}

// Dequeue removes and returns the first value from the queue.
func (q Queue) Dequeue() models.QueueEntry {
	return <-q.Channel
}

// NewQueue creates a new queue with the given capacity and returns a pointer to it.
func NewQueue(capacity int) Queue {
	return Queue{
		Channel: make(chan models.QueueEntry, capacity),
	}
}

// TestQueue tests the functionality of the Queue
func TestQueue() {
	// Create a new queue with a capacity of 5
	queue := NewQueue(5)

	// Create sample QueueEntries
	entries := []models.QueueEntry{
		{CustomerName: "Sai Ram", WeekDay: models.Monday, Time: 10},
		{CustomerName: "John Doe", WeekDay: models.Tuesday, Time: 11},
		{CustomerName: "Jane Smith", WeekDay: models.Wednesday, Time: 12},
	}

	// Test Enqueue
	for _, entry := range entries {
		err := queue.Enqueue(entry)
		if err != nil {
			fmt.Println("Enqueue failed:", err)
		}
	}

	// Print the elements present in the queue
	fmt.Println("Elements in the queue:")
	for i := 0; i < len(entries); i++ {
		dequeuedEntry := queue.Dequeue()
		fmt.Printf("CustomerName: %v, WeekDay: %v, Time: %v\n", dequeuedEntry.CustomerName, dequeuedEntry.WeekDay, dequeuedEntry.Time)
	}
}
