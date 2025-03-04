package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"library/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type RecordListResponse struct {
	Records []models.RecordResponse `json:"records"`
	Total   int                     `json:"total"`
}

var mockRecordListExpectedReturn = RecordListResponse{
	Records: []models.RecordResponse{
		{
			ID:     1,
			Name:   "Mock",
			Title:  "Mock Book 1",
			DueAt:  parsedOverdueAt.Local(),
			Status: "Overdue",
		},
		{
			ID:     3,
			Name:   "Mock",
			Title:  "Mock Book 3",
			DueAt:  parsedDueAt.Local(),
			Status: "Borrowed",
		},
		{
			ID:     4,
			Name:   "Mock",
			Title:  "Mock Book 1",
			DueAt:  parsedDueAt.Local(),
			Status: "Returned",
		},
	},
	Total: 3,
}

func TestGetRecordList(t *testing.T) {
	db := SetupMockDB()
	PrepareMockRecordDB(db)
	defer db.ConnPool.(*sql.DB).Close()
	router := SetupMockRouter(db)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"title":     "",
		"page_size": 10,
		"page":      0,
		"status":    0,
	})
	req, _ := http.NewRequest("POST", "/record/list", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var recordListResponse RecordListResponse
	err := json.NewDecoder(w.Body).Decode(&recordListResponse)
	require.NoError(t, err) // Ensure JSON decoding is successful

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, mockRecordListExpectedReturn.Total, recordListResponse.Total)
	for index, _ := range recordListResponse.Records {
		assert.Equal(t, mockRecordListExpectedReturn.Records[index].Title, recordListResponse.Records[index].Title)
		assert.Equal(t, mockRecordListExpectedReturn.Records[index].Name, recordListResponse.Records[index].Name)
		assert.Equal(t, mockRecordListExpectedReturn.Records[index].ID, recordListResponse.Records[index].ID)
		assert.Equal(t, mockRecordListExpectedReturn.Records[index].DueAt, recordListResponse.Records[index].DueAt)
		assert.Equal(t, mockRecordListExpectedReturn.Records[index].Status, recordListResponse.Records[index].Status)
	}
}

func TestExtendRecordsOtherUserRecord(t *testing.T) {
	db := SetupMockDB()
	PrepareMockRecordDB(db)
	defer db.ConnPool.(*sql.DB).Close()
	router := SetupMockRouter(db)

	var ids = []int{2}
	requestBody, _ := json.Marshal(map[string][]int{
		"ids": ids,
	})

	req, _ := http.NewRequest("POST", "/record/extend", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
func TestExtendRecordsOverdueRecord(t *testing.T) {
	db := SetupMockDB()
	PrepareMockRecordDB(db)
	defer db.ConnPool.(*sql.DB).Close()
	router := SetupMockRouter(db)

	var ids = []int{1}
	requestBody, _ := json.Marshal(map[string][]int{
		"ids": ids,
	})

	req, _ := http.NewRequest("POST", "/record/extend", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
func TestExtendRecordsClosedRecord(t *testing.T) {
	db := SetupMockDB()
	PrepareMockRecordDB(db)
	defer db.ConnPool.(*sql.DB).Close()
	router := SetupMockRouter(db)

	var ids = []int{4}
	requestBody, _ := json.Marshal(map[string][]int{
		"ids": ids,
	})

	req, _ := http.NewRequest("POST", "/record/extend", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
func TestExtendRecordsSuccess(t *testing.T) {
	db := SetupMockDB()
	PrepareMockRecordDB(db)
	defer db.ConnPool.(*sql.DB).Close()
	router := SetupMockRouter(db)

	var ids = []int{3}
	requestBody, _ := json.Marshal(map[string][]int{
		"ids": ids,
	})

	req, _ := http.NewRequest("POST", "/record/extend", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestReturnRecordsOtherUserRecord(t *testing.T) {
	db := SetupMockDB()
	PrepareMockRecordDB(db)
	defer db.ConnPool.(*sql.DB).Close()
	router := SetupMockRouter(db)

	var ids = []int{2}
	requestBody, _ := json.Marshal(map[string][]int{
		"ids": ids,
	})

	req, _ := http.NewRequest("POST", "/record/extend", bytes.NewBuffer(requestBody))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestReturnRecordsClosedRecord(t *testing.T) {
	db := SetupMockDB()
	PrepareMockRecordDB(db)
	defer db.ConnPool.(*sql.DB).Close()
	router := SetupMockRouter(db)

	var ids = []int{4}
	requestBody, _ := json.Marshal(map[string][]int{
		"ids": ids,
	})

	req, _ := http.NewRequest("POST", "/record/extend", bytes.NewBuffer(requestBody))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestReturnRecordsSuccess(t *testing.T) {
	db := SetupMockDB()
	PrepareMockRecordDB(db)
	defer db.ConnPool.(*sql.DB).Close()
	router := SetupMockRouter(db)

	var ids = []int{3}
	requestBody, _ := json.Marshal(map[string][]int{
		"ids": ids,
	})

	req, _ := http.NewRequest("POST", "/record/extend", bytes.NewBuffer(requestBody))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
