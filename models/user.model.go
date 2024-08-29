package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name            string     `gorm:"type:varchar(100);not null"`
	Email           string     `gorm:"type:varchar(100);uniqueIndex:idx_email;not null"`
	Password        string     `gorm:"type:varchar(100);not null"`
	Role            string     `gorm:"type:varchar(50);default:'user';not null"`
	DefaultLanguage string     `gorm:"type:varchar(50);default:'eng';not null"`
	DefaultAddress  string     `gorm:"type:varchar(50);"`
	Photo           string     `gorm:"type:varchar(100);not null;default:'default.png'"`
	IsEmailVerified bool       `gorm:"type:boolean;not null;default:false"`
	IsBlocked       bool       `gorm:"type:boolean;not null;default:false"`
	IsLoggedIn      bool       `gorm:"type:boolean;not null;default:false"`
	ValidLoggedInAt *time.Time `gorm:"type:timestamp;not null;default:now()"`
	CreatedAt       *time.Time `gorm:"type:timestamp;not null;default:now()"`
	UpdatedAt       *time.Time `gorm:"type:timestamp;not null;default:now()"`
}
