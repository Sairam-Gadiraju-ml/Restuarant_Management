package views

import (
	"fmt"
	"log"
	"os"
	"time"
)

func CreateLogFile() *os.File {
	year, month, day := time.Now().Local().Date()
	dateStr := fmt.Sprintf("%d-%02d-%02d", year, month, day)

	// Ensure the logs directory exists
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err = os.Mkdir("logs", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Create/Open the log file in the logs directory
	file, err := os.OpenFile("logs/"+dateStr+" logFile.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	return file
}
