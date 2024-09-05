package models

import (
	"time"
)

type UserOtp struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID       string    `gorm:"type:uuid;not null;unique" json:"account_id"`
	Account         Account   `gorm:"foreignKey:AccountID;references:ID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE" json:"account"`
	Otp             string    `gorm:"type:varchar(10);not null" json:"otp"`
	IsExpired       bool      `gorm:"type:boolean;default:false" json:"is_expired"`
	ExpiredDateTime time.Time `gorm:"type:timestamp" json:"expired_date_time"`
	CreatedAt       time.Time `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt       time.Time `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
}
