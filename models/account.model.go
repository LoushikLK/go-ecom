package models

import (
	"time"
)

type Account struct {
	ID                  string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name                string    `gorm:"type:varchar(100)" json:"name"`
	Email               string    `gorm:"type:varchar(50);unique;default:null" json:"email"`
	Mobile              string    `gorm:"type:varchar(30);unique;not null" json:"mobile"`
	Password            string    `gorm:"type:varchar(500)" json:"password"`
	IsBlocked           bool      `gorm:"type:boolean;default:false" json:"is_blocked"`
	IsBlacklisted       bool      `gorm:"type:boolean;default:false" json:"is_blacklisted"`
	IsBlockedReason     string    `gorm:"type:varchar(30)" json:"is_blocked_reason"`
	IsBlacklistedReason string    `gorm:"type:varchar(30)" json:"is_blacklisted_reason"`
	Lang                string    `gorm:"type:varchar(10);default:'en'" json:"lang"`
	CountryCode         string    `gorm:"type:varchar(10);default:'+91'" json:"country_code"`
	Lat                 float64   `gorm:"type:float" json:"lat"`
	Long                float64   `gorm:"type:float" json:"long"`
	IsLoggedIn          bool      `gorm:"type:boolean;default:false" json:"is_logged_in"`
	IsMobileVerified    bool      `gorm:"type:bool;default:false" json:"is_mobile_verified"`
	IsEmailVerified     bool      `gorm:"type:bool;default:false" json:"is_email_verified"`
	GoogleID            string    `gorm:"type:varchar(30)" json:"google_id"`
	AppleID             string    `gorm:"type:varchar(30)" json:"apple_id"`
	ProfileImage        string    `gorm:"type:varchar(30)" json:"profile_image"`
	CreatedAt           time.Time `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt           time.Time `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
}
