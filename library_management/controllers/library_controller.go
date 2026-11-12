package controllers

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"library_management/models"
	"library_management/services"
)

// LibraryController handles console interactions and delegates business logic
// to the LibraryManager service.
type LibraryController struct {
	manager services.LibraryManager
	library *services.Library // concrete reference for member management helpers
	reader  *bufio.Reader
}

// NewLibraryController creates a new controller instance.
func NewLibraryController(manager services.LibraryManager, r io.Reader) *LibraryController {
	if r == nil {
		r = os.Stdin
	}
	lib, _ := manager.(*services.Library)
	return &LibraryController{
		manager: manager,
		library: lib,
		reader:  bufio.NewReader(r),
	}
}

// Run starts the interactive console loop.
func (c *LibraryController) Run() {
	for {
		c.displayMenu()
		choice := c.readLine("Enter your choice (1-9): ")

		switch choice {
		case "1":
			c.AddBook()
		case "2":
			c.RemoveBook()
		case "3":
			c.BorrowBook()
		case "4":
			c.ReturnBook()
		case "5":
			c.ListAvailableBooks()
		case "6":
			c.ListBorrowedBooks()
		case "7":
			c.RegisterMember()
		case "8":
			c.ListMembers()
		case "9":
			fmt.Println("\nThank you for using the Library Management System. Goodbye!")
			return
		default:
			fmt.Println("\n[!] Invalid option. Please enter a number between 1 and 9.")
		}

		c.pause()
	}
}

// displayMenu prints the console menu.
func (c *LibraryController) displayMenu() {
	fmt.Println("\n==================================================")
	fmt.Println("           LIBRARY MANAGEMENT SYSTEM              ")
	fmt.Println("==================================================")
	fmt.Println("1. Add a New Book")
	fmt.Println("2. Remove an Existing Book")
	fmt.Println("3. Borrow a Book")
	fmt.Println("4. Return a Book")
	fmt.Println("5. List All Available Books")
	fmt.Println("6. List All Borrowed Books by a Member")
	fmt.Println("7. Register a New Member")
	fmt.Println("8. List All Members")
	fmt.Println("9. Exit")
	fmt.Println("==================================================")
}

// AddBook prompts the user for book details and adds it to the library.
func (c *LibraryController) AddBook() {
	fmt.Println("\n--- [1] Add a New Book ---")

	id, err := c.readInt("Enter Book ID: ")
	if err != nil {
		fmt.Printf("[x] Error: %v\n", err)
		return
	}

	// Check if book ID already exists if concrete library is available
	if c.library != nil {
		if existing, exists := c.library.GetBook(id); exists {
			fmt.Printf("[x] Error: A book with ID %d already exists (\"%s\").\n", id, existing.Title)
			return
		}
	}

	title := c.readLine("Enter Book Title: ")
	if title == "" {
		fmt.Println("[x] Error: Book title cannot be empty.")
		return
	}

	author := c.readLine("Enter Book Author: ")
	if author == "" {
		fmt.Println("[x] Error: Book author cannot be empty.")
		return
	}

	book := models.Book{
		ID:     id,
		Title:  title,
		Author: author,
		Status: models.StatusAvailable,
	}

	c.manager.AddBook(book)
	fmt.Printf("[✓] Book \"%s\" (ID: %d) successfully added to the library!\n", title, id)
}

// RemoveBook prompts for a book ID and removes it from the library.
func (c *LibraryController) RemoveBook() {
	fmt.Println("\n--- [2] Remove an Existing Book ---")

	id, err := c.readInt("Enter Book ID to remove: ")
	if err != nil {
		fmt.Printf("[x] Error: %v\n", err)
		return
	}

	if c.library != nil {
		book, exists := c.library.GetBook(id)
		if !exists {
			fmt.Printf("[x] Error: Book with ID %d not found.\n", id)
			return
		}
		c.manager.RemoveBook(id)
		fmt.Printf("[✓] Book \"%s\" (ID: %d) successfully removed.\n", book.Title, id)
		return
	}

	c.manager.RemoveBook(id)
	fmt.Printf("[✓] Book with ID %d removed.\n", id)
}

// BorrowBook allows a member to borrow an available book.
func (c *LibraryController) BorrowBook() {
	fmt.Println("\n--- [3] Borrow a Book ---")

	bookID, err := c.readInt("Enter Book ID: ")
	if err != nil {
		fmt.Printf("[x] Error: %v\n", err)
		return
	}

	memberID, err := c.readInt("Enter Member ID: ")
	if err != nil {
		fmt.Printf("[x] Error: %v\n", err)
		return
	}

	err = c.manager.BorrowBook(bookID, memberID)
	if err != nil {
		fmt.Printf("[x] Failed to borrow book: %v\n", err)
		return
	}

	fmt.Printf("[✓] Book #%d successfully borrowed by Member #%d!\n", bookID, memberID)
}

