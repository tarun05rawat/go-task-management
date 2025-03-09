package model

import (
	"time"

	"gorm.io/gorm"
)

// Task struct
type Task struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"` // ✅ Auto-increment added
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	Status      string         `gorm:"default:pending" json:"status"` // ✅ Fixed default value
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
