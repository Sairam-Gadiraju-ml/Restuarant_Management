package main

import (
	"Project/routes"
	"Project/views"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Creates a logFile.txt if not created already and Opens it to write logs
	file := views.CreateLogFile()

	// Intialize the tables for the day
	views.IntializeTables()
	go views.QueueProcessor()

	// Initialize routes
	routes.InitializeRoutes(router)

	log.Println("Intializing the server at http://localhost:5000/:")
	if err := http.ListenAndServe(":5000", router); err != nil {
		log.Panicln("Error Starting server")
		fmt.Println("Error starting server:", err)
	}

	// Closes the logFile.txt
	defer file.Close()
}
