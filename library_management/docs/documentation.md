# Console-Based Library Management System in Go

## 1. Overview & Objective

The **Console-Based Library Management System** is a modular command-line application designed and implemented in Go. Its primary objective is to demonstrate idiomatic Go language features and software engineering best practices, including:
- **Structs**: Defining domain models (`Book` and `Member`).
- **Interfaces**: Decoupling abstraction from implementation with `LibraryManager`.
- **Methods**: Associating behaviors and business rules with structs.
- **Slices**: Dynamically managing collections of items (such as a member's borrowed books).
- **Maps**: Providing fast, $O(1)$ key-value lookups for books and members by their unique IDs.
- **Error Handling**: Graceful, explicit validation and informative feedback for edge cases.
- **Docker Containerization**: Packaging the application into a lightweight, isolated Linux container for seamless execution across any environment.

---

## 2. Project Directory Structure

The project follows a clean, modular package structure compliant with Go conventions:

```
library_management/
├── main.go                     # Application entry point and bootstrapping
├── controllers/
│   ├── library_controller.go   # Console UI interaction, input parsing, and routing
│   └── library_controller_test.go # Unit tests for controller workflows
├── models/
│   ├── book.go                 # Book domain model & status definitions
│   └── member.go               # Member domain model
├── services/
│   ├── library_service.go      # LibraryManager interface & business logic implementation
│   └── library_service_test.go # Comprehensive unit tests for library operations
├── docs/
│   └── documentation.md        # Complete system documentation and guides
├── Dockerfile                  # Multi-stage Docker build configuration
├── .dockerignore               # Files excluded from Docker build context
├── docker-compose.yml          # Interactive Docker Compose setup
└── go.mod                      # Go module definition
```

---

## 3. Data Models (`models/`)

### 3.1 `Book` Struct (`models/book.go`)

Represents a book in the library catalog.

```go
type Book struct {
    ID     int
    Title  string
    Author string
    Status string // "Available" or "Borrowed"
}
```

- **`ID`** (`int`): Unique numeric identifier for the book.
- **`Title`** (`string`): The title of the book.
- **`Author`** (`string`): The author's name.
- **`Status`** (`string`): Availability status, represented by constants `StatusAvailable` ("Available") and `StatusBorrowed` ("Borrowed").

### 3.2 `Member` Struct (`models/member.go`)

Represents a library patron who can borrow and return books.

```go
type Member struct {
    ID            int
    Name          string
    BorrowedBooks []Book // a slice to hold borrowed books
}
```

- **`ID`** (`int`): Unique numeric identifier for the member.
- **`Name`** (`string`): The full name of the member.
- **`BorrowedBooks`** (`[]Book`): A dynamic slice storing all `Book` records currently checked out by the member.

---

## 4. Interfaces & Implementation (`services/`)

### 4.1 `LibraryManager` Interface

The `LibraryManager` interface establishes the contract for library management operations:

```go
type LibraryManager interface {
    AddBook(book models.Book)
    RemoveBook(bookID int)
    BorrowBook(bookID int, memberID int) error
    ReturnBook(bookID int, memberID int) error
    ListAvailableBooks() []models.Book
    ListBorrowedBooks(memberID int) []models.Book
}
```

### 4.2 `Library` Struct

The `Library` struct implements `LibraryManager`:

```go
type Library struct {
    Books   map[int]models.Book   // stores all books keyed by book ID
    Members map[int]models.Member // stores all members keyed by member ID
}
```

- **Map Storage**: Both books and members are stored in hash maps (`map[int]models.Book` and `map[int]models.Member`), ensuring constant time $O(1)$ lookups, insertions, and deletions.
- **Compile-Time Verification**: Uses `var _ LibraryManager = (*Library)(nil)` to guarantee that `Library` satisfies the `LibraryManager` interface at compile time.

### 4.3 Core Methods & Business Rules

1. **`AddBook(book models.Book)`**:
   - Adds or updates the book in `Books` map.
   - If `book.Status` is empty, defaults to `"Available"`.

2. **`RemoveBook(bookID int)`**:
   - Deletes the book from `Books` map.
   - Also cleans up any reference to this book from active members' `BorrowedBooks` slices to maintain relational integrity.

3. **`BorrowBook(bookID int, memberID int) error`**:
   - Validates that the book exists in the library catalog. If not, returns `book with ID <id> not found`.
   - Validates that the member exists. If not, returns `member with ID <id> not found`.
   - Checks if `book.Status == "Available"`. If it is already borrowed, returns an error stating that the book is not available.
   - Updates `book.Status` to `"Borrowed"` in the library catalog.
   - Appends the book to the member's `BorrowedBooks` slice.

4. **`ReturnBook(bookID int, memberID int) error`**:
   - Validates that the member exists.
   - Validates that the book exists in the library records.
   - Verifies that the member actually borrowed this book by searching their `BorrowedBooks` slice.
   - Removes the book from `member.BorrowedBooks` using idiomatic Go slice manipulation:
     ```go
     member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
     ```
   - Updates `book.Status` back to `"Available"` in the library catalog.

5. **`ListAvailableBooks() []models.Book`**:
   - Iterates through the `Books` map and filters entries where `Status == "Available"`.
   - Sorts the results ascending by book ID for consistent, user-friendly presentation.

6. **`ListBorrowedBooks(memberID int) []models.Book`**:
   - Returns a safe copy of `member.BorrowedBooks`.
   - Returns an empty slice if the member does not exist or has no borrowed books.

7. **Member Management Extensions**:
   - `AddMember(member models.Member) error`
   - `GetMember(memberID int) (models.Member, bool)`
   - `ListMembers() []models.Member`

---

## 5. Console Controller (`controllers/`)

The `LibraryController` struct coordinates user I/O and bridges user commands with the service layer:

```
==================================================
           LIBRARY MANAGEMENT SYSTEM              
==================================================
1. Add a New Book
2. Remove an Existing Book
3. Borrow a Book
4. Return a Book
5. List All Available Books
6. List All Borrowed Books by a Member
7. Register a New Member
8. List All Members
9. Exit
==================================================
Enter your choice (1-9):
```

### Features:
- **Input Sanitization**: Handles both Unix (`\n`) and Windows (`\r\n`) line endings using `strings.TrimSpace`.
- **Validation**: Rejects invalid numeric inputs gracefully with helpful error messages.
- **Tabular Formatting**: Formats book lists and member records with aligned columns.
- **Dependency Injection**: Accepts `services.LibraryManager` and `io.Reader`, allowing unit tests to simulate automated terminal inputs.

---

## 6. How to Run Locally

### Prerequisites
- Go 1.22+ installed on your system.

### Running Directly
From the root workspace directory or inside `library_management/`:

```bash
# Navigate to the project directory
cd library_management

# Run the application
go run main.go
```

### Running Unit Tests
```bash
# Run all unit tests with verbose output
cd library_management
go test -v ./...
```

---

## 7. How to Run on Docker

The project includes an optimized multi-stage `Dockerfile` that builds a minimal Alpine Linux image.

### 7.1 Build the Docker Image
From the `library_management/` directory (or workspace root):

```bash
# Build the image
docker build -t library-management ./library_management
```

*(Or from within `library_management/`: `docker build -t library-management .`)*

### 7.2 Run Interactively
Because this is an interactive console application, you must run the container with `-it` (interactive TTY) flags:

```bash
docker run -it --rm library-management
```

### 7.3 Run Using Docker Compose (Alternative)
A `docker-compose.yml` file is provided:

```bash
# Build and run with interactive console
docker compose run --rm library-app
```

---

## 8. Sample Interactive Walkthrough

1. **Listing Initial Books**:
   - Choose `5` to view pre-seeded books (e.g., *The Go Programming Language*, *Clean Architecture*).
2. **Borrowing a Book**:
   - Choose `3`.
   - Enter Book ID: `1`.
   - Enter Member ID: `101` (Alice Johnson).
   - System confirms: `[✓] Book #1 successfully borrowed by Member #101!`
3. **Checking Member Borrowed List**:
   - Choose `6`.
   - Enter Member ID: `101`.
   - Displays Book #1 under Alice's borrowed books.
4. **Verifying Book Availability**:
   - Choose `5`.
   - Book #1 is no longer listed in Available Books.
5. **Returning the Book**:
   - Choose `4`.
   - Enter Book ID: `1`.
   - Enter Member ID: `101`.
   - System confirms: `[✓] Book #1 successfully returned by Member #101!`
6. **Exiting**:
   - Choose `9` to exit the application safely.

