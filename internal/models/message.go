package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	MessageId  string    `json:"messageId" gorm:"primaryKey;type:varchar(36)"`
	SenderId   string    `json:"senderId" gorm:"type:varchar(36);not null"`
	ReceiverId string    `json:"receiverId" gorm:"type:varchar(36);not null"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	MediaUrl   *string   `json:"mediaUrl,omitempty" gorm:"type:text"`
	IsRead     bool      `json:"isRead" gorm:"default:false"`
	CreatedAt  time.Time `json:"createdAt" gorm:"default:now()"`
	IsDeleted  bool      `json:"isDeleted" gorm:"default:false"`

	// Relations
	Sender   User `json:"sender,omitempty" gorm:"foreignKey:SenderId;references:UserId"`
	Receiver User `json:"receiver,omitempty" gorm:"foreignKey:ReceiverId;references:UserId"`
}

// TableName specifies the table name for GORM
func (Message) TableName() string {
	return "messages"
}

// BeforeCreate hook to generate UUID if not set
func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.MessageId == "" {
		m.MessageId = uuid.New().String()
	}
	return nil
}
