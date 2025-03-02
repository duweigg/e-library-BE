package controllers

import (
	"library/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Define a struct to hold the database instance
type UserController struct {
	DB *gorm.DB
}

// Constructor function to create a new BookController
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var signUpPayload models.SignUpPayload

	// Validate request payload
	if err := c.ShouldBindJSON(&signUpPayload); err != nil {
		log.Printf("Invalid signup request: %v\n", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	// Check if the username already exists
	var existingUser models.User
	if err := uc.DB.Where("username = ?", signUpPayload.Username).First(&existingUser).Error; err == nil {
		log.Printf("Username already in use: %s\n", signUpPayload.Username)
		c.JSON(http.StatusConflict, gin.H{"error": "Username already in use"})
		return
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(signUpPayload.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Nickname: signUpPayload.Nickname,
		Username: signUpPayload.Username,
		Password: string(passwordHash),
	}

	// Use a transaction for safety
	tx := uc.DB.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to create user: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	tx.Commit()

	log.Printf("User %s created successfully\n", user.Username)
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "data": user})
}

func (uc *UserController) SignIn(c *gin.Context) {
	var signInPayload models.SignInPayload

	// Validate request payload
	if err := c.ShouldBindJSON(&signInPayload); err != nil {
		log.Printf("Invalid signin request: %v\n", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	// Find user by username
	var userFound models.User
	if err := uc.DB.Where("username = ?", signInPayload.Username).First(&userFound).Error; err != nil {
		log.Printf("User not found: %s\n", signInPayload.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(signInPayload.Password)); err != nil {
		log.Printf("Invalid password for user: %s\n", signInPayload.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	token, err := generateJWT(userFound.ID)
	if err != nil {
		log.Printf("Failed to generate token for user: %s\n", signInPayload.Username)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	log.Printf("User %s signed in successfully\n", userFound.Username)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (uc *UserController) GetUserInfo(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// generateJWT creates a JWT token
func generateJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Println("JWT secret is not set")
		return "", http.ErrServerClosed
	}

	return token.SignedString([]byte(secret))
}
