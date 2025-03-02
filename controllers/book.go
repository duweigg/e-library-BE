package controllers

import (
	"library/models"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Define a struct to hold the database instance
type BookController struct {
	DB *gorm.DB
}

// Constructor function to create a new BookController
func NewBookController(db *gorm.DB) *BookController {
	return &BookController{DB: db}
}

func (bc *BookController) GetBookList(c *gin.Context) {
	var bookTypes []models.BookType
	var bookRequest models.BookRequest

	// Bind JSON and return 422 Unprocessable Entity on failure
	if err := c.ShouldBindJSON(&bookRequest); err != nil {
		log.Printf("Invalid request payload: %v\n", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	var query = bc.DB.Offset(bookRequest.Page * bookRequest.PageSize).Limit(bookRequest.PageSize)
	if bookRequest.Title != "" {
		query = query.Where("title ilike ?", "%"+bookRequest.Title+"%")
	}

	// Fetch book types and return 500 Internal Server Error on failure
	if err := query.Find(&bookTypes).Error; err != nil {
		log.Printf("Database error fetching book list: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch book list"})
		return
	}

	if len(bookTypes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No books found"})
		return
	}

	var bookTypeIDs []uint
	for _, bookType := range bookTypes {
		bookTypeIDs = append(bookTypeIDs, bookType.ID)
	}

	var totalCounts []struct {
		BookTypeID uint
		Count      int
	}
	if err := bc.DB.Model(&models.Book{}).
		Select("book_type_id, COUNT(*) as count").
		Where("book_type_id IN ?", bookTypeIDs).
		Group("book_type_id").
		Scan(&totalCounts).Error; err != nil {
		log.Printf("Error fetching total book count: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total book count"})
		return
	}

	var availableCounts []struct {
		BookTypeID uint
		Count      int
	}
	if err := bc.DB.Model(&models.Book{}).
		Select("book_type_id, COUNT(*) as count").
		Where("book_type_id IN ? AND status = 1", bookTypeIDs).
		Group("book_type_id").
		Scan(&availableCounts).Error; err != nil {
		log.Printf("Error fetching available book count: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available book count"})
		return
	}

	// Prepare response
	booksResponse := PrepareBookResponses(bookTypes, totalCounts, availableCounts)
	var count int64
	bc.DB.Model(&models.BookType{}).Count(&count)
	c.JSON(http.StatusOK, gin.H{"books": booksResponse, "total": count})
}

func (bc *BookController) BorrowBooks(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		log.Println("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userData, ok := user.(models.UserResponse)
	if !ok {
		log.Println("Invalid user data in request")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user data"})
		return
	}
	userID := userData.ID

	var bookTypeIDs models.BookIDsPayload
	if err := c.ShouldBindJSON(&bookTypeIDs); err != nil {
		log.Printf("Invalid borrow request payload: %v\n", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	if len(bookTypeIDs.BookTypeIDs) == 0 {
		log.Println("Empty book borrow request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No book IDs provided"})
		return
	}

	// Begin transaction
	tx := bc.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("Transaction panic, rolled back")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		}
	}()

	var bookIDs []uint
	var records []models.Record

	for _, bookTypeID := range bookTypeIDs.BookTypeIDs {
		var book models.Book

		// Fetch an available book
		if err := tx.Where("book_type_id = ? AND status = 1", bookTypeID).
			First(&book).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				tx.Rollback()
				log.Printf("No available books for book_type_id: %d\n", bookTypeID)
				c.JSON(http.StatusNotFound, gin.H{"error": "No available books found for book_type_id", "book_type_id": bookTypeID})
				return
			}
			tx.Rollback()
			log.Printf("Error fetching available book: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available book"})
			return
		}

		// Append to borrow list
		bookIDs = append(bookIDs, book.ID)
		records = append(records, models.Record{
			UserID: userID,
			BookID: book.ID,
			DueAt:  time.Now().AddDate(0, 0, 28), // 4 weeks
		})
	}

	// Update book status
	if err := tx.Model(&models.Book{}).
		Where("id IN ?", bookIDs).
		Update("status", 2).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating book status: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book status"})
		return
	}

	// Create borrowing records
	if err := tx.Create(&records).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating borrow records: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create borrow records"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("Transaction commit failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}

	// Return response
	log.Printf("User %d successfully borrowed %d books\n", userData.ID, len(records))
	c.JSON(http.StatusOK, gin.H{"message": "Books borrowed successfully", "data": records})
}

// Prepare book responses
func PrepareBookResponses(bookTypes []models.BookType, totalCounts, availableCounts []struct {
	BookTypeID uint
	Count      int
}) []models.BookResponse {
	// Convert counts into maps for fast lookup
	totalMap := make(map[uint]int)
	availableMap := make(map[uint]int)

	for _, item := range totalCounts {
		totalMap[item.BookTypeID] = item.Count
	}
	for _, item := range availableCounts {
		availableMap[item.BookTypeID] = item.Count
	}

	// Generate response using ToResponse() method
	var booksResponse []models.BookResponse
	for _, bookType := range bookTypes {
		booksResponse = append(booksResponse, bookType.ToResponse(totalMap[bookType.ID], availableMap[bookType.ID]))
	}

	// Sort responses by ID
	sort.Slice(booksResponse, func(i, j int) bool {
		return booksResponse[i].ID < booksResponse[j].ID
	})

	return booksResponse
}
