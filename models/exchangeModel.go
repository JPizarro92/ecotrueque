package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Exchange struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Message        string         `json:"message"`
	Observations   string         `json:"observations"`
	Status         string         `json:"status" validate:"required"`
	ProposedUserID uuid.UUID      `json:"proposed_user_id" validate:"required"`
	ProposedUser   *User          `json:"user" gorm:"foreignKey:ProposedUserID"`
	ProposedPostID uint           `json:"proposed_post_id"`
	ProposedPost   *Post          `json:"proposed_post" gorm:"foreignKey:ProposedPostID"`
	PostUserID     uuid.UUID      `json:"post_user_id" validate:"required"`
	PostUser       *User          `json:"post_user" gorm:"foreignKey:PostUserID"`
	PostSharedID   uint           `json:"post_shared_id"`
	PostShared     *Post          `json:"post_shared" gorm:"foreignKey:PostSharedID"`
}
