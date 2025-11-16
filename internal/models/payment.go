package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MidtransResponse represents JSON structure for Midtrans response
type MidtransResponse map[string]interface{}

// Value implements driver.Valuer interface for JSON
func (m MidtransResponse) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan implements sql.Scanner interface for JSON
func (m *MidtransResponse) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, m)
}

// MidtransAction represents JSON structure for Midtrans action
type MidtransAction map[string]interface{}

// Value implements driver.Valuer interface for JSON
func (m MidtransAction) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan implements sql.Scanner interface for JSON
func (m *MidtransAction) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, m)
}

// PaymentStatus enum
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusSuccess   PaymentStatus = "SUCCESS"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusCancelled PaymentStatus = "CANCELLED"
	PaymentStatusExpired   PaymentStatus = "EXPIRED"
)

type Payment struct {
	Id                    string            `json:"id" gorm:"primaryKey;type:varchar(36)"`
	OrderId               string            `json:"orderId" gorm:"type:varchar(100);uniqueIndex;not null"`
	UserId                string            `json:"userId" gorm:"type:varchar(36);not null"`
	Role                  *string           `json:"role,omitempty" gorm:"type:varchar(100)"`
	Amount                int               `json:"amount" gorm:"not null"`
	AdminFee              int               `json:"adminFee" gorm:"default:0"`
	TotalAmount           int               `json:"totalAmount" gorm:"not null"`
	PaymentMethod         string            `json:"paymentMethod" gorm:"type:varchar(50);not null"`
	PaymentType           *string           `json:"paymentType,omitempty" gorm:"type:varchar(50)"`
	Status                string            `json:"status" gorm:"type:varchar(20);default:'PENDING'"`
	SnapToken             *string           `json:"snapToken,omitempty" gorm:"type:text"`
	SnapRedirectUrl       *string           `json:"snapRedirectUrl,omitempty" gorm:"type:text"`
	TransactionId         *string           `json:"transactionId,omitempty" gorm:"type:varchar(255)"`
	MidtransTransactionId *string           `json:"midtransTransactionId,omitempty" gorm:"type:varchar(255)"`
	TransactionStatus     *string           `json:"transactionStatus,omitempty" gorm:"type:varchar(50)"`
	FraudStatus           *string           `json:"fraudStatus,omitempty" gorm:"type:varchar(50)"`
	PaymentCode           *string           `json:"paymentCode,omitempty" gorm:"type:varchar(100)"`
	VaNumber              *string           `json:"vaNumber,omitempty" gorm:"type:varchar(100)"`
	BankType              *string           `json:"bankType,omitempty" gorm:"type:varchar(50)"`
	ExpiryTime            *time.Time        `json:"expiryTime,omitempty" gorm:"type:timestamp"`
	PaidAt                *time.Time        `json:"paidAt,omitempty" gorm:"type:timestamp"`
	MidtransResponse      *MidtransResponse `json:"midtransResponse,omitempty" gorm:"type:jsonb"`
	MidtransAction        *MidtransAction   `json:"midtransAction,omitempty" gorm:"type:jsonb"`
	Star                  *int              `json:"star,omitempty"`
	Type                  *string           `json:"type,omitempty" gorm:"type:varchar(20)"`
	CreatedAt             time.Time         `json:"createdAt" gorm:"default:now()"`
	UpdatedAt             time.Time         `json:"updatedAt" gorm:"default:now()"`

	// Relations
	User    User  `json:"user,omitempty" gorm:"foreignKey:UserId;references:UserId"`
	RoleRef *Role `json:"roleRef,omitempty" gorm:"foreignKey:Role;references:Name"`
}

// TableName specifies the table name for GORM
func (Payment) TableName() string {
	return "payments"
}

// BeforeCreate hook to generate UUID if not set
func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.Id == "" {
		p.Id = uuid.New().String()
	}
	return nil
}

// BeforeUpdate hook to update UpdatedAt
func (p *Payment) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}
