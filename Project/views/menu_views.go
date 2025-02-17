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
			fmt.Fprintf(w, "Error Decoding the Request Body", r)
			return
		}
		for _, val := range MenuItems {
			if val.Name == MenuItem.Name {
				log.Println("Menu Item already Exists")
				fmt.Fprintf(w, "Menu Item already Exists, MenuItem ID: %v \n", val.Id)
				return
			}
			MenuItem.Id = uint(len(MenuItems) + 1)
			MenuItems = append(MenuItems, models.MenuItem{Id: MenuItem.Id, Name: MenuItem.Name, Price: MenuItem.Price})
			log.Printf("Menu Item Created: \n Id: %v, Name: %v, Price: %v \n", MenuItem.Id, MenuItem.Name, MenuItem.Price)
			fmt.Fprintf(w, "Menu Item Created: \n Id: %v, Name: %v, Price: %v \n", MenuItem.Id, MenuItem.Name, MenuItem.Price)
			return
		}
	}
}

func (s *MenuServiceImplementation) HandleMenuItem(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)["Id"]
	Id, err := strconv.Atoi(param)
	if err == nil {
		switch r.Method {
		case http.MethodGet:
			for _, val := range MenuItems {
				if val.Id == uint(Id) {
					log.Printf("Item Id: %v, Item Name: %v,Item Price : %v", val.Id, val.Name, val.Price)
					fmt.Fprintf(w, "Item Id: %v, Item Name: %v,Item Price : %v", val.Id, val.Name, val.Price)
					return
				}
			}
		}
	}

}

type MenuService interface {
	GetMenu(w http.ResponseWriter, r *http.Request)
	HandleMenuItem(w http.ResponseWriter, r *http.Request)
}

type MenuServiceImplementation struct{}
