package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"
)

type Book struct {
	ID     int    `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   string `json:year`
}

var books []Book
var db *sql.DB

func init() {
	gotenv.Load()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)

	}
}

func main() {
	pgUrl, errorPg := pq.ParseURL(os.Getenv("ELEPHANT_URL"))
	logFatal(errorPg)

	db, errorPg = sql.Open("postgres", pgUrl)
	logFatal(errorPg)

	errorPg = db.Ping()
	logFatal(errorPg)

	router := mux.NewRouter()

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var bookItem Book
	books = []Book{}

	sqlQuery := "select * from books"

	row, errorPsql := db.Query(sqlQuery)
	logFatal(errorPsql)

	for row.Next() {
		errorPsql := row.Scan(&bookItem.ID, &bookItem.Title, &bookItem.Author, &bookItem.Year)
		logFatal(errorPsql)
		books = append(books, bookItem)
	}

	defer row.Close()

	json.NewEncoder(w).Encode(books)

}

func getBook(w http.ResponseWriter, r *http.Request) {
}

func addBook(w http.ResponseWriter, r *http.Request) {
}

func updateBook(w http.ResponseWriter, r *http.Request) {
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
}
