package main

import (
	"auth-service/config"
	"auth-service/internal/models"
	"log"
)

func main() {
	config.InitDB()

	err := config.DB.AutoMigrate(&models.User{}, &models.RefreshToken{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migration completed successfully")
}
