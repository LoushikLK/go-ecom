package models

import "time"

type Address struct {
	ID                string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	AccountID         string    `gorm:"type:uuid;not null;" json:"account_id"`
	Account           Account   `gorm:"foreignKey:AccountID;references:ID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE" json:"account"`
	IsDefault         bool      `gorm:"default:false" json:"is_default"`
	IsDeleted         bool      `gorm:"default:false" json:"is_deleted"`
	FullName          string    `gorm:"type:varchar(255);not null" json:"full_name"`
	PhoneNumber       string    `gorm:"type:varchar(255);not null" json:"phone_number"`
	CountryCode       string    `gorm:"type:varchar(50);not null" json:"country_code"`
	AddressLine1      string    `gorm:"type:varchar(255);not null" json:"address_line1"`
	AddressLine2      string    `gorm:"type:varchar(255)" json:"address_line2"`
	City              string    `gorm:"type:varchar(100);not null" json:"city"`
	State             string    `gorm:"type:varchar(100);not null" json:"state"`
	Country           string    `gorm:"type:varchar(100);not null" json:"country"`
	PostalCode        string    `gorm:"type:varchar(50);not null" json:"postal_code"`
	IsShippingAddress bool      `gorm:"default:true" json:"is_shipping_address"`
	IsBillingAddress  bool      `gorm:"default:true" json:"is_billing_address"`
	AddressTitle      string    `gorm:"type:varchar(255);" json:"address_title"`
	Lat               float64   `gorm:"type:float" json:"lat"`
	Long              float64   `gorm:"type:float" json:"long"`
	DeletedAt         time.Time `gorm:"type:timestamp" json:"deleted_at"`
	CreatedAt         time.Time `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt         time.Time `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
}
