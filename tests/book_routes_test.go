package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"library/models"
	"net/http"
	"net/http/httptest"
	"testing"

	// Adjust package import based on your project structure

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type BookListResponse struct {
	Books []models.BookResponse `json:"books"`
	Total int                   `json:"total"`
}

var testBookList = BookListResponse{
	Books: []models.BookResponse{
		{
			ID:             1,
			Name:           "All Available Test",
			TotalCount:     2,
			AvailableCount: 2,
		},
	},
	Total: 1,
}

func TestGetBookList(t *testing.T) {
	db := SetupMockDB()
	defer db.ConnPool.(*sql.DB).Close()
	router := SetupMockRouter(db)
	requestBody, _ := json.Marshal(map[string]interface{}{
		"title":     "",
		"page_size": 10,
		"page":      0,
	})

	req, _ := http.NewRequest("POST", "/book/list", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var bookListResponse BookListResponse
	err := json.NewDecoder(w.Body).Decode(&bookListResponse)
	require.NoError(t, err) // Ensure JSON decoding is successful

	assert.Equal(t, testBookList, bookListResponse)
}

// func TestBorrowBooks(t *testing.T) {
// db := SetupMockDB()
// defer db.ConnPool.(*sql.DB).Close()
// router := SetupMockRouter(db)

// requestBody, _ := json.Marshal(map[string]interface{}{
// 	"book_id": 1,
// })

// req, _ := http.NewRequest("POST", "/book/borrow", bytes.NewBuffer(requestBody))

// w := httptest.NewRecorder()
// router.ServeHTTP(w, req)

// assert.Equal(t, http.StatusOK, w.Code)
// }
