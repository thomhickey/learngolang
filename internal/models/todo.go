package models

import (
	"time"
)

type Todo struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string    `json:"title" validate:"required,max=100" gorm:"type:varchar(100);not null"`
	Description string    `json:"description" validate:"required,max=1000" gorm:"type:text;not null"`
	Completed   bool      `json:"completed" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime:true"`
}
