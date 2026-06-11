package main

import (
	"net/http"
	"strconv"
	"sync"
	"github.com/gin-gonic/gin"
)

// In-memory database with a mutex to handle concurrent request safety
var (
	books  = make(map[int]Book)
	nextID = 1
	mu     sync.Mutex
)

// 1. READ ALL (GET /books)
func getBooks(c *gin.Context) {
	 mu.Lock()
	 defer mu.Unlock()

	// Convert map values into a slice
	bookList := make([]Book, 0, len(books))
	for _, book := range books {
		bookList = append(bookList, book)
	}

	c.JSON(http.StatusOK, gin.H{"result": bookList})
}

// 2. READ ONE (GET /books/{id})
func getBookById(c *gin.Context) {	
	// Extract the wildcard variable from the path
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		 c.JSON(400, gin.H{"error": "id must be an integer"})
        return
	}

	mu.Lock()
	book, exists := books[id]
	mu.Unlock()

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": book})
}

// 3. CREATE (POST /books)
func createBook(c *gin.Context) {

	var book Book
	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	book.ID = nextID
	books[nextID] = book
	nextID++
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"result": books})
}

// 4. UPDATE (PUT /books/{id})
func updateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		 c.JSON(400, gin.H{"error": "id must be an integer"})
        return
	}

	var updatedBook Book
	if err := c.BindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request body"})
		return
	}
	
	mu.Lock()
	_, exists := books[id]
	if !exists {
		mu.Unlock()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Book not found"})
		return
	}

	updatedBook.ID = id
	books[id] = updatedBook
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"result": updatedBook})
}

// 5. DELETE (DELETE /books/{id})
func deleteBook(c *gin.Context) {
	idStr  := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		 c.JSON(400, gin.H{"error": "id must be an integer"})
        return
	}

	mu.Lock()
	_, exists := books[id]
	if !exists {
		mu.Unlock()
		c.JSON(http.StatusInternalServerError, gin.H{"error": ""})
		return
	}

	delete(books, id)
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"result": "1"})
}
