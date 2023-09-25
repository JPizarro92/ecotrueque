package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID          uint           `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description" validate:"required"`
	Status      string         `json:"status" validate:"required"`
	Image       string         `json:"image" validate:"required"`
}
