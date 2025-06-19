// NOTE: Make sure to install gorm.io/datatypes if not present: go get gorm.io/datatypes
package models

import (
	"time"

	"gorm.io/datatypes"
)

// Only define Category struct once in the codebase. Remove this if already defined elsewhere.
// type Category struct {
// 	ID          uint   `gorm:"primaryKey"`
// 	Name        string `gorm:"unique;not null"`
// 	Description string
// 	CreatedAt   time.Time
// 	UpdatedAt   time.Time
// }

type LearningPlanStructure struct {
	ID         int64 `gorm:"primaryKey"`
	UserID     int64 // nullable for generic plans
	Goal       string
	TotalWeeks int
	Structure  datatypes.JSON // JSONB: milestones, weekly themes, prerequisites, etc.
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
}

type GeneratedContent struct {
	ID               int64 `gorm:"primaryKey"`
	PlanID           int64 // FK to LearningPlanStructure
	WeekNumber       int
	ContentData      datatypes.JSON // JSONB: detailed lessons, quizzes, etc.
	GeneratedBasedOn datatypes.JSON // JSONB: snapshot of user progress
	CreatedAt        time.Time      `gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime"`
}

type ContentAdaptationFlag struct {
	PlanID            int64
	WeekNumber        int
	NeedsRegeneration bool
	Reason            string
}
