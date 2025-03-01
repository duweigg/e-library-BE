package controllers

import (
	"library/initializers"
	"library/models"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRecordList(c *gin.Context) {
	user, _ := c.Get("user")
	userData, _ := user.(models.UserResponse)
	userID := userData.ID

	var recordSearchRequest models.RecordSearchRequest
	if err := c.ShouldBindJSON(&recordSearchRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var records []models.Record
	query := initializers.DB.
		Offset(recordSearchRequest.Page*recordSearchRequest.PageSize).
		Limit(recordSearchRequest.PageSize).
		Preload("User").
		Preload("Book.BookType").
		Where("user_id = ?", userID)

	var isClosed = false
	if recordSearchRequest.Status == 2 {
		isClosed = true
	}
	if recordSearchRequest.Status != 0 {
		query = query.Where("is_closed = ?", isClosed)
	}
	if recordSearchRequest.Title != "" {
		query = query.
			Joins("JOIN books ON books.id = records.book_id").
			Joins("JOIN book_types ON book_types.id = books.book_type_id").
			Where("book_types.title ilike ?", "%"+recordSearchRequest.Title+"%")
	}

	if err := query.Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	recordsResponse := make([]models.RecordResponse, len(records))
	for i, record := range records {
		recordsResponse[i] = record.ToResponse()
	}

	sort.Slice(recordsResponse, func(i, j int) bool {
		return recordsResponse[i].ID < recordsResponse[j].ID
	})

	var total int64
	query = initializers.DB.Model(&models.Record{}).Where("user_id = ?", userID)

	if recordSearchRequest.Status != 0 {
		query = query.Where("is_closed = ?", isClosed)
	}

	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"records": recordsResponse, "total": total})
}
func ExtendRecords(c *gin.Context) {
	var recordIDs models.RecordRequest
	if err := c.ShouldBindJSON(&recordIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := initializers.DB.Model(&models.Record{}).
		Where("id IN ?", recordIDs.IDs).
		Update("due_at", gorm.Expr("due_at + INTERVAL '2 weeks'")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": nil})
}

func ReturnRecords(c *gin.Context) {
	var recordIDs models.RecordRequest
	if err := c.ShouldBindJSON(&recordIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := initializers.DB.Model(&models.Record{}).
		Where("id IN ?", recordIDs.IDs).
		Update("returned_at", time.Now()).
		Update("is_closed", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var records []models.Record
	if err := initializers.DB.Where("id IN ?", recordIDs.IDs).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
		return
	}

	var bookIDs []uint
	for _, record := range records {
		bookIDs = append(bookIDs, record.BookID)
	}

	if err := initializers.DB.Model(&models.Book{}).
		Where("id IN ?", bookIDs).
		Update("status", 1).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": nil})
}
