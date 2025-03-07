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
	dsn := "host=localhost user=tarunrawat dbname=task_management port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to Connect to database:", err)
	}

	err = DB.AutoMigrate(&model.User{}, &model.UserData{})
	if err != nil {
		log.Fatal("Failed to auto-migrate database: ", err)
	}

	fmt.Println("Database connected and migrated successfully!")

}
