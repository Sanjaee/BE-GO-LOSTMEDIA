package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Follower struct {
	FollowId    string    `json:"followId" gorm:"primaryKey;type:varchar(36)"`
	FollowerId  string    `json:"followerId" gorm:"type:varchar(36);not null"`
	FollowingId string    `json:"followingId" gorm:"type:varchar(36);not null"`
	FollowedAt  time.Time `json:"followedAt" gorm:"default:now()"`
	IsActive    bool      `json:"isActive" gorm:"default:true"`

	// Relations
	Follower  User `json:"follower,omitempty" gorm:"foreignKey:FollowerId;references:UserId"`
	Following User `json:"following,omitempty" gorm:"foreignKey:FollowingId;references:UserId"`
}

// TableName specifies the table name for GORM
func (Follower) TableName() string {
	return "followers"
}

// BeforeCreate hook to generate UUID if not set
func (f *Follower) BeforeCreate(tx *gorm.DB) error {
	if f.FollowId == "" {
		f.FollowId = uuid.New().String()
	}
	return nil
}
