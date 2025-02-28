package main

import (
	"library/initializers"
	"library/models"
)

func init() {
	initializers.GetEnvs()
	initializers.ConnectDB()
}

func main() {
     initializers.DB.AutoMigrate(&models.User{})
}

//go mod migrate/migrate.go