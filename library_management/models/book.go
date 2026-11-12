package models

import "fmt"

// Book status constants
const (
	StatusAvailable = "Available"
	StatusBorrowed  = "Borrowed"
)

// Book represents a book in the library system.
type Book struct {
	ID     int
	Title  string
	Author string
	Status string // can be "Available" or "Borrowed"
}

// String returns a formatted representation of the Book.
func (b Book) String() string {
	return fmt.Sprintf("[%d] \"%s\" by %s (Status: %s)", b.ID, b.Title, b.Author, b.Status)
}

