package models

import (
	"time"
)

type Book struct {
	ID        		uint       `json:"id" gorm:"primary_key"`
	Name      		string     `json:"name"`	
	Status    		uint	   // 1: available, 2: rent
	AvailableDate   time.Time  
	UserID    		string     `json:"user_id" gorm:"index"` // Foreign key, indexed for performance
	User      		User       `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt 		time.Time
	UpdatedAt 		time.Time
}

type BookPayload struct {
	ID []uint
}