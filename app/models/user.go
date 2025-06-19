package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              uint           `gorm:"primaryKey"`
	Email           string         `gorm:"uniqueIndex;not null"`
	Password        string         `gorm:"not null"`
	FirstName       *string        `gorm:"default:null"`
	LastName        *string        `gorm:"default:null"`
	DailyCommitment int            `gorm:"not null"`
	LearningGoal    string         `gorm:"not null"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
