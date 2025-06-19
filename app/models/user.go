package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                int64    `gorm:"primaryKey"`
	Email             string  `gorm:"uniqueIndex;not null"`
	Password          string  `gorm:"not null"`
	FirstName         *string `gorm:"default:null"`
	LastName          *string `gorm:"default:null"`
	DailyCommitment   int     `gorm:"not null"`
	LearningGoal      string  `gorm:"not null"`
	Age               *int
	Level             *string        // e.g., "beginner", "intermediate", "advanced"
	Background        *string        // e.g., "Student", "Software Engineer"
	PreferredLanguage *string        // e.g., "English"
	Interests         *string        // or pq.StringArray/JSON for multiple
	CreatedAt         time.Time      `gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}
