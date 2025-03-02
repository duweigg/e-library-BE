package tests

import (
	"library/controllers"
	"library/models"
	"log"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupMockDB initializes a mock PostgreSQL database using pgxmock
func SetupMockDB() *gorm.DB {
	DB, err := gorm.Open(postgres.Open("host=localhost user=postgres password=admin dbname=library_test port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	return DB
}

func SetupMockRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	bookController := controllers.NewBookController(db)

	bookRouter := router.Group("/book")
	{
		bookRouter.POST("/list", bookController.GetBookList)
		bookRouter.POST("/borrow", MockCheckAuth, bookController.BorrowBooks)
	}

	return router
}

// TestMain runs before any test starts
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func MockCheckAuth(c *gin.Context) {
	var user = models.User{
		ID:       1,
		Username: "test",
		Password: "123",
		Nickname: "Test",
	}
	c.Set("user", user)
}
