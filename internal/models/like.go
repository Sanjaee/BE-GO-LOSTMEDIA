package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Like struct {
	LikeId    string    `json:"likeId" gorm:"primaryKey;type:varchar(36)"`
	UserId    string    `json:"userId" gorm:"type:varchar(36);not null"`
	PostId    *string   `json:"postId,omitempty" gorm:"type:varchar(36)"`
	CommentId *string   `json:"commentId,omitempty" gorm:"type:varchar(36)"`
	LikeType  string    `json:"likeType" gorm:"type:varchar(20);not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:now()"`

	// Relations
	User    User     `json:"user,omitempty" gorm:"foreignKey:UserId;references:UserId"`
	Post    *Post    `json:"post,omitempty" gorm:"foreignKey:PostId;references:PostId"`
	Comment *Comment `json:"comment,omitempty" gorm:"foreignKey:CommentId;references:CommentId"`
}

// TableName specifies the table name for GORM
func (Like) TableName() string {
	return "likes"
}

// BeforeCreate hook to generate UUID if not set
func (l *Like) BeforeCreate(tx *gorm.DB) error {
	if l.LikeId == "" {
		l.LikeId = uuid.New().String()
	}
	return nil
}
