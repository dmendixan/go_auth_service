package main

import (
	"auth-service/config"
	"auth-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()

	r := gin.Default()
	api := r.Group("/api")

	api.POST("/register", handlers.Register)
	api.POST("/login", handlers.Login)
	api.POST("/refresh", handlers.Refresh)

	api.GET("/profile", handlers.AuthMiddleware(""), handlers.Profile)

	admin := api.Group("/admin")
	admin.Use(handlers.AuthMiddleware("admin")) // Только админ
	{
		admin.GET("/users", handlers.GetAllUsers)
		admin.DELETE("/users/:id", handlers.DeleteUser)
		admin.PUT("/users/:id", handlers.UpdateUser)
	}

	_ = r.Run(":8080")

}
