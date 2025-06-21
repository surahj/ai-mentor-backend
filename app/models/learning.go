package models

import (
	"time"

	"gorm.io/datatypes"
)

// LearningPlanStructure represents the high-level structure of a learning plan
type LearningPlanStructure struct {
	BaseModel
	UserID     int64          `json:"user_id" example:"1"`
	Goal       string         `json:"goal" example:"Learn React and TypeScript"`
	TotalWeeks int            `json:"total_weeks" example:"8"`
	Structure  datatypes.JSON `json:"structure" swaggertype:"object"` // JSONB: stores the complete structure
}

// WeeklyTheme represents a week's learning theme and objectives
type WeeklyTheme struct {
	WeekNumber    int      `json:"week_number" example:"1"`
	Theme         string   `json:"theme" example:"React Fundamentals"`
	Objectives    []string `json:"objectives" example:"[\"Understand JSX\",\"Learn Components\"]"`
	KeyConcepts   []string `json:"key_concepts" example:"[\"JSX\",\"Components\",\"Props\"]"`
	Prerequisites []string `json:"prerequisites" example:"[\"JavaScript Basics\"]"`
}

// DailyMilestone represents a daily learning milestone
type DailyMilestone struct {
	BaseModel
	DayNumber       int    `json:"day_number" example:"1"`
	Topic           string `json:"topic" example:"Introduction to Go"`
	Description     string `json:"description" example:"Learn the basics of Go syntax."`
	DurationMinutes int    `json:"duration_minutes" example:"60"`
	Difficulty      string `json:"difficulty" example:"Easy"`
}

// Exercise represents a practice exercise or quiz
type Exercise struct {
	BaseModel
	Type        string   `json:"type" example:"quiz"` // quiz, coding, reading, etc.
	Question    string   `json:"question" example:"What is JSX?"`
	Options     []string `json:"options,omitempty" example:"[\"JavaScript XML\",\"React Syntax\",\"HTML in JS\"]"` // for multiple choice
	Answer      string   `json:"answer" example:"JavaScript XML"`
	Explanation string   `json:"explanation" example:"JSX stands for JavaScript XML"`
	Difficulty  string   `json:"difficulty" example:"beginner"`
}

// Resource represents learning resources (videos, articles, etc.)
type Resource struct {
	BaseModel
	Type        string `json:"type" example:"video"` // video, article, book, etc.
	Title       string `json:"title" example:"React Tutorial"`
	URL         string `json:"url" example:"https://example.com/react-tutorial"`
	Description string `json:"description" example:"Complete React tutorial for beginners"`
	Duration    int    `json:"duration_minutes,omitempty" example:"60"`
}

// WeeklyContent represents the structure of content generated for a week
type WeeklyContent struct {
	Theme           string           `json:"theme"`
	Objectives      []string         `json:"objectives"`
	KeyConcepts     []string         `json:"key_concepts"`
	Prerequisites   []string         `json:"prerequisites"`
	DailyMilestones []DailyMilestone `json:"daily_milestones"`
	AdaptiveNotes   string           `json:"adaptive_notes"`
}

// CompleteLearningPlan represents the full structure
type CompleteLearningPlan struct {
	ID              int64               `gorm:"primaryKey" json:"id" example:"1"`
	Goal            string              `json:"goal" example:"Learn React and TypeScript"`
	TotalWeeks      int                 `json:"total_weeks" example:"8"`
	DailyCommitment int                 `json:"daily_commitment_minutes" example:"30"`
	WeeklyThemes    []WeeklyTheme       `json:"weekly_themes"`
	Prerequisites   map[string][]string `json:"prerequisites" example:"{\"week1\":[\"JavaScript Basics\"]}"`
	AdaptiveRules   map[string]string   `json:"adaptive_rules" example:"{\"difficulty\":\"auto\"}"` // rules for content adaptation
}

// GeneratedWeeklyContent represents the content stored in the database
type GeneratedWeeklyContent struct {
	ID               int64          `gorm:"primaryKey" json:"id" example:"1"`
	PlanID           int64          `json:"plan_id" example:"1"` // FK to LearningPlanStructure
	UserID           int64          `json:"user_id" example:"1"`
	WeekNumber       int            `json:"week_number" example:"1"`
	ContentData      datatypes.JSON `json:"content_data" swaggertype:"object"`       // JSONB: stores int64
	GeneratedBasedOn datatypes.JSON `json:"generated_based_on" swaggertype:"object"` // JSONB: snapshot of user progress
	CreatedAt        time.Time      `json:"created_at"`
}

// ContentAdaptationFlag represents flags for content regeneration
type ContentAdaptationFlag struct {
	PlanID            uint   `json:"plan_id" example:"1"`
	WeekNumber        int    `json:"week_number" example:"1"`
	NeedsRegeneration bool   `json:"needs_regeneration" example:"false"`
	Reason            string `json:"reason" example:"User struggling with concepts"`
}

type DailyContent struct {
	BaseModel
	PlanID     int64          `json:"plan_id" example:"1"`
	UserID     int64          `json:"user_id" example:"1"`
	WeekNumber int            `json:"week_number" example:"1"`
	DayNumber  int            `json:"day_number" example:"1"`
	Content    datatypes.JSON `json:"content" swaggertype:"object"`   // The main lesson/content for the day
	Exercises  datatypes.JSON `json:"exercises" swaggertype:"object"` // Exercises for the day
	Resources  datatypes.JSON `json:"resources" swaggertype:"object"` // List of resource links
}
