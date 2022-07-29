package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Book
type Book struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Writer    string `json:"writer"`
	Language  string `json:"language"`
	Category  string `json:"category"`
	ISBN10    string `json:"isbn10"`
	ISBN13    string `json:"isbn13"`
	PageCount int    `json:"page_count"`
	Rating    int    `json:"rating"`
}

//Shelve
type Shelve struct {
	ID    int    `json:"id"`
	Books []Book `json:"books"`
}

// Temporary storage for all shelves
var library = []Shelve{}

// Handle for all shelves
func handleAllShelves(w http.ResponseWriter, r *http.Request) {
	log.Println("API request to /shelves")
	for i := 0; i < len(library); i++ {
		w.Header().Set("Content-Type", "text/json")
		json.NewEncoder(w).Encode(library)
	}
}

// Handle for single shelve
func handleShelve(w http.ResponseWriter, r *http.Request) {
	log.Println("API request to /shelve/*")
	suffix := r.URL.Path[len("/shelve/"):]
	log.Printf("API request for shelve with id '%v'\n", suffix)
	id, err := strconv.Atoi(suffix)
	if err != nil {
		fmt.Fprintf(w, "Invalid shelve id!")
		log.Println("Invalid shelve id!")
		return
	}
	if id <= 0 {
		fmt.Fprintf(w, "Invalid shelve id!")
		log.Println("Invalid shelve id!")
		return
	}

	for _, shelve := range library {
		if shelve.ID == id {
			w.Header().Set("Content-Type", "text/json")
			json.NewEncoder(w).Encode(shelve)
		}
	}
}

// Handle for single book
func handleBook(w http.ResponseWriter, r *http.Request) {
	log.Println("API request to /book/*")
	suffix := r.URL.Path[len("/book/"):]
	log.Printf("API request for book with id '%v'\n", suffix)
	id, err := strconv.Atoi(suffix)
	if err != nil {
		fmt.Fprintf(w, "Invalid book id!")
		log.Println("Invalid book id!")
		return
	}
	if id <= 0 {
		fmt.Fprintf(w, "Invalid book id!")
		log.Println("Invalid book id!")
		return
	}

	switch r.Method {
	case "GET":
		log.Printf("GET | Book | %v\n", id)
		for _, shelve := range library {
			for _, book := range shelve.Books {
				if book.ID == id {
					w.Header().Set("Content-Type", "text/json")
					json.NewEncoder(w).Encode(book)
					return
				}
			}
		}
		json.NewEncoder(w).Encode("No book found!")
	case "POST":
		log.Printf("POST | Book | %v\n", id)
		rawBody, _ := ioutil.ReadAll(r.Body)
		var newBook Book
		err := json.Unmarshal(rawBody, &newBook)
		if err != nil {
			fmt.Fprintf(w, "Invalid body!")
			log.Printf("Invalid body: %v\n", r.Body)
			if e, ok := err.(*json.SyntaxError); ok {
				log.Printf("syntax error at byte offset %d", e.Offset)
			}
			return
		}
		library[0].Books = append(library[0].Books, newBook)
		log.Printf("Added new book: %v\n", newBook)
	default:
		json.NewEncoder(w).Encode("Invalid request!")
	}
}

// Temporary data seed for testing purposes
func seedData() {
	b1 := Book{
		ID:        1,
		Title:     "Erstes Buch",
		Writer:    "Max Mustermann",
		Language:  "Deutsch",
		Category:  "Roman",
		ISBN10:    "1234567890",
		ISBN13:    "1234567890123",
		PageCount: 500,
		Rating:    1,
	}
	b2 := Book{ID: 3, Title: "Animal"}
	b3 := Book{ID: 5, Title: "Ein Sommer in Nienburg"}
	bookList1 := []Book{}
	bookList1 = append(bookList1, b1, b2, b3)
	shelve1 := Shelve{ID: 1, Books: bookList1}

	b4 := Book{ID: 6, Title: "Später"}
	bookList2 := []Book{}
	bookList2 = append(bookList2, b4)
	shelve2 := Shelve{ID: 3, Books: bookList2}

	library = append(library, shelve1, shelve2)
}

// Main function
func main() {
	log.Println("Starting liberator-backend ...")

	seedData()

	// Setup handlers
	http.HandleFunc("/shelves", handleAllShelves)
	http.HandleFunc("/shelve/", handleShelve)
	http.HandleFunc("/book/", handleBook)
	// Server API
	log.Fatal(http.ListenAndServe(":8080", nil))
}
