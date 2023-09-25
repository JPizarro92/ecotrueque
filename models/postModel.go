package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Title            string         `json:"title" validate:"required,min=2,max=250"`
	Price            float64        `json:"price"`
	ShortDescription string         `json:"short_description" validate:"required"`
	LongDescription  string         `json:"long_description"`
	ExchangeRate     string         `json:"exchange_rate"`
	Tags             string         `json:"tags"`
	ProductStatus    string         `json:"product_status" validate:"required"`
	PostStatus       bool           `json:"post_status" validate:"required"`
	UserID           uuid.UUID      `json:"user_id" validate:"required"`
	User             *User          `json:"user" gorm:"foreignKey:UserID"`
	CategoryID       uint           `json:"category_id" validate:"required"`
	Category         *Category      `json:"category" gorm:"foreignKey:CategoryID"`
	SubCategoryID    *uint          `json:"sub_category_id"`
	Images           []PostImage    `json:"images" gorm:"foreignKey:PostID"`
}

type PostImage struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Img       string         `json:"img"`
	PostID    uint           `json:"post_id"`
}
