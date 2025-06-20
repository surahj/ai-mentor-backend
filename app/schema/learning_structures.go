package schema

import (
	"time"

	"gorm.io/datatypes"
)

// LearningPlanStructure represents the high-level structure of a learning plan
type LearningPlanStructure struct {
	BaseModel
	UserID     int64 // nullable for generic plans
	Goal       string
	TotalWeeks int
	Structure  datatypes.JSON // JSONB: stores the complete structure
}

// WeeklyTheme represents a week's learning theme and objectives
type WeeklyTheme struct {
	WeekNumber    int      `json:"week_number"`
	Theme         string   `json:"theme"`
	Objectives    []string `json:"objectives"`
	KeyConcepts   []string `json:"key_concepts"`
	Prerequisites []string `json:"prerequisites"`
}

// DailyMilestone represents a daily learning milestone
type DailyMilestone struct {
	BaseModel
	DayNumber   int    `json:"day_number"`
	Topic       string `json:"topic"`
	Description string `json:"description"`
	Duration    int    `json:"duration_minutes"` // in minutes
	Difficulty  string `json:"difficulty"`       // beginner, intermediate, advanced
}

// Exercise represents a practice exercise or quiz
type Exercise struct {
	BaseModel
	Type        string   `json:"type"` // quiz, coding, reading, etc.
	Question    string   `json:"question"`
	Options     []string `json:"options,omitempty"` // for multiple choice
	Answer      string   `json:"answer"`
	Explanation string   `json:"explanation"`
	Difficulty  string   `json:"difficulty"`
}

// Resource represents learning resources (videos, articles, etc.)
type Resource struct {
	BaseModel
	Type        string `json:"type"` // video, article, book, etc.
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Duration    int    `json:"duration_minutes,omitempty"`
}

// WeeklyContent represents detailed content for a specific week
type WeeklyContent struct {
	ID              int64            `gorm:"primaryKey"`
	WeekNumber      int              `json:"week_number"`
	Theme           WeeklyTheme      `json:"theme"`
	DailyMilestones []DailyMilestone `json:"daily_milestones"`
	Exercises       []Exercise       `json:"exercises"`
	Resources       []Resource       `json:"resources"`
	AdaptiveNotes   string           `json:"adaptive_notes"` // AI notes based on user progress
}

// CompleteLearningPlan represents the full structure
type CompleteLearningPlan struct {
	ID              int64               `gorm:"primaryKey"`
	Goal            string              `json:"goal"`
	TotalWeeks      int                 `json:"total_weeks"`
	DailyCommitment int                 `json:"daily_commitment_minutes"`
	WeeklyThemes    []WeeklyTheme       `json:"weekly_themes"`
	Prerequisites   map[string][]string `json:"prerequisites"`
	AdaptiveRules   map[string]string   `json:"adaptive_rules"` // rules for content adaptation
}

// GeneratedContent represents the detailed content stored in the database
type GeneratedContent struct {
	ID               int64 `gorm:"primaryKey"`
	PlanID           int64 // FK to LearningPlanStructure
	WeekNumber       int
	ContentData      datatypes.JSON // JSONB: stores int64
	GeneratedBasedOn datatypes.JSON // JSONB: snapshot of user progress
	CreatedAt        time.Time
}

// ContentAdaptationFlag represents flags for content regeneration
type ContentAdaptationFlag struct {
	PlanID            uint
	WeekNumber        int
	NeedsRegeneration bool
	Reason            string
}
