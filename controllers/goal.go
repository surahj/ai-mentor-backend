package controllers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/surahj/ai-mentor-backend/models"
	"gorm.io/gorm"
)

type GoalController struct {
	DB *gorm.DB
}

func NewGoalController(db *gorm.DB) *GoalController {
	return &GoalController{DB: db}
}

// Create a new goal
func (gc *GoalController) Create(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	// Parse input
	var input struct {
		Title         string    `json:"title"`
		Description   string    `json:"description"`
		TargetMinutes int       `json:"target_minutes"`
		IsWeekly      bool      `json:"is_weekly"`
		StartDate     time.Time `json:"start_date"`
		EndDate       time.Time `json:"end_date"`
		IsRecurring   bool      `json:"is_recurring"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	// Validate required fields
	if input.Title == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Title is required"})
	}
	if input.TargetMinutes <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Target minutes must be positive"})
	}

	// Use current time if start date not provided
	startDate := input.StartDate
	if startDate.IsZero() {
		startDate = time.Now()
	}

	// Create goal
	goal := models.Goal{
		UserID:        userID.(uint),
		Title:         input.Title,
		Description:   input.Description,
		TargetMinutes: input.TargetMinutes,
		IsWeekly:      input.IsWeekly,
		StartDate:     startDate,
		EndDate:       input.EndDate,
		IsRecurring:   input.IsRecurring,
	}

	if err := gc.DB.Create(&goal).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create goal"})
	}

	return c.JSON(http.StatusCreated, goal)
}

// Get all goals for the current user
func (gc *GoalController) List(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	var goals []models.Goal
	if err := gc.DB.Where("user_id = ?", userID).Find(&goals).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch goals"})
	}

	return c.JSON(http.StatusOK, goals)
}

// Get a specific goal
func (gc *GoalController) Get(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	id := c.Param("id")

	var goal models.Goal
	if err := gc.DB.Where("id = ? AND user_id = ?", id, userID).First(&goal).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Goal not found"})
	}

	return c.JSON(http.StatusOK, goal)
}

// Update a goal
func (gc *GoalController) Update(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	id := c.Param("id")

	var goal models.Goal
	if err := gc.DB.Where("id = ? AND user_id = ?", id, userID).First(&goal).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Goal not found"})
	}

	var input struct {
		Title         string    `json:"title"`
		Description   string    `json:"description"`
		TargetMinutes int       `json:"target_minutes"`
		IsWeekly      bool      `json:"is_weekly"`
		StartDate     time.Time `json:"start_date"`
		EndDate       time.Time `json:"end_date"`
		IsRecurring   bool      `json:"is_recurring"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	// Update fields if provided
	if input.Title != "" {
		goal.Title = input.Title
	}
	if input.Description != "" {
		goal.Description = input.Description
	}
	if input.TargetMinutes > 0 {
		goal.TargetMinutes = input.TargetMinutes
	}

	// These fields can be explicitly set to different values
	goal.IsWeekly = input.IsWeekly
	goal.IsRecurring = input.IsRecurring

	if !input.StartDate.IsZero() {
		goal.StartDate = input.StartDate
	}
	if !input.EndDate.IsZero() {
		goal.EndDate = input.EndDate
	}

	if err := gc.DB.Save(&goal).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update goal"})
	}

	return c.JSON(http.StatusOK, goal)
}

// Delete a goal
func (gc *GoalController) Delete(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	id := c.Param("id")

	var goal models.Goal
	if err := gc.DB.Where("id = ? AND user_id = ?", id, userID).First(&goal).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Goal not found"})
	}

	if err := gc.DB.Delete(&goal).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete goal"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Goal deleted successfully"})
}

// Get progress toward goals
func (gc *GoalController) GetProgress(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	// Get all active goals
	var goals []models.Goal
	if err := gc.DB.Where("user_id = ?", userID).Find(&goals).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch goals"})
	}

	progress := []map[string]interface{}{}

	// For each goal, calculate progress
	for _, goal := range goals {
		// Calculate date range based on goal type
		var startDate, endDate time.Time
		now := time.Now()

		if goal.IsWeekly {
			// Weekly goal - get current week
			startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
			startDate = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, now.Location())
			endDate = startDate.AddDate(0, 0, 7)
		} else {
			// Daily goal - today only
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			endDate = startDate.AddDate(0, 0, 1)
		}

		// Get total minutes for this period
		var totalMinutes int64
		gc.DB.Model(&models.Session{}).
			Where("user_id = ? AND date >= ? AND date < ?", userID, startDate, endDate).
			Select("COALESCE(SUM(duration_minutes), 0)").
			Row().Scan(&totalMinutes)

		// Calculate percentage
		var percentage float64
		if goal.TargetMinutes > 0 {
			percentage = float64(totalMinutes) / float64(goal.TargetMinutes) * 100
			if percentage > 100 {
				percentage = 100
			}
		}

		// Add to progress list
		progress = append(progress, map[string]interface{}{
			"goal_id":      goal.ID,
			"title":        goal.Title,
			"target":       goal.TargetMinutes,
			"actual":       totalMinutes,
			"percentage":   percentage,
			"is_weekly":    goal.IsWeekly,
			"period_start": startDate,
			"period_end":   endDate,
		})
	}

	return c.JSON(http.StatusOK, progress)
}
