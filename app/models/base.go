package models

import (
	"time"
)

// BaseModel is the base model for all models
type BaseModel struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
