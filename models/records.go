package models

import "time"

type Record struct {
	ID         uint
	UserID     uint
	BookID     uint
	ReturnedAt time.Time
	DueAt      time.Time
	User       User `gorm:"foreignKey:UserID"` // Automatically fetch User
	Book       Book `gorm:"foreignKey:BookID"` // Automatically fetch Book
	IsClosed   bool `json:"is_closed" gorm:"column:is_closed;default:false"`
	CommonTime
}
type RecordResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	DueAt     time.Time `json:"due_at"`
	Status    string    `json:"status"`
}
type RecordRequest struct {
	IDs []uint `json:"ids"`
}

type RecordSearchRequest struct {
	Title  string `json:"title"`
	Status int    `json:"status"` //0: all, 1: open, 2: closed
	Pagination
}

func (r *Record) ToResponse() (rr RecordResponse) {
	var status = "Returned"
	if !r.IsClosed {
		if r.DueAt.Before(time.Now()) {
			status = "Overdue"
		} else {
			status = "Borrowed"
		}
	}

	rr = RecordResponse{
		ID:        r.ID,
		Name:      r.User.Nickname,
		Title:     r.Book.BookType.Title,
		CreatedAt: r.CreatedAt,
		DueAt:     r.DueAt,
		Status:    status,
	}
	return rr
}
