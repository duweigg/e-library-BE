package tests

// func TestGetRecordList(t *testing.T) {
// 	db := SetupTestDB()
// 	router := SetupTestRouter(db)

// 	req, _ := http.NewRequest("POST", "/record/list", nil)
// 	req.Header.Set("Authorization", "Bearer valid_token")

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)
// }

// func TestExtendRecords(t *testing.T) {
// 	db := SetupTestDB()
// 	router := SetupTestRouter(db)

// 	requestBody, _ := json.Marshal(map[string]interface{}{
// 		"record_id": 456,
// 	})

// 	req, _ := http.NewRequest("POST", "/record/extend", bytes.NewBuffer(requestBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer valid_token")

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)
// }
