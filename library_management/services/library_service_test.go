package services

import (
	"testing"

	"library_management/models"
)

func setupTestLibrary() *Library {
	lib := NewLibrary()
	lib.AddBook(models.Book{ID: 1, Title: "The Go Programming Language", Author: "Alan Donovan", Status: models.StatusAvailable})
	lib.AddBook(models.Book{ID: 2, Title: "Clean Code", Author: "Robert C. Martin", Status: models.StatusAvailable})
	lib.AddBook(models.Book{ID: 3, Title: "Design Patterns", Author: "Erich Gamma", Status: models.StatusAvailable})

	lib.AddMember(models.Member{ID: 101, Name: "Alice"})
	lib.AddMember(models.Member{ID: 102, Name: "Bob"})
	return lib
}

func TestAddBook(t *testing.T) {
	lib := NewLibrary()
	book := models.Book{ID: 1, Title: "Learning Go", Author: "Jon Bodner"}
	lib.AddBook(book)

	stored, exists := lib.GetBook(1)
	if !exists {
		t.Fatalf("expected book with ID 1 to exist")
	}
	if stored.Title != "Learning Go" {
		t.Errorf("expected title 'Learning Go', got '%s'", stored.Title)
	}
	if stored.Status != models.StatusAvailable {
		t.Errorf("expected default status to be 'Available', got '%s'", stored.Status)
	}
}

func TestRemoveBook(t *testing.T) {
	lib := setupTestLibrary()
	lib.RemoveBook(2)

	_, exists := lib.GetBook(2)
	if exists {
		t.Errorf("expected book with ID 2 to be removed")
	}

	avail := lib.ListAvailableBooks()
	if len(avail) != 2 {
		t.Errorf("expected 2 available books, got %d", len(avail))
	}
}

func TestBorrowBook_Success(t *testing.T) {
	lib := setupTestLibrary()

	err := lib.BorrowBook(1, 101)
	if err != nil {
		t.Fatalf("unexpected error borrowing book: %v", err)
	}

	// Verify book status in library
	book, _ := lib.GetBook(1)
	if book.Status != models.StatusBorrowed {
		t.Errorf("expected book status to be 'Borrowed', got '%s'", book.Status)
	}

	// Verify member's borrowed list
	borrowed := lib.ListBorrowedBooks(101)
	if len(borrowed) != 1 {
		t.Fatalf("expected 1 borrowed book for member 101, got %d", len(borrowed))
	}
	if borrowed[0].ID != 1 {
		t.Errorf("expected borrowed book ID to be 1, got %d", borrowed[0].ID)
	}
}

func TestBorrowBook_BookNotFound(t *testing.T) {
	lib := setupTestLibrary()

	err := lib.BorrowBook(999, 101)
	if err == nil {
		t.Fatal("expected error for non-existent book, got nil")
	}
}

func TestBorrowBook_MemberNotFound(t *testing.T) {
	lib := setupTestLibrary()

	err := lib.BorrowBook(1, 999)
	if err == nil {
		t.Fatal("expected error for non-existent member, got nil")
	}
}

func TestBorrowBook_AlreadyBorrowed(t *testing.T) {
	lib := setupTestLibrary()

	err := lib.BorrowBook(1, 101)
	if err != nil {
		t.Fatalf("unexpected error borrowing first time: %v", err)
	}

	// Attempt to borrow same book by another member
	err2 := lib.BorrowBook(1, 102)
	if err2 == nil {
		t.Fatal("expected error when borrowing already borrowed book, got nil")
	}
}

func TestReturnBook_Success(t *testing.T) {
	lib := setupTestLibrary()

	_ = lib.BorrowBook(1, 101)
	err := lib.ReturnBook(1, 101)
	if err != nil {
		t.Fatalf("unexpected error returning book: %v", err)
	}

	// Verify status updated in library
	book, _ := lib.GetBook(1)
	if book.Status != models.StatusAvailable {
		t.Errorf("expected status 'Available', got '%s'", book.Status)
	}

	// Verify removed from member's borrowed books
	borrowed := lib.ListBorrowedBooks(101)
	if len(borrowed) != 0 {
		t.Errorf("expected 0 borrowed books after return, got %d", len(borrowed))
	}
}

func TestReturnBook_NotBorrowedByMember(t *testing.T) {
	lib := setupTestLibrary()

	// Member 101 borrows Book 1
	_ = lib.BorrowBook(1, 101)

	// Member 102 attempts to return Book 1
	err := lib.ReturnBook(1, 102)
	if err == nil {
		t.Fatal("expected error when returning book not borrowed by member, got nil")
	}
}

func TestListAvailableBooks(t *testing.T) {
	lib := setupTestLibrary()

	avail := lib.ListAvailableBooks()
	if len(avail) != 3 {
		t.Fatalf("expected 3 available books, got %d", len(avail))
	}

	_ = lib.BorrowBook(1, 101)

	availAfter := lib.ListAvailableBooks()
	if len(availAfter) != 2 {
		t.Fatalf("expected 2 available books, got %d", len(availAfter))
	}
}

func TestListBorrowedBooks(t *testing.T) {
	lib := setupTestLibrary()

	// Initially empty
	borrowed := lib.ListBorrowedBooks(101)
	if len(borrowed) != 0 {
		t.Errorf("expected 0 borrowed books initially, got %d", len(borrowed))
	}

	_ = lib.BorrowBook(1, 101)
	_ = lib.BorrowBook(3, 101)

	borrowed = lib.ListBorrowedBooks(101)
	if len(borrowed) != 2 {
		t.Fatalf("expected 2 borrowed books, got %d", len(borrowed))
	}
}

