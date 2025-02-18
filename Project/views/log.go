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

	file, err := os.OpenFile(dateStr+" logFile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	return file
}
