package tests

import (
	"database/sql"
	"library/controllers"
	"library/initializers"
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

	userController := controllers.NewUserController(db)
	userRouter := router.Group("/user")
	{
		userRouter.POST("/signup", userController.CreateUser)
		userRouter.POST("/signin", userController.SignIn)
		userRouter.GET("/info", MockCheckAuth, userController.GetUserInfo)
	}

	bookController := controllers.NewBookController(db)
	bookRouter := router.Group("/book")
	{
		bookRouter.POST("/list", bookController.GetBookList)
		bookRouter.POST("/borrow", MockCheckAuth, bookController.BorrowBooks)
	}

	recordController := controllers.NewRecordController(db)
	recordRouter := router.Group("/record")
	{
		recordRouter.POST("/list", MockCheckAuth, recordController.GetRecordList)
		recordRouter.POST("/extend", MockCheckAuth, recordController.ExtendRecords)
		recordRouter.POST("/return", MockCheckAuth, recordController.ReturnRecords)
	}
	return router
}

// TestMain runs before any test starts
func TestMain(m *testing.M) {
	initializers.GetEnvs()
	gin.SetMode(gin.TestMode)
	db := SetupMockDB()
	PrepareMockUserDB(db)
	defer db.ConnPool.(*sql.DB).Close()
	m.Run()
}

func MockCheckAuth(c *gin.Context) {
	var user = models.UserResponse{
		ID:       1,
		Nickname: "Test",
	}
	c.Set("user", user)
}

func PrepareMockUserDB(db *gorm.DB) {
	db.Migrator().DropTable(&models.User{})
	db.Migrator().AutoMigrate(&models.User{})
	db.Save(&MockUser)

}
func PrepareMockBookDB(db *gorm.DB) {
	db.Migrator().DropTable(&models.Book{})
	db.Migrator().DropTable(&models.BookType{})
	db.Migrator().AutoMigrate(&models.Book{})
	db.Migrator().AutoMigrate(&models.BookType{})
	db.Save(&MockBookType)
	db.Save(&MockBook)
}
func PrepareMockRecordDB(db *gorm.DB) {
	db.Migrator().DropTable(&models.Record{})
	db.Migrator().DropTable(&models.Book{})
	db.Migrator().DropTable(&models.BookType{})
	db.Migrator().AutoMigrate(&models.Record{})
	db.Migrator().AutoMigrate(&models.Book{})
	db.Migrator().AutoMigrate(&models.BookType{})
	db.Save(&MockBookType)
	db.Save(&MockBook)
	db.Save(&MockRecord)
}
