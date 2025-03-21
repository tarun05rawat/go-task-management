package model

type TaskAttachment struct {
	ID       uint `gorm:"primaryKey"`
	TaskID   string
	FileURL  string
	Filename string
}
