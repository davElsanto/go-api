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
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var bookItem Book
	books = []Book{}

	sqlQuery := "select * from books"

	row, errorSql := db.Query(sqlQuery)
	logFatal(errorSql)

	for row.Next() {
		errorSql := row.Scan(&bookItem.ID, &bookItem.Title, &bookItem.Author, &bookItem.Year)
		logFatal(errorSql)
		books = append(books, bookItem)
	}

	defer row.Close()

	json.NewEncoder(w).Encode(books)

}

func getBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	params := mux.Vars(r)
	sqlQuery := "select * from books where id = $1"

	row := db.QueryRow(sqlQuery, params["id"])

	errorSql := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	logFatal(errorSql)

	json.NewEncoder(w).Encode(book)

}

func addBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	json.NewDecoder(r.Body).Decode(&book)

	sqlQuery := "insert into books (title, author, year) values ($1, $2, $3) RETURNING id;"

	errorSql := db.QueryRow(sqlQuery, &book.Title, &book.Author, &book.Year).Scan(&book.ID)
	logFatal(errorSql)

	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	sqlQuery := "update books set title=$1, author=$2, year=$3 where id=$4;"

	json.NewDecoder(r.Body).Decode(&book)

	row, errorSql := db.Exec(sqlQuery, &book.Title, &book.Author, &book.Year, &book.ID)
	logFatal(errorSql)

	_, errorSql = row.RowsAffected()
	logFatal(errorSql)

	json.NewEncoder(w).Encode(&book)

}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	sqlQuery := "delete from books where id = $1"

	row, errorSql := db.Exec(sqlQuery, params["id"])
	logFatal(errorSql)

	_, errorSql = row.RowsAffected()
	logFatal(errorSql)

	json.NewEncoder(w).Encode("deleted")
}
