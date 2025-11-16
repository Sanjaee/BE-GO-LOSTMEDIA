package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	PostId      string     `json:"postId" gorm:"primaryKey;type:varchar(36)"`
	UserId      string     `json:"userId" gorm:"type:varchar(36);not null;index:idx_post_user"`
	Title       string     `json:"title" gorm:"type:varchar(500);not null"`
	Description *string    `json:"description,omitempty" gorm:"type:text"`
	Content     *string    `json:"content,omitempty" gorm:"type:text"`
	Category    string     `json:"category" gorm:"type:varchar(100);not null"`
	MediaUrl    *string    `json:"mediaUrl,omitempty" gorm:"type:text"`
	Blurred     bool       `json:"blurred" gorm:"default:false"`
	ViewsCount  int        `json:"viewsCount" gorm:"default:0"`
	LikesCount  int        `json:"likesCount" gorm:"default:0"`
	SharesCount int        `json:"sharesCount" gorm:"default:0"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"default:now()"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"default:now()"`
	IsDeleted   bool       `json:"isDeleted" gorm:"default:false"`
	IsPublished bool       `json:"isPublished" gorm:"default:false"`
	ScheduledAt *time.Time `json:"scheduledAt,omitempty" gorm:"type:timestamp"`
	IsScheduled bool       `json:"isScheduled" gorm:"default:false"`

	// Relations
	User     User             `json:"user,omitempty" gorm:"foreignKey:UserId;references:UserId"`
	Comments []Comment        `json:"comments,omitempty" gorm:"foreignKey:PostId;references:PostId"`
	Likes    []Like           `json:"likes,omitempty" gorm:"foreignKey:PostId;references:PostId"`
	Sections []ContentSection `json:"sections,omitempty" gorm:"foreignKey:PostId;references:PostId"`
}

// TableName specifies the table name for GORM
func (Post) TableName() string {
	return "posts"
}

// BeforeCreate hook to generate UUID if not set
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.PostId == "" {
		p.PostId = uuid.New().String()
	}
	return nil
}

// BeforeUpdate hook to update UpdatedAt
func (p *Post) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}
