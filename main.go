package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	// "github.com/Isotop7/liberator-backend/models"
)

var DB *gorm.DB
var books = []Book{}
var shelves = []Shelve{}

// Helper function
func containsBook(books []Book, id int) bool {
	for _, book := range books {
		if book.ID == id {
			return true
		}
	}
	return false
}

// List all Books
func listBooksEndpoint(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, books)
}

// List single book
func listBookEndpoint(ctx *gin.Context) {
	// Get parameter from request
	idParam := ctx.Param("id")

	// Parse id to int
	id, err := strconv.Atoi(idParam)
	// If no int
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "book not found",
		})
		return
	}
	// If negative value
	if id <= 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "book not found",
		})
		return
	}

	// Find correct item
	for _, book := range books {
		if book.ID == id {
			ctx.IndentedJSON(http.StatusOK, book)
			return
		}
	}

	// Default handler when no book is found
	ctx.JSON(http.StatusNotFound, gin.H{
		"message": "book not found",
	})
}

// Create single book
func createBookEndpoint(ctx *gin.Context) {
	var newElement Book

	err := ctx.BindJSON(&newElement)
	if err != nil {
		switch err.(type) {
		case *json.UnmarshalTypeError:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request",
			})
		}
		return
	}

	if containsBook(books, newElement.ID) {
		ctx.JSON(http.StatusConflict, gin.H{
			"message": "duplicate element with id",
			"data":    newElement,
		})
	} else {
		books = append(books, newElement)
		ctx.JSON(http.StatusCreated, newElement)
	}
}

// List all shelves
func listShelvesEndpoint(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, shelves)
}

// Temporary data seed for testing purposes
func seedData() {
	b1 := Book{
		ID:        1,
		Title:     "Erstes Buch",
		Author:    "Max Mustermann",
		Language:  "Deutsch",
		Category:  "Roman",
		ISBN10:    "1234567890",
		ISBN13:    "1234567890123",
		PageCount: 500,
		Rating:    1,
	}
	b2 := Book{ID: 3, Title: "Animal"}
	b3 := Book{ID: 5, Title: "Ein Sommer in Nienburg"}
	books = append(books, b1, b2, b3)
	shelves = append(shelves, Shelve{
		ID:       1,
		Location: "Schlafzimmer",
		Content:  books,
	})
}

func ConnectDB() {
	db, err := gorm.Open("sqlite3", "liberator.db")

	if err != nil {
		log.Panic("Failed to connect to database!")
	}

	db.AutoMigrate()
}

// Main function
func main() {
	log.Println("Starting liberator-backend ...")

	seedData()

	// Setup handlers
	router := gin.Default()
	router.GET("/books", listBooksEndpoint)
	router.POST("/books", createBookEndpoint)
	router.GET("/books/:id", listBookEndpoint)
	router.GET("/shelves", listShelvesEndpoint)

	// Server API
	log.Fatal(router.Run())
}
