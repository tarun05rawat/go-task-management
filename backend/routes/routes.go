package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tarun05rawat/go-task-management/handlers"
	"github.com/tarun05rawat/go-task-management/middleware"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")

	// ✅ Secure Route (Admin Only)
	protected := api.Group("/")
	protected.Use(middleware.RequireAuth)
	protected.GET("/users", handlers.GetAllUsers) // ✅ Only Admins can view all users
}
