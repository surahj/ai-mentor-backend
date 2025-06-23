package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                int64          `gorm:"primaryKey"`
	Email             string         `gorm:"uniqueIndex;not null"`
	Password          *string        `gorm:"default:null"`
	FirstName         *string        `gorm:"default:null"`
	LastName          *string        `gorm:"default:null"`
	DailyCommitment   int            `gorm:"not null"`
	LearningGoal      string         `gorm:"not null"`
	Age               *int           `gorm:"default:null"`
	Level             *string        `gorm:"default:null"`
	Background        *string        `gorm:"default:null"`
	PreferredLanguage *string        `gorm:"default:null"`
	Interests         *string        `gorm:"default:null"`
	Country           *string        `gorm:"default:null"`
	CreatedAt         time.Time      `gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	IsVerified        bool           `gorm:"default:false" json:"is_verified"`
	OTP               *string        `json:"-"`
	OTPExpiresAt      *time.Time     `json:"-"`
	AuthProvider      string         `gorm:"default:'email'"` // 'email' or 'google'
}
