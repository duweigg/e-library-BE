package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserRepeatedUsername(t *testing.T) {
	db := SetupMockDB()
	defer db.ConnPool.(*sql.DB).Close()
	router := SetupMockRouter(db)

	requestBody, _ := json.Marshal(map[string]string{
		"username": "mock",
		"password": "admin123",
		"nickname": "Mock",
	})

	req, _ := http.NewRequest("POST", "/user/signup", bytes.NewBuffer(requestBody))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestCreateUserSuccess(t *testing.T) {
	db := SetupMockDB()
	defer db.ConnPool.(*sql.DB).Close()
	router := SetupMockRouter(db)

	requestBody, _ := json.Marshal(map[string]string{
		"username": "mock_success",
		"password": "admin123",
		"nickname": "Mock",
	})

	req, _ := http.NewRequest("POST", "/user/signup", bytes.NewBuffer(requestBody))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSignInWrongInfo(t *testing.T) {
	db := SetupMockDB()
	defer db.ConnPool.(*sql.DB).Close()
	router := SetupMockRouter(db)

	requestBody, _ := json.Marshal(map[string]string{
		"username": "mock",
		"password": "admin123",
	})

	req, _ := http.NewRequest("POST", "/user/signin", bytes.NewBuffer(requestBody))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
func TestSignInSuccess(t *testing.T) {
	db := SetupMockDB()
	defer db.ConnPool.(*sql.DB).Close()
	router := SetupMockRouter(db)

	requestBody, _ := json.Marshal(map[string]string{
		"username": "mock",
		"password": "admin",
	})

	req, _ := http.NewRequest("POST", "/user/signin", bytes.NewBuffer(requestBody))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
func TestGetInfo(t *testing.T) {
	db := SetupMockDB()
	defer db.ConnPool.(*sql.DB).Close()
	router := SetupMockRouter(db)

	req, _ := http.NewRequest("GET", "/user/info", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
