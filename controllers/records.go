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
type RecordController struct {
	DB *gorm.DB
}

// Constructor function to create a new BookController
func NewRecordController(db *gorm.DB) *RecordController {
	return &RecordController{DB: db}
}

func (rc *RecordController) GetRecordList(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		log.Println("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userData, _ := user.(models.UserResponse)

	var recordSearchRequest models.RecordSearchRequest
	if err := c.ShouldBindJSON(&recordSearchRequest); err != nil {
		log.Printf("Invalid record search request: %v\n", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	var records []models.Record
	query := rc.DB.
		Offset(recordSearchRequest.Page*recordSearchRequest.PageSize).
		Limit(recordSearchRequest.PageSize).
		Preload("User").
		Preload("Book.BookType").
		Where("user_id = ?", userData.ID)

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
		log.Printf("Failed to fetch records: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
		return
	}
	if len(records) == 0 {
		log.Println("No records found for the given filters")
		c.JSON(http.StatusNotFound, gin.H{"error": "No records found"})
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
	countQuery := rc.DB.Model(&models.Record{}).Where("user_id = ?", userData.ID)
	if recordSearchRequest.Status != 0 {
		countQuery = countQuery.Where("is_closed = ?", isClosed)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		log.Printf("Failed to count records: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total count"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"records": recordsResponse, "total": total})
}

func (rc *RecordController) ExtendRecords(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		log.Println("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userData, _ := user.(models.UserResponse)
	// need to verify if the record belong to user
	var recordIDs models.RecordRequest
	if err := c.ShouldBindJSON(&recordIDs); err != nil {
		log.Printf("Invalid extend request payload: %v\n", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}
	if len(recordIDs.IDs) == 0 {
		log.Println("Empty extend request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No record IDs provided"})
		return
	}
	// Fetch records to verify ownership
	var records []models.Record
	if err := rc.DB.Where("id IN ?", recordIDs.IDs).Find(&records).Error; err != nil {
		log.Printf("Failed to fetch records for extend: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
		return
	}

	// Ensure all records belong to the user
	for _, record := range records {
		if record.UserID != userData.ID {
			log.Printf("User %d attempted to extend non-owned record %d\n", userData.ID, record.ID)
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to extend these records"})
			return
		}
	}

	// Update due dates for all verified records
	if err := rc.DB.Model(&models.Record{}).
		Where("id IN ?", recordIDs.IDs).
		Update("due_at", gorm.Expr("due_at + INTERVAL '3 weeks'")).Error; err != nil {
		log.Printf("Failed to extend records: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extend records"})
		return
	}

	log.Printf("User %d successfully extended %d records\n", userData.ID, len(records))
	c.JSON(http.StatusOK, gin.H{"message": "Records extended successfully"})
}

func (rc *RecordController) ReturnRecords(c *gin.Context) {
	// need to verify if the record belong to user
	// Get the authenticated user
	user, exists := c.Get("user")
	if !exists {
		log.Println("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userData, _ := user.(models.UserResponse)

	var recordIDs models.RecordRequest
	if err := c.ShouldBindJSON(&recordIDs); err != nil {
		log.Printf("Invalid return request payload: %v\n", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}
	if len(recordIDs.IDs) == 0 {
		log.Println("Empty return request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No record IDs provided"})
		return
	}

	// Fetch records to verify ownership
	var recordsChecking []models.Record
	if err := rc.DB.Where("id IN ?", recordIDs.IDs).Find(&recordsChecking).Error; err != nil {
		log.Printf("Failed to fetch records for return: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
		return
	}

	// Ensure all records belong to the user
	for _, record := range recordsChecking {
		if record.UserID != userData.ID {
			log.Printf("User %d attempted to return non-owned record %d\n", userData.ID, record.ID)
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to return these records"})
			return
		}
	}

	// Begin transaction
	tx := rc.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("Transaction panic, rolled back")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		}
	}()
	// Update records as returned
	if err := tx.Model(&models.Record{}).
		Where("id IN ?", recordIDs.IDs).
		Update("returned_at", time.Now()).
		Update("is_closed", true).Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to update return records: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to return records"})
		return
	}

	// Fetch updated records
	var records []models.Record
	if err := tx.Where("id IN ?", recordIDs.IDs).Find(&records).Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to fetch updated records: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated records"})
		return
	}

	// Extract book IDs
	var bookIDs []uint
	for _, record := range records {
		bookIDs = append(bookIDs, record.BookID)
	}
	// Update book status to available
	if err := tx.Model(&models.Book{}).
		Where("id IN ?", bookIDs).
		Update("status", 1).Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to update book status: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book availability"})
		return
	}
	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("Transaction commit failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}

	log.Printf("User %d successfully extended %d records\n", userData.ID, len(records))
	c.JSON(http.StatusOK, gin.H{"message": "Records returned successfully"})
}
