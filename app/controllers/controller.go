package controllers

import (
	"github.com/surahj/ai-mentor-backend/app/configs"
	"github.com/surahj/ai-mentor-backend/app/services"
	"gorm.io/gorm"
)

type Controller struct {
	DB          *gorm.DB
	EmailClient services.EmailServiceProvider
	Config      *configs.Config
}
