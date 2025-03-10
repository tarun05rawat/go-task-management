package database

import (
	"fmt"
	"log"

	"github.com/tarun05rawat/go-task-management/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {
	var err error

	// PostgreSQL connection string (update accordingly)
	dsn := "host=localhost user=tarunrawat dbname=task_management port=5432 sslmode=disable"

	// Open DB connection
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	// Check if connection is alive (Ping test)
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("❌ Failed to retrieve database instance:", err)
	}
	if err = sqlDB.Ping(); err != nil {
		log.Fatal("❌ Database connection not responding:", err)
	}

	// AutoMigrate models (Ensure all required tables exist)
	err = DB.AutoMigrate(&model.User{}, &model.UserData{}, &model.Task{}) // ✅ Added Task model
	if err != nil {
		log.Fatal("❌ Failed to auto-migrate database:", err)
	}

	fmt.Println("✅ Database connected and migrated successfully!")
}
