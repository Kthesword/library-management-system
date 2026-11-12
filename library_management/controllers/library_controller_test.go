package controllers

import (
	"bytes"
	"strings"
	"testing"

	"library_management/models"
	"library_management/services"
)

func TestController_AddAndListAvailableBooks(t *testing.T) {
	lib := services.NewLibrary()

	// Simulate user input for adding a book: ID "10", Title "Concurrency in Go", Author "Katherine Cox-Buday"
	input := "10\nConcurrency in Go\nKatherine Cox-Buday\n"
	inReader := strings.NewReader(input)

	controller := NewLibraryController(lib, inReader)
	controller.AddBook()

	books := lib.ListAvailableBooks()
	if len(books) != 1 {
		t.Fatalf("expected 1 book in library, got %d", len(books))
	}
	if books[0].ID != 10 || books[0].Title != "Concurrency in Go" {
		t.Errorf("unexpected book content: %+v", books[0])
	}
}

func TestController_BorrowAndReturnBook(t *testing.T) {
	lib := services.NewLibrary()
	lib.AddBook(models.Book{ID: 1, Title: "Test Book", Author: "Tester", Status: models.StatusAvailable})
	_ = lib.AddMember(models.Member{ID: 100, Name: "Dave"})

	// Simulate borrowing: Book ID "1", Member ID "100"
	borrowInput := strings.NewReader("1\n100\n")
	controller := NewLibraryController(lib, borrowInput)
	controller.BorrowBook()

	borrowed := lib.ListBorrowedBooks(100)
	if len(borrowed) != 1 {
		t.Fatalf("expected 1 borrowed book, got %d", len(borrowed))
	}

	// Simulate returning: Book ID "1", Member ID "100"
	returnInput := strings.NewReader("1\n100\n")
	controller2 := NewLibraryController(lib, returnInput)
	controller2.ReturnBook()

	borrowedAfter := lib.ListBorrowedBooks(100)
	if len(borrowedAfter) != 0 {
		t.Errorf("expected 0 borrowed books after return, got %d", len(borrowedAfter))
	}
}

func TestController_Exit(t *testing.T) {
	lib := services.NewLibrary()
	// Choice "9" should exit the loop immediately
	input := strings.NewReader("9\n")
	var out bytes.Buffer

	controller := NewLibraryController(lib, input)
	_ = out
	controller.Run()
}

