package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	CommentId  string    `json:"commentId" gorm:"primaryKey;type:varchar(36)"`
	PostId     string    `json:"postId" gorm:"type:varchar(36);not null"`
	UserId     string    `json:"userId" gorm:"type:varchar(36);not null"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	LikesCount int       `json:"likesCount" gorm:"default:0"`
	CreatedAt  time.Time `json:"createdAt" gorm:"default:now()"`
	IsDeleted  bool      `json:"isDeleted" gorm:"default:false"`
	ParentId   *string   `json:"parentId,omitempty" gorm:"type:varchar(36)"`

	// Relations
	Post    Post      `json:"post,omitempty" gorm:"foreignKey:PostId;references:PostId"`
	User    User      `json:"user,omitempty" gorm:"foreignKey:UserId;references:UserId"`
	Parent  *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentId;references:CommentId"`
	Replies []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentId;references:CommentId"`
	Likes   []Like    `json:"likes,omitempty" gorm:"foreignKey:CommentId;references:CommentId"`
}

// TableName specifies the table name for GORM
func (Comment) TableName() string {
	return "comments"
}

// BeforeCreate hook to generate UUID if not set
func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.CommentId == "" {
		c.CommentId = uuid.New().String()
	}
	return nil
}
