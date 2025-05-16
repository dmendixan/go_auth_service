package main

import (
	"auth-service/config"
	"auth-service/internal/handlers"
	"auth-service/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	config.InitDB()

	// ✅ Автоматическая миграция таблиц
	if err := config.DB.AutoMigrate(&models.User{}, &models.RefreshToken{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	r := gin.Default()
	api := r.Group("/api")

	api.POST("/register", handlers.RegisterWithDB(config.DB))
	api.POST("/login", handlers.Login)
	api.POST("/refresh", handlers.Refresh)

	api.GET("/profile", handlers.AuthMiddleware(""), handlers.ProfileHandler(config.DB))

	admin := api.Group("/admin")
	admin.Use(handlers.AuthMiddleware("admin")) // Только админ
	{
		admin.GET("/users", handlers.GetAllUsersWithDB(config.DB))
		admin.DELETE("/users/:id", handlers.DeleteUserWithDB(config.DB))
		admin.PUT("/users/:id", handlers.UpdateUserWithDB(config.DB))
	}

	_ = r.Run(":8080")
	fmt.Println("hello")
}
