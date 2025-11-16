package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Benefit represents JSON structure for role benefits
type Benefit struct {
	Features []string `json:"features,omitempty"`
	Limits   map[string]interface{} `json:"limits,omitempty"`
}

// Value implements driver.Valuer interface for JSON
func (b Benefit) Value() (driver.Value, error) {
	return json.Marshal(b)
}

// Scan implements sql.Scanner interface for JSON
func (b *Benefit) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, b)
}

type Role struct {
	Id        string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name      string     `json:"name" gorm:"type:varchar(100);uniqueIndex;not null;index:idx_role_name"`
	Price     int        `json:"price" gorm:"not null"`
	Benefit   *Benefit   `json:"benefit,omitempty" gorm:"type:jsonb"`
	Image     *string    `json:"image,omitempty" gorm:"type:text"`
	CreatedAt time.Time  `json:"createdAt" gorm:"default:now()"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"default:now()"`

	// Relations
	Payments []Payment `json:"payments,omitempty" gorm:"foreignKey:Role;references:Name"`
}

// TableName specifies the table name for GORM
func (Role) TableName() string {
	return "roles"
}

// BeforeCreate hook to generate UUID if not set
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.Id == "" {
		r.Id = uuid.New().String()
	}
	return nil
}

// BeforeUpdate hook to update UpdatedAt
func (r *Role) BeforeUpdate(tx *gorm.DB) error {
	r.UpdatedAt = time.Now()
	return nil
}

