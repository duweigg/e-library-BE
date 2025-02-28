package models

import (
	"time"
)

type Book struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`	
	UserID    string    `json:"user_id" gorm:"index"` // Foreign key, indexed for performance
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
