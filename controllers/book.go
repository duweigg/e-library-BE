package controllers

import (
	"library/initializers"
	"library/models"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

func GetBookList(c *gin.Context) {
	var bookTypes []models.BookType
	var bookRequest models.BookRequest
	if err := c.ShouldBindJSON(&bookRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var query = initializers.DB.Offset(bookRequest.Page * bookRequest.PageSize).Limit(bookRequest.PageSize)
	if bookRequest.Title != "" {
		query = query.Where("title ilike ?", "%"+bookRequest.Title+"%")
	}
	query.Find(&bookTypes)

	var bookTypeIDs []uint
	for _, bookType := range bookTypes {
		bookTypeIDs = append(bookTypeIDs, bookType.ID)
	}

	var totalCounts []struct {
		BookTypeID uint
		Count      int
	}
	initializers.DB.Model(&models.Book{}).
		Select("book_type_id, COUNT(*) as count").
		Where("book_type_id IN ?", bookTypeIDs).
		Group("book_type_id").
		Scan(&totalCounts)
	var availableCounts []struct {
		BookTypeID uint
		Count      int
	}
	initializers.DB.Model(&models.Book{}).
		Select("book_type_id, COUNT(*) as count").
		Where("book_type_id IN ? AND status = 1", bookTypeIDs).
		Group("book_type_id").
		Scan(&availableCounts)

	totalMap := make(map[uint]int)
	availableMap := make(map[uint]int)

	for _, item := range totalCounts {
		totalMap[item.BookTypeID] = item.Count
	}
	for _, item := range availableCounts {
		availableMap[item.BookTypeID] = item.Count
	}

	var booksResponse []models.BookResponse
	for _, bookType := range bookTypes {
		var bookResponse = models.BookResponse{
			ID:             bookType.ID,
			Name:           bookType.Title,
			TotalCount:     totalMap[bookType.ID],
			AvailableCount: availableMap[bookType.ID],
		}
		booksResponse = append(booksResponse, bookResponse)
	}
	sort.Slice(booksResponse, func(i, j int) bool {
		return booksResponse[i].ID < booksResponse[j].ID
	})
	var count int64
	initializers.DB.Model(&models.BookType{}).Count(&count)
	c.JSON(http.StatusOK, gin.H{"books": booksResponse, "total": count})
}

func BorrowBooks(c *gin.Context) {
	user, _ := c.Get("user")
	userData, _ := user.(models.UserResponse)
	userID := userData.ID

	var bookTypeIDs models.BookIDsPayload
	if err := c.ShouldBindJSON(&bookTypeIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var bookIDs []uint
	records := make([]models.Record, len(bookTypeIDs.BookTypeIDs))
	for i, bookTypeID := range bookTypeIDs.BookTypeIDs {
		var book models.Book
		initializers.DB.Where("book_type_id = ?", bookTypeID).Where("status = 1").First(&book)
		bookIDs = append(bookIDs, book.ID)
		records[i] = models.Record{
			UserID: userID,
			BookID: book.ID,
			DueAt:  time.Now().AddDate(0, 0, 4*7),
		}
	}
	if err := initializers.DB.Model(&models.Book{}).
		Where("id IN ?", bookIDs).
		Update("status", 2).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := initializers.DB.Create(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": nil})
}
