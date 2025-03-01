package main

import (
	"library/initializers"
	"library/models"

	"log"
)

func init() {
	initializers.GetEnvs()
	initializers.ConnectDB()
}

func main() {
	if initializers.DB == nil {
		log.Fatal("Database connection is nil")
	}

	err := initializers.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate User table:", err)
	}

	err = initializers.DB.AutoMigrate(&models.BookType{})
	if err != nil {
		log.Fatal("Failed to migrate Book table:", err)
	}

	err = initializers.DB.AutoMigrate(&models.Book{})
	if err != nil {
		log.Fatal("Failed to migrate Record table:", err)
	}

	err = initializers.DB.AutoMigrate(&models.Record{})
	if err != nil {
		log.Fatal("Failed to migrate Record table:", err)
	}

}

//go mod migrate/migrate.go
