package models

import (
	"time"
)

type UserLogin struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID    string    `gorm:"type:uuid;not null uniqueIndex:idx_unique_device" json:"account_id"`
	Account      Account   `gorm:"foreignKey:AccountID;references:ID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE" json:"account"`
	FCM          string    `gorm:"type:varchar(50);uniqueIndex:idx_unique_device" json:"fcm"`
	DeviceName   string    `gorm:"type:varchar(30)" json:"device_name"`
	Lang         string    `gorm:"type:varchar(10);default:'en'" json:"lang"`
	RefreshToken string    `gorm:"type:varchar(500);uniqueIndex:idx_unique_device" json:"refresh_token"`
	Platform     string    `gorm:"type:varchar(30)" json:"platform"`
	IsExpired    time.Time `gorm:"type:timestamp" json:"is_expired"`
	IsActive     bool      `gorm:"type:bool;default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
}
