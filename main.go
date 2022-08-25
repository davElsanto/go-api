package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     int    `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   string `json:year`
}

var books []Book

func main() {
	router := mux.NewRouter()

	books = append(books,
		Book{
			ID:     1,
			Title:  "La culpa es de la vaca",
			Author: "Varios Autores",
			Year:   "2000-08-12",
		},
		Book{
			ID:     2,
			Title:  "El Marquéz de sade",
			Author: "Donatien Alphonse ",
			Year:   "1900-01-02",
		},
	)

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for _, book := range books {
		if book.ID == id {
			json.NewEncoder(w).Encode(book)
		}
	}
}
