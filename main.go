package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tarun05rawat/go-task-management/controllers"
	"github.com/tarun05rawat/go-task-management/database"
	"github.com/tarun05rawat/go-task-management/handlers"
	"github.com/tarun05rawat/go-task-management/middleware"
)

func main() {
	// ✅ Connect to Database
	database.ConnectToDb()

	// ✅ Initialize Gin Router
	r := gin.Default()

	// ✅ Public Routes (No Authentication Required)
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)

	// ✅ Protected Routes (Require Authentication)
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth)

	// ✅ Authentication validation
	protected.GET("/validate", controllers.Validate)

	// ✅ Admin-only Route to View All Users
	protected.GET("/users", controllers.GetAllUsers)

	// ✅ Task Management Routes (For Authenticated Users)
	taskRoutes := protected.Group("/tasks") // ✅ This groups all task routes under `/tasks`
	{
		taskRoutes.POST("/", handlers.CreateTask)      // ✅ Create Task
		taskRoutes.GET("/", handlers.GetTasks)         // ✅ Get All Tasks (for logged-in user)
		taskRoutes.GET("/:id", handlers.GetTaskByID)   // ✅ Get Specific Task
		taskRoutes.PUT("/:id", handlers.UpdateTask)    // ✅ Update Task
		taskRoutes.DELETE("/:id", handlers.DeleteTask) // ✅ Delete Task
	}

	// ✅ Start Server
	fmt.Println("✅ Server is running on port 8080")
	err := r.Run(":8080")
	if err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
