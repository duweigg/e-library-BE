package models

type BookType struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Title string `json:"title"`
	CommonTime
}

type Book struct {
	ID         uint     `json:"id" gorm:"primary_key"`
	BookTypeID uint     `josn:"book_type_id"`
	Status     uint     `json:"status"` //1: avaiable, 2: rent out
	BookType   BookType `gorm:"foreignKey:BookTypeID"`
	CommonTime
}

type BookResponse struct {
	ID             uint   `json:"id" gorm:"primary_key"`
	Name           string `json:"name"`
	TotalCount     int    `json:"total_count"`
	AvailableCount int    `json:"available_count"`
}
type BookIDsPayload struct {
	BookTypeIDs []uint `json:"ids"`
}
type BookRequest struct {
	Title string `json:"title"`
	Pagination
}
