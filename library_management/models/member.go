package models

import "fmt"

// Member represents a registered member of the library.
type Member struct {
	ID            int
	Name          string
	BorrowedBooks []Book // a slice to hold borrowed books
}

// String returns a formatted representation of the Member.
func (m Member) String() string {
	return fmt.Sprintf("Member #%d: %s (Borrowed Books: %d)", m.ID, m.Name, len(m.BorrowedBooks))
}

