package models

import "time"

type Category struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique;not null"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type LearningPath struct {
	ID              uint   `gorm:"primaryKey"`
	Name            string `gorm:"unique;not null"`
	Description     string
	CategoryID      uint
	DailyCommitment int    `gorm:"not null"`
	PlanJSON        string `gorm:"type:jsonb"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