// ReturnBook allows a member to return a previously borrowed book.
func (c *LibraryController) ReturnBook() {
	fmt.Println("\n--- [4] Return a Book ---")

	bookID, err := c.readInt("Enter Book ID: ")
	if err != nil {
		fmt.Printf("[x] Error: %v\n", err)
		return
	}

	memberID, err := c.readInt("Enter Member ID: ")
	if err != nil {
		fmt.Printf("[x] Error: %v\n", err)
		return
	}

	err = c.manager.ReturnBook(bookID, memberID)
	if err != nil {
		fmt.Printf("[x] Failed to return book: %v\n", err)
		return
	}

	fmt.Printf("[✓] Book #%d successfully returned by Member #%d!\n", bookID, memberID)
}

// ListAvailableBooks lists all books with "Available" status.
func (c *LibraryController) ListAvailableBooks() {
	fmt.Println("\n--- [5] Available Books ---")
	books := c.manager.ListAvailableBooks()

	if len(books) == 0 {
		fmt.Println("No books are currently available in the library.")
		return
	}

	c.printBooksTable(books)
}

// ListBorrowedBooks lists all books borrowed by a specific member.
func (c *LibraryController) ListBorrowedBooks() {
	fmt.Println("\n--- [6] Borrowed Books by Member ---")

	memberID, err := c.readInt("Enter Member ID: ")
	if err != nil {
		fmt.Printf("[x] Error: %v\n", err)
		return
	}

	if c.library != nil {
		member, exists := c.library.GetMember(memberID)
		if !exists {
			fmt.Printf("[x] Error: Member with ID %d does not exist.\n", memberID)
			return
		}
		fmt.Printf("Member: %s (ID: %d)\n", member.Name, member.ID)
	}

	books := c.manager.ListBorrowedBooks(memberID)
	if len(books) == 0 {
		fmt.Printf("No books currently borrowed by Member #%d.\n", memberID)
		return
	}

	c.printBooksTable(books)
}

// RegisterMember registers a new library member.
func (c *LibraryController) RegisterMember() {
	fmt.Println("\n--- [7] Register a New Member ---")
	if c.library == nil {
		fmt.Println("[x] Member registration is not supported by the underlying manager.")
		return
	}

	id, err := c.readInt("Enter Member ID: ")
	if err != nil {
		fmt.Printf("[x] Error: %v\n", err)
		return
	}

	name := c.readLine("Enter Member Name: ")
	if name == "" {
		fmt.Println("[x] Error: Member name cannot be empty.")
		return
	}

	member := models.Member{
		ID:            id,
		Name:          name,
		BorrowedBooks: []models.Book{},
	}

	err = c.library.AddMember(member)
	if err != nil {
		fmt.Printf("[x] Failed to register member: %v\n", err)
		return
	}

	fmt.Printf("[✓] Member \"%s\" (ID: %d) successfully registered!\n", name, id)
}

// ListMembers lists all registered members.
func (c *LibraryController) ListMembers() {
	fmt.Println("\n--- [8] Registered Members ---")
	if c.library == nil {
		fmt.Println("[x] Member listing is not supported by the underlying manager.")
		return
	}

	members := c.library.ListMembers()
	if len(members) == 0 {
		fmt.Println("No members are currently registered.")
		return
	}

	fmt.Printf("%-6s | %-25s | %-16s\n", "ID", "Name", "Borrowed Count")
	fmt.Println(strings.Repeat("-", 53))
	for _, m := range members {
		fmt.Printf("%-6d | %-25s | %-16d\n", m.ID, m.Name, len(m.BorrowedBooks))
	}
}

// -----------------------------------------------------------------------------
// Helper Methods for I/O and Formatting
// -----------------------------------------------------------------------------

func (c *LibraryController) printBooksTable(books []models.Book) {
	fmt.Printf("%-6s | %-32s | %-22s | %-12s\n", "ID", "Title", "Author", "Status")
	fmt.Println(strings.Repeat("-", 78))
	for _, b := range books {
		title := b.Title
		if len(title) > 30 {
			title = title[:27] + "..."
		}
		author := b.Author
		if len(author) > 20 {
			author = author[:17] + "..."
		}
		fmt.Printf("%-6d | %-32s | %-22s | %-12s\n", b.ID, title, author, b.Status)
	}
}

func (c *LibraryController) readLine(prompt string) string {
	fmt.Print(prompt)
	line, err := c.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return ""
	}
	return strings.TrimSpace(line)
}

func (c *LibraryController) readInt(prompt string) (int, error) {
	valStr := c.readLine(prompt)
	if valStr == "" {
		return 0, fmt.Errorf("input cannot be empty")
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("'%s' is not a valid integer number", valStr)
	}
	return val, nil
}

func (c *LibraryController) pause() {
	fmt.Print("\nPress Enter to continue...")
	_, _ = c.reader.ReadString('\n')
}

