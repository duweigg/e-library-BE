package models

import "time"

type Pagination struct {
	Page     int
	PageSize int
}

type CommonTime struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
