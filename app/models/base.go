package models

import (
	"time"

)

// BaseModel is the base model for all models
type BaseModel struct {
	ID        int64 		`gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}