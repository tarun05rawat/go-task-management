package model

import (
	"time"
)

// Task struct
type Task struct {
	UserID      uint      `gorm:"primaryKey" json:"user_id"` // ✅ Composite Primary Key
	TaskID      uint      `gorm:"primaryKey" json:"id"`      // ✅ Composite Primary Key
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	Status      string    `gorm:"default:pending" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
