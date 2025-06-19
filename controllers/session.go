package controllers

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/surahj/ai-mentor-backend/models"
	"gorm.io/gorm"
)

type SessionController struct {
	DB *gorm.DB
}

func NewSessionController(db *gorm.DB) *SessionController {
	return &SessionController{DB: db}
}

// Create a new learning session
func (sc *SessionController) Create(c echo.Context) error {
	// Get user ID from context
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	// Parse input
	var input struct {
		Topic           string    `json:"topic"`
		DurationMinutes int       `json:"duration_minutes"`
		Notes           string    `json:"notes"`
		Date            time.Time `json:"date"`
		Rating          int       `json:"rating"`     // 1-5 rating
		Tags            string    `json:"tags"`       // Comma-separated tags
		Reflection      string    `json:"reflection"` // Post-session reflection
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	// Validate required fields
	if input.Topic == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Topic is required"})
	}
	if input.DurationMinutes <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Duration must be positive"})
	}

	// Use current time if date not provided
	sessionDate := input.Date
	if sessionDate.IsZero() {
		sessionDate = time.Now()
	}

	// Create session
	session := models.Session{
		UserID:          userID.(uint),
		Topic:           input.Topic,
		DurationMinutes: input.DurationMinutes,
		Notes:           input.Notes,
		Date:            sessionDate,
		Rating:          input.Rating,     // Add this line
		Tags:            input.Tags,       // Add this line
		Reflection:      input.Reflection, // Add this line
	}

	if err := sc.DB.Create(&session).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create session"})
	}

	return c.JSON(http.StatusCreated, session)
}

// Get a specific session
func (sc *SessionController) Get(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	id := c.Param("id")

	var session models.Session
	if err := sc.DB.Where("id = ? AND user_id = ?", id, userID).First(&session).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Session not found"})
	}

	return c.JSON(http.StatusOK, session)
}

// List all sessions for a user with filtering
func (sc *SessionController) List(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	// Get query parameters for filtering
	topic := c.QueryParam("topic")
	fromDate := c.QueryParam("from_date")
	toDate := c.QueryParam("to_date")
	tag := c.QueryParam("tag")
	ratingStr := c.QueryParam("rating")

	// Build query
	query := sc.DB.Where("user_id = ?", userID)

	// Apply filters if provided
	if topic != "" {
		query = query.Where("topic LIKE ?", "%"+topic+"%") // Changed to LIKE for partial matches
	}

	if fromDate != "" {
		query = query.Where("DATE(date) >= DATE(?)", fromDate) // Added DATE() function
	}

	if toDate != "" {
		query = query.Where("DATE(date) <= DATE(?)", toDate) // Added DATE() function
	}

	// Filter by tag
	if tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}

	// Filter by rating
	if ratingStr != "" {
		var rating int
		if _, err := fmt.Sscanf(ratingStr, "%d", &rating); err == nil && rating > 0 {
			query = query.Where("rating = ?", rating)
		}
	}

	// Debug: Print the SQL query
	fmt.Println("DEBUG SQL:", query.Statement.SQL.String())

	// Execute query
	var sessions []models.Session
	if err := query.Order("date DESC").Find(&sessions).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch sessions"})
	}

	return c.JSON(http.StatusOK, sessions)
}

// Update an existing session
func (sc *SessionController) Update(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	// Get session ID from URL parameter
	id := c.Param("id")

	// Check if session exists and belongs to user
	var session models.Session
	if err := sc.DB.Where("id = ? AND user_id = ?", id, userID).First(&session).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Session not found"})
	}

	// Parse input
	var input struct {
		Topic           string    `json:"topic"`
		DurationMinutes int       `json:"duration_minutes"`
		Notes           string    `json:"notes"`
		Date            time.Time `json:"date"`
		Rating          int       `json:"rating"`     // 1-5 rating
		Tags            string    `json:"tags"`       // Comma-separated tags
		Reflection      string    `json:"reflection"` // Post-session reflection
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	// Update fields if provided
	if input.Topic != "" {
		session.Topic = input.Topic
	}
	if input.DurationMinutes > 0 {
		session.DurationMinutes = input.DurationMinutes
	}
	if input.Notes != "" {
		session.Notes = input.Notes
	}
	if !input.Date.IsZero() {
		session.Date = input.Date
	}
	if input.Rating > 0 {
		session.Rating = input.Rating
	}
	if input.Tags != "" {
		session.Tags = input.Tags
	}
	if input.Reflection != "" {
		session.Reflection = input.Reflection
	}

	// Save updates
	if err := sc.DB.Save(&session).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update session"})
	}

	return c.JSON(http.StatusOK, session)
}

// Delete a session
func (sc *SessionController) Delete(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	// Get session ID from URL parameter
	id := c.Param("id")

	// Check if session exists and belongs to user
	var session models.Session
	if err := sc.DB.Where("id = ? AND user_id = ?", id, userID).First(&session).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Session not found"})
	}

	// Delete the session
	if err := sc.DB.Delete(&session).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete session"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Session deleted successfully"})
}

// GetTags returns all unique tags used by the user
func (sc *SessionController) GetTags(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	// Get distinct tags used by this user
	var sessions []models.Session
	if err := sc.DB.Where("user_id = ? AND tags <> ''", userID).Find(&sessions).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch sessions"})
	}

	// Extract and deduplicate tags
	tagMap := make(map[string]bool)
	for _, session := range sessions {
		tags := strings.Split(session.Tags, ",")
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tagMap[tag] = true
			}
		}
	}

	// Convert to array
	var tagList []string
	for tag := range tagMap {
		tagList = append(tagList, tag)
	}

	// Sort alphabetically
	sort.Strings(tagList)

	return c.JSON(http.StatusOK, tagList)
}
