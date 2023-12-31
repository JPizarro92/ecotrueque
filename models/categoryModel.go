package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Name      string         `json:"name" validate:"required" gorm:"unique;not null"`
	Icon      string         `json:"icon"`
	Posts     []Post         `json:"posts" gorm:"foreignKey:CategoryID"`
}
