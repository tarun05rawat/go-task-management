package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tarun05rawat/go-task-management/controllers"
	"github.com/tarun05rawat/go-task-management/database"
)

func main() {
	database.ConnectToDb()

	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	fmt.Println("Server is running on port 8080") // Ensure this prints
	err := r.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
