package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ImageDetail represents JSON structure for image details
type ImageDetail struct {
	Width   int    `json:"width,omitempty"`
	Height  int    `json:"height,omitempty"`
	Alt     string `json:"alt,omitempty"`
	Caption string `json:"caption,omitempty"`
}

// Value implements driver.Valuer interface for JSON
func (i ImageDetail) Value() (driver.Value, error) {
	return json.Marshal(i)
}

// Scan implements sql.Scanner interface for JSON
func (i *ImageDetail) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, i)
}

type ContentSection struct {
	SectionId   string       `json:"sectionId" gorm:"primaryKey;type:varchar(36)"`
	Type        string       `json:"type" gorm:"type:varchar(50);not null"`
	Content     *string      `json:"content,omitempty" gorm:"type:text"`
	Src         *string      `json:"src,omitempty" gorm:"type:text"`
	ImageDetail *ImageDetail `json:"imageDetail,omitempty" gorm:"type:jsonb"`
	Order       int          `json:"order" gorm:"not null"`
	PostId      string       `json:"postId" gorm:"type:varchar(36);not null"`
	CreatedAt   time.Time    `json:"createdAt" gorm:"default:now()"`
	UpdatedAt   time.Time    `json:"updatedAt" gorm:"default:now()"`

	// Relations
	Post Post `json:"post,omitempty" gorm:"foreignKey:PostId;references:PostId"`
}

// TableName specifies the table name for GORM
func (ContentSection) TableName() string {
	return "content_sections"
}

// BeforeCreate hook to generate UUID if not set
func (c *ContentSection) BeforeCreate(tx *gorm.DB) error {
	if c.SectionId == "" {
		c.SectionId = uuid.New().String()
	}
	return nil
}

// BeforeUpdate hook to update UpdatedAt
func (c *ContentSection) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}
