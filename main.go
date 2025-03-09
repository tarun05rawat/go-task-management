package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tarun05rawat/go-task-management/controllers"
	"github.com/tarun05rawat/go-task-management/database"
	"github.com/tarun05rawat/go-task-management/middleware"
)

func main() {
	// Connect to Database
	database.ConnectToDb()

	// Initialize Gin Router
	r := gin.Default()

	// Public Routes (No Authentication Required)
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)

	// Protected Routes (Require Authentication)
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth)
	protected.GET("/validate", controllers.Validate)

	// Start Server
	fmt.Println("✅ Server is running on port 8080")
	err := r.Run(":8080")
	if err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
