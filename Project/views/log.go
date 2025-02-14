package views

import (
	"log"
	"os"
)

func CreateLogFile() *os.File {

	file, err := os.OpenFile("logFile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	return file
}
