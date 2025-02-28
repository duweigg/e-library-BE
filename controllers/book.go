package controllers

import (
	"library/initializers"
	"library/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetBookList(c *gin.Context) {

	var books []models.Book
	initializers.DB.Scan(&books)

	c.JSON(http.StatusOK, gin.H{"data": books})

}

func BorrowBooks(c *gin.Context) {

	var books []models.Book
	var BookPayload models.BookPayload

	if err := c.ShouldBindJSON(&BookPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	initializers.DB.Where("id in (?)", &BookPayload.ID).Scan(&books)
	if len(books) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
	}
	for _,book := range books {
		if book.Status != 1{
			c.JSON(http.StatusNotFound, gin.H{"error": book.Name+" is not available"})
		}
	}
	for _,book := range books {
		book.Status = 2
		book.AvailableDate = time.Now().Add(28 * 24 * time.Hour)
		initializers.DB.Save(book)
	}
	c.JSON(http.StatusOK, gin.H{"data": books})
}

func ExtendBooks(c *gin.Context) {

	var books []models.Book
	var BookPayload models.BookPayload

	if err := c.ShouldBindJSON(&BookPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	initializers.DB.Where("id in (?)", &BookPayload.ID).Scan(&books)
	if len(books) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
	}
	for _,book := range books {
		if book.Status == 1{
			c.JSON(http.StatusNotFound, gin.H{"error": book.Name+" is not borrowed"})
		}
	}
	for _,book := range books {
		book.AvailableDate = time.Now().Add(14 * 24 * time.Hour)
		initializers.DB.Save(book)
	}
	c.JSON(http.StatusOK, gin.H{"data": books})
}

func ReturnBooks(c *gin.Context) {

	var books []models.Book
	var BookPayload models.BookPayload

	if err := c.ShouldBindJSON(&BookPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	initializers.DB.Where("id in (?)", &BookPayload.ID).Scan(&books)
	if len(books) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
	}
	for _,book := range books {
		if book.Status == 1{
			c.JSON(http.StatusNotFound, gin.H{"error": book.Name+" is not borrowed"})
		}
	}
	for _,book := range books {
		book.Status = 1
		book.AvailableDate = time.Now()
		initializers.DB.Save(book)
	}
	c.JSON(http.StatusOK, gin.H{"data": books})

}