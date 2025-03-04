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

var mockBookListExpectredReturn = BookListResponse{
	Books: []models.BookResponse{
		{
			ID:             1,
			Name:           "Mock Book 1",
			TotalCount:     2,
			AvailableCount: 1,
		},
		{
			ID:             2,
			Name:           "Mock Book 2",
			TotalCount:     3,
			AvailableCount: 2,
		},
		{
			ID:             3,
			Name:           "Mock Book 3",
			TotalCount:     1,
			AvailableCount: 0,
		},
	},
	Total: 3,
}

type BorrowResponse struct {
	Record  []models.Record `json:"books"`
	Message string          `json:"message"`
}

var mockBorrowSuccessExpectedReturn = BorrowResponse{
	Record: []models.Record{
		{
			ID:       5,
			UserID:   1,
			BookID:   1,
			IsClosed: false,
		},
	},
	Message: "Books borrowed successfully",
}

func TestGetBookList(t *testing.T) {
	db := SetupMockDB()
	PrepareMockBookDB(db)
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

	assert.Equal(t, mockBookListExpectredReturn, bookListResponse)
}

func TestBorrowBooksNotAvailable(t *testing.T) {
	db := SetupMockDB()
	PrepareMockBookDB(db)
	defer db.ConnPool.(*sql.DB).Close()
	router := SetupMockRouter(db)

	var ids = []int{3}
	requestBody, _ := json.Marshal(map[string][]int{
		"ids": ids,
	})
	req, _ := http.NewRequest("POST", "/book/borrow", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBorrowBooksSuccess(t *testing.T) {
	db := SetupMockDB()
	PrepareMockRecordDB(db)
	defer db.ConnPool.(*sql.DB).Close()
	router := SetupMockRouter(db)

	var ids = []int{1}
	requestBody, _ := json.Marshal(map[string][]int{
		"ids": ids,
	})
	req, _ := http.NewRequest("POST", "/book/borrow", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var borrowResponse BorrowResponse
	err := json.NewDecoder(w.Body).Decode(&borrowResponse)
	require.NoError(t, err) // Ensure JSON decoding is successful

	assert.Equal(t, mockBorrowSuccessExpectedReturn.Message, borrowResponse.Message)
	for index, _ := range borrowResponse.Record {
		assert.Equal(t, mockBorrowSuccessExpectedReturn.Record[index].BookID, borrowResponse.Record[index].BookID)
		assert.Equal(t, mockBorrowSuccessExpectedReturn.Record[index].UserID, borrowResponse.Record[index].UserID)
		assert.Equal(t, mockBorrowSuccessExpectedReturn.Record[index].ID, borrowResponse.Record[index].ID)
	}
}
