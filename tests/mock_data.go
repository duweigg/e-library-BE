package tests

import (
	"library/models"
	"time"
)

var MockUser = []models.User{
	{
		Username: "mock",
		Password: "$2a$10$Wx2A8AtGjiCBXia94By9V.fJPBsfyuQHwSblQg3fkPU.P5Ivt.tbe", // admin
		Nickname: "Mock",
	},
	{
		Username: "mock2",
		Password: "$2a$10$Wx2A8AtGjiCBXia94By9V.fJPBsfyuQHwSblQg3fkPU.P5Ivt.tbe", // admin
		Nickname: "Mock2",
	},
}

var MockBookType = []models.BookType{
	{
		Title: "Mock Book 1",
	},
	{
		Title: "Mock Book 2",
	},
	{
		Title: "Mock Book 3",
	},
}

var MockBook = []models.Book{
	{
		BookTypeID: 1,
		Status:     1,
	},
	{
		BookTypeID: 1,
		Status:     2,
	},
	{
		BookTypeID: 2,
		Status:     1,
	},
	{
		BookTypeID: 2,
		Status:     1,
	},
	{
		BookTypeID: 2,
		Status:     2,
	},
	{
		BookTypeID: 3,
		Status:     2,
	},
}

var layout = "2006-01-02 15:04:05" // Go's reference time format
var overdueAt = "2024-02-01 10:20:40"
var dueAt = "2025-04-01 10:20:40"
var returnedAt = "2023-04-01 10:20:40"
var parsedOverdueAt, _ = time.Parse(layout, overdueAt)
var parsedDueAt, _ = time.Parse(layout, dueAt)
var parsedReturnedAt, _ = time.Parse(layout, returnedAt)

var MockRecord = []models.Record{
	// overdue
	{
		UserID:     1,
		BookID:     2,
		ReturnedAt: nil,
		DueAt:      parsedOverdueAt,
		IsClosed:   false,
	},

	{
		UserID:     2,
		BookID:     5,
		ReturnedAt: nil,
		DueAt:      parsedDueAt,
		IsClosed:   false,
	},
	{
		UserID:     1,
		BookID:     6,
		ReturnedAt: nil,
		DueAt:      parsedDueAt,
		IsClosed:   false,
	},
	{
		UserID:     1,
		BookID:     1,
		ReturnedAt: &parsedReturnedAt,
		DueAt:      parsedDueAt,
		IsClosed:   true,
	},
}
