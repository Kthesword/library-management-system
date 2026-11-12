package main

import (
	"fmt"
	"os"

	"library_management/controllers"
	"library_management/models"
	"library_management/services"
)

// seedInitialData populates the library with sample books and members for quick demonstration.
func seedInitialData(lib *services.Library) {
	sampleBooks := []models.Book{
		{ID: 1, Title: "The Go Programming Language", Author: "Alan Donovan & Brian Kernighan", Status: models.StatusAvailable},
		{ID: 2, Title: "Clean Architecture", Author: "Robert C. Martin", Status: models.StatusAvailable},
		{ID: 3, Title: "Designing Data-Intensive Applications", Author: "Martin Kleppmann", Status: models.StatusAvailable},
		{ID: 4, Title: "The Pragmatic Programmer", Author: "Andrew Hunt & David Thomas", Status: models.StatusAvailable},
		{ID: 5, Title: "Introduction to Algorithms", Author: "Cormen, Leiserson, Rivest, Stein", Status: models.StatusAvailable},
	}

	for _, b := range sampleBooks {
		lib.AddBook(b)
	}

	sampleMembers := []models.Member{
		{ID: 101, Name: "Alice Johnson", BorrowedBooks: []models.Book{}},
		{ID: 102, Name: "Bob Smith", BorrowedBooks: []models.Book{}},
		{ID: 103, Name: "Charlie Brown", BorrowedBooks: []models.Book{}},
	}

	for _, m := range sampleMembers {
		_ = lib.AddMember(m)
	}
}

func main() {
	fmt.Println("Initializing Library Management System...")

	// Create a new Library service instance implementing LibraryManager interface
	library := services.NewLibrary()

	// Seed with default initial data
	seedInitialData(library)

	// Create and launch the console controller
	controller := controllers.NewLibraryController(library, os.Stdin)
	controller.Run()
}

