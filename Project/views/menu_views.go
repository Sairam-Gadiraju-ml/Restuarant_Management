package views

import (
	"Project/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var MenuItems []models.MenuItem

func (s *MenuServiceImplementation) GetMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		log.Println("Listing all the Menu Items")
		for _, val := range MenuItems {
			w.WriteHeader(http.StatusOK)
			log.Printf("Menu Id : %v , Item Name: %v , Item Price : %v", val.Id, val.Name, val.Price)
			fmt.Fprintf(w, "Menu Id : %v , Item Name: %v , Item Price : %v", val.Id, val.Name, val.Price)
			return
		}
	} else if r.Method == http.MethodPost {
		// check if item already exists in the menu
		var MenuItem models.MenuItem
		err := json.NewDecoder(r.Body).Decode(&MenuItem)
		if err != nil {
			log.Println("Error Decoding the Request Body")
			fmt.Fprintln(w, "Error Decoding the Request Body", r)
			return
		}
		for _, val := range MenuItems {
			if val.Name == MenuItem.Name {
				log.Println("Menu Item already Exists")
				fmt.Fprintf(w, "Menu Item already Exists, MenuItem ID: %v \n", val.Id)
				return
			}
		}
		MenuItem.Id = uint(len(MenuItems) + 1)
		MenuItems = append(MenuItems, models.MenuItem{Id: MenuItem.Id, Name: MenuItem.Name, Price: MenuItem.Price})
		log.Printf("Menu Item Created: \n Id: %v, Name: %v, Price: %v \n", MenuItem.Id, MenuItem.Name, MenuItem.Price)
		fmt.Fprintf(w, "Menu Item Created: \n Id: %v, Name: %v, Price: %v \n", MenuItem.Id, MenuItem.Name, MenuItem.Price)
	}
}

func (s *MenuServiceImplementation) HandleMenuItem(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)["Id"]
	Id, err := strconv.Atoi(param)
	if err != nil {
		log.Println("Error Parsing ID from Request")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Find menu item by ID
	var menuItemIndex int
	itemFound := false
	for i, val := range MenuItems {
		if val.Id == uint(Id) {
			menuItemIndex = i
			itemFound = true
			break
		}
	}

	if !itemFound {
		log.Printf("Menu Item with ID %v not found", Id)
		http.Error(w, "Menu Item Not Found", http.StatusNotFound)
		return
	}

	// Handle HTTP methods
	switch r.Method {
	case http.MethodGet:
		// Get the details of a menu item
		item := MenuItems[menuItemIndex]
		log.Printf("Item Id: %v, Item Name: %v, Item Price: %v", item.Id, item.Name, item.Price)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Item Id: %v, Item Name: %v, Item Price: %v", item.Id, item.Name, item.Price)

	case http.MethodPut:
		// Update the menu item
		var updatedItem models.MenuItem
		err := json.NewDecoder(r.Body).Decode(&updatedItem)
		if err != nil {
			log.Println("Error Decoding the Request Body for Update")
			http.Error(w, "Error Decoding the Request Body", http.StatusBadRequest)
			return
		}

		// Update the item details
		MenuItems[menuItemIndex].Name = updatedItem.Name
		MenuItems[menuItemIndex].Price = updatedItem.Price

		log.Printf("Menu Item Updated: Id: %v, Name: %v, Price: %v", updatedItem.Id, updatedItem.Name, updatedItem.Price)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Menu Item Updated: Id: %v, Name: %v, Price: %v", updatedItem.Id, updatedItem.Name, updatedItem.Price)

	case http.MethodDelete:
		// Delete the menu item
		MenuItems = append(MenuItems[:menuItemIndex], MenuItems[menuItemIndex+1:]...)
		log.Printf("Menu Item Deleted: Id: %v", Id)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Menu Item Deleted: Id: %v", Id)
	}
}

type MenuService interface {
	GetMenu(w http.ResponseWriter, r *http.Request)
	HandleMenuItem(w http.ResponseWriter, r *http.Request)
}

type MenuServiceImplementation struct{}
