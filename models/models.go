package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	Email           string         `gorm:"uniqueIndex;not null" json:"email"`
	Password        string         `gorm:"not null" json:"-"`
	Name            string         `gorm:"not null" json:"name"`
	DailyCommitment int            `gorm:"not null" json:"daily_commitment"`
	LearningGoal    string         `gorm:"not null" json:"learning_goal"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
