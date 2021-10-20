package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

type BookOwner struct {
	gorm.Model

	Name  string `json:"name"`
	Email string `json:"email"` //`gorm:"typevarchar(100);unique_index"`
	Books []Book `json:"books"`
}

type Book struct {
	gorm.Model

	Title      string
	Author     string
	CallNumber int //`gorm:"unique_index"`
	PersonID   int
}

var (
	owner = BookOwner{
		Model: gorm.Model{},
		Name:  "Euston Francis",
		Email: "efrancis@solewantgroup.com",
	}
	book = []*Book{
		{Title: "Legend of the Seeker", Author: "Neil Warnock", CallNumber: 1234, PersonID: 1},
		{Title: "Merlin", Author: "Candice Lisa", CallNumber: 1211, PersonID: 1},
		{Title: "Harry Porter", Author: "Obediah Brooks", CallNumber: 1714, PersonID: 15},
	}
)

var db *gorm.DB
var err error

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.")
	}

	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, password, dbPort)

	db, err := gorm.Open(dialect, dbUri)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfuly connected to the database")
	}

	defer db.Close()

	db.AutoMigrate(&BookOwner{})
	db.AutoMigrate(&Book{})

	db.Create(owner)
	for _, idx := range book {
		db.Create(idx)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", getOwners).Methods("GET")
	r.HandleFunc("/:id", getOwner).Methods("")

	http.ListenAndServe(":8080", r)

}

func getOwners(w http.ResponseWriter, r *http.Request) {
	var owners []*BookOwner
	db.Find(owners)
	json.NewEncoder(w).Encode(owners)
}

func getOwner(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var owner BookOwner
	var book []Book

	db.First(&owner, params["id"])
	db.Model(&owner).Related(&book)

	owner.Books = book

	json.NewEncoder(w).Encode(owner)
}
