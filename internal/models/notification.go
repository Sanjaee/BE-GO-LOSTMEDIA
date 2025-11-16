package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	NotifId   string    `json:"notifId" gorm:"primaryKey;type:varchar(36)"`
	UserId    *string   `json:"userId,omitempty" gorm:"type:varchar(36)"`
	ActorId   string    `json:"actorId" gorm:"type:varchar(36);not null"`
	Type      string    `json:"type" gorm:"type:varchar(50);not null"`
	Content   *string   `json:"content,omitempty" gorm:"type:text"`
	ActionUrl *string   `json:"actionUrl,omitempty" gorm:"type:text"`
	IsRead    bool      `json:"isRead" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:now()"`

	// Relations
	User  *User `json:"user,omitempty" gorm:"foreignKey:UserId;references:UserId"`
	Actor User  `json:"actor,omitempty" gorm:"foreignKey:ActorId;references:UserId"`
}

// TableName specifies the table name for GORM
func (Notification) TableName() string {
	return "notifications"
}

// BeforeCreate hook to generate UUID if not set
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.NotifId == "" {
		n.NotifId = uuid.New().String()
	}
	return nil
}
