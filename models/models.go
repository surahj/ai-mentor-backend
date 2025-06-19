package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email             string         `gorm:"unique;not null" json:"email"`
	Password          string         `gorm:"not null" json:"-"`
	Name              string         `gorm:"not null" json:"name"`
	DailyCommitment   int            `gorm:"not null" json:"daily_commitment"`
	LearningGoal      string         `gorm:"not null" json:"learning_goal"`
	ResetToken        string         `json:"reset_token"`
	ResetTokenExpires time.Time      `json:"reset_token_expires"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

type Session struct {
	gorm.Model
	UserID          uint      `json:"user_id"`
	Topic           string    `json:"topic"`
	DurationMinutes int       `json:"duration_minutes"`
	Notes           string    `json:"notes"`
	Date            time.Time `json:"date"`
	Rating          int       `json:"rating"`     // 1-5 rating for productivity/quality
	Tags            string    `json:"tags"`       // Comma-separated tags
	Reflection      string    `json:"reflection"` // Post-session reflection
}

type Goal struct {
	gorm.Model
	UserID        uint
	Title         string
	Description   string
	TargetMinutes int  // Minutes per day/week
	IsWeekly      bool // true = weekly goal, false = daily goal
	StartDate     time.Time
	EndDate       time.Time // Optional for recurring goals
	IsRecurring   bool      // Does goal repeat?
}
