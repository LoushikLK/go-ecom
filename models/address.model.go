package models

import (
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID          *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name        string     `gorm:"type:varchar(100);not null"`
	PhoneNumber string     `gorm:"type:varchar(100);not null"`
	Email       string     `gorm:"type:varchar(100);not null"`
	UserId      *uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt   *time.Time `gorm:"not null;default:now()"`
	UpdatedAt   *time.Time `gorm:"not null;default:now()"`
}
