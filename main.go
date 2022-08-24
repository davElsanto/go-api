package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     string `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   string `json:year`
}

var books []Book

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/books", getBooks)

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	log.Println("test get mux")
}
