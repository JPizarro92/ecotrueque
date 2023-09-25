package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FormatDate   = "2006-01-02"
type User struct {
	ID                uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Name              string         `json:"name" validate:"required,min=2,max=100"`
	Surname           string         `json:"surname" validate:"required,min=2,max=100"`
	Email             string         `json:"email" binding:"required,email" gorm:"unique;not null"`
	Password          string         `json:"password" validate:"required,min=6"`
	Phone             string         `json:"phone"`
	Location          string         `json:"location"`
	Token             string         `json:"token"`
	Refresh_token     string         `json:"refresh_token"`
	Verified          bool           `json:"verified"`
	VerificationCode  bool           `json:"verification_code"`
	Role              string         `json:"role" validate:"required,eq=ADMIN|eq=USER"`
	Avatar            string         `json:"avatar"`
	BirthDate         time.Time      `json:"birth_date"`
	Age               uint8          `json:"age"`
	Status            string         `json:"status" validate:"required"`
	RatingScore       float64        `json:"rating_score"`
	Ratings           []Rating       `json:"ratings" gorm:"foreignKey:RatedID"`
	Posts             []Post         `json:"posts" gorm:"foreignKey:UserID"`
	ExchangesSend     []Exchange     `json:"exchanges_send" gorm:"foreignKey:ProposedUserID"`
	ExchangesReceived []Exchange     `json:"exchanges_received" gorm:"foreignKey:PostUserID"`
	//ExchangesAccepted []Exchange     `json:"exchanges_accepted" gorm:"foreignkey:UserID"`
}

type Rating struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Value       float32        `json:"value" validate:"required"`
	RatedByID   uuid.UUID      `json:"rated_by_id" validate:"required"`
	RatedByUser *User          `json:"rated_by_ user" gorm:"foreignKey:RatedByID" `
	RatedID     uuid.UUID      `json:"rated_id" validate:"required"`
	RatedUser   *User          `json:"rated_user" gorm:"foreignKey:RatedID"`
}

type SignUpInput struct {
	Name            string `json:"name" binding:"required"`
	Surname         string `json:"surname" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
	Avatar          string `json:"avatar" binding:"required"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	Avatar    string    `json:"avatar,omitempty"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
