package services

import (
	"errors"
	"fmt"
	"sort"

	"library_management/models"
)

// LibraryManager defines the contract for library management operations.
type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

// Library implements the LibraryManager interface.
// It stores books and members using maps keyed by their respective IDs.
type Library struct {
	Books   map[int]models.Book   // stores all books keyed by book ID
	Members map[int]models.Member // stores all members keyed by member ID
}

// Compile-time check to ensure Library implements LibraryManager.
var _ LibraryManager = (*Library)(nil)

// NewLibrary initializes and returns a new Library instance with empty maps.
func NewLibrary() *Library {
	return &Library{
		Books:   make(map[int]models.Book),
		Members: make(map[int]models.Member),
	}
}

// AddBook adds a new book to the library.
// If the book status is empty, it defaults to StatusAvailable ("Available").
func (l *Library) AddBook(book models.Book) {
	if l.Books == nil {
		l.Books = make(map[int]models.Book)
	}
	if book.Status == "" {
		book.Status = models.StatusAvailable
	}
	l.Books[book.ID] = book
}

// RemoveBook removes a book from the library by its ID.
// If the book is currently in any member's borrowed list, it is removed as well.
func (l *Library) RemoveBook(bookID int) {
	if l.Books == nil {
		return
	}
	delete(l.Books, bookID)

	// Clean up book reference from any member's borrowed books
	for memberID, member := range l.Members {
		var updatedBorrowed []models.Book
		for _, b := range member.BorrowedBooks {
			if b.ID != bookID {
				updatedBorrowed = append(updatedBorrowed, b)
			}
		}
		if len(updatedBorrowed) != len(member.BorrowedBooks) {
			member.BorrowedBooks = updatedBorrowed
			l.Members[memberID] = member
		}
	}
}

// BorrowBook allows a member to borrow a book if it exists and is available.
// It updates the book status to "Borrowed" and adds it to the member's BorrowedBooks slice.
func (l *Library) BorrowBook(bookID int, memberID int) error {
	if l.Books == nil || l.Members == nil {
		return errors.New("library storage is not initialized")
	}

	// 1. Verify book exists
	book, bookExists := l.Books[bookID]
	if !bookExists {
		return fmt.Errorf("book with ID %d not found in the library", bookID)
	}

	// 2. Verify member exists
	member, memberExists := l.Members[memberID]
	if !memberExists {
		return fmt.Errorf("member with ID %d not found", memberID)
	}

	// 3. Verify book availability
	if book.Status != models.StatusAvailable {
		return fmt.Errorf("book '%s' (ID %d) is already borrowed and not available", book.Title, bookID)
	}

	// 4. Update book status to Borrowed in the library map
	book.Status = models.StatusBorrowed
	l.Books[bookID] = book

	// 5. Add book to member's BorrowedBooks slice and save back to map
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member

	return nil
}

// ReturnBook allows a member to return a borrowed book.
// It updates the book status back to "Available" and removes it from the member's BorrowedBooks slice.
func (l *Library) ReturnBook(bookID int, memberID int) error {
	if l.Books == nil || l.Members == nil {
		return errors.New("library storage is not initialized")
	}

	// 1. Verify member exists
	member, memberExists := l.Members[memberID]
	if !memberExists {
		return fmt.Errorf("member with ID %d not found", memberID)
	}

	// 2. Verify book exists in library
	book, bookExists := l.Books[bookID]
	if !bookExists {
		return fmt.Errorf("book with ID %d does not exist in the library records", bookID)
	}

	// 3. Verify book is in member's borrowed list
	foundIndex := -1
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			foundIndex = i
			break
		}
	}
	if foundIndex == -1 {
		return fmt.Errorf("book '%s' (ID %d) was not borrowed by member '%s' (ID %d)", book.Title, bookID, member.Name, memberID)
	}

	// 4. Remove book from member's BorrowedBooks slice using slice manipulation
	member.BorrowedBooks = append(member.BorrowedBooks[:foundIndex], member.BorrowedBooks[foundIndex+1:]...)
	l.Members[memberID] = member

	// 5. Update book status in the library map back to Available
	book.Status = models.StatusAvailable
	l.Books[bookID] = book

	return nil
}

// ListAvailableBooks returns a slice of all books with Status "Available".
// Results are sorted by book ID for consistent ordering.
func (l *Library) ListAvailableBooks() []models.Book {
	var available []models.Book
	for _, book := range l.Books {
		if book.Status == models.StatusAvailable {
			available = append(available, book)
		}
	}

	// Sort by ID ascending
	sort.Slice(available, func(i, j int) bool {
		return available[i].ID < available[j].ID
	})

	if available == nil {
		return []models.Book{}
	}
	return available
}

// ListBorrowedBooks returns a slice of all books borrowed by a specific member.
func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	member, exists := l.Members[memberID]
	if !exists {
		return []models.Book{}
	}

	// Return a copy of the slice to avoid external mutation
	borrowed := make([]models.Book, len(member.BorrowedBooks))
	copy(borrowed, member.BorrowedBooks)
	return borrowed
}

// -----------------------------------------------------------------------------
// Helper Methods for Member & General Data Management
// -----------------------------------------------------------------------------

// AddMember registers a new member to the library.
func (l *Library) AddMember(member models.Member) error {
	if l.Members == nil {
		l.Members = make(map[int]models.Member)
	}
	if _, exists := l.Members[member.ID]; exists {
		return fmt.Errorf("member with ID %d already exists", member.ID)
	}
	if member.BorrowedBooks == nil {
		member.BorrowedBooks = []models.Book{}
	}
	l.Members[member.ID] = member
	return nil
}

// GetMember retrieves a member by ID.
func (l *Library) GetMember(memberID int) (models.Member, bool) {
	m, exists := l.Members[memberID]
	return m, exists
}

// ListMembers returns all registered members sorted by ID.
func (l *Library) ListMembers() []models.Member {
	var list []models.Member
	for _, m := range l.Members {
		list = append(list, m)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})
	if list == nil {
		return []models.Member{}
	}
	return list
}

// GetBook retrieves a book by ID.
func (l *Library) GetBook(bookID int) (models.Book, bool) {
	b, exists := l.Books[bookID]
	return b, exists
}

// ListAllBooks returns all books in the library sorted by ID.
func (l *Library) ListAllBooks() []models.Book {
	var list []models.Book
	for _, b := range l.Books {
		list = append(list, b)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})
	if list == nil {
		return []models.Book{}
	}
	return list
}

