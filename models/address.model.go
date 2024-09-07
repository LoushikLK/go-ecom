package models

import "time"

type Address struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID       string    `gorm:"type:uuid;not null;unique" json:"account_id"`
	Account         Account   `gorm:"foreignKey:AccountID;references:ID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE" json:"account"`
	IsDefault       bool      `gorm:"default:false" json:"is_default"`
	IsDeleted       bool      `gorm:"default:false" json:"is_deleted"`
	FullName        string    `gorm:"type:varchar(255);not null" json:"full_name"`
	PhoneNumber     string    `gorm:"type:varchar(255);not null" json:"phone_number"`
	AddressLine1    string    `gorm:"type:varchar(255);not null" json:"address_line1"`
	AddressLine2    string    `gorm:"type:varchar(255)" json:"address_line2"`
	City            string    `gorm:"type:varchar(255);not null" json:"city"`
	State           string    `gorm:"type:varchar(255);not null" json:"state"`
	Country         string    `gorm:"type:varchar(255);not null" json:"country"`
	PostalCode      string    `gorm:"type:varchar(255);not null" json:"postal_code"`
	ShippingAddress bool      `gorm:"default:false" json:"shipping_address"`
	BillingAddress  bool      `gorm:"default:false" json:"billing_address"`
	AddressTitle    string    `gorm:"type:varchar(255);" json:"address_title"`
	CreatedAt       time.Time `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt       time.Time `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
}
