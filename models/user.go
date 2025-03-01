package models

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Nickname string
	CommonTime
}

type SignInPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignUpPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

type UserResponse struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Nickname string `json:"nickname"`
}
