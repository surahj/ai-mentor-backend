package controllers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/surahj/ai-mentor-backend/models"
	"gorm.io/gorm"
)

type StatsController struct {
	DB *gorm.DB
}

func NewStatsController(db *gorm.DB) *StatsController {
	return &StatsController{DB: db}
}

// Statistics response
type DashboardStats struct {
	TotalSessions int64                `json:"total_sessions"`
	TotalMinutes  int                  `json:"total_minutes"`
	Streak        int                  `json:"streak"`
	TopTopics     []TopicStats         `json:"top_topics"`
	DailyActivity []DailyActivityStats `json:"daily_activity"`
}

type TopicStats struct {
	Topic    string `json:"topic"`
	Minutes  int    `json:"minutes"`
	Sessions int    `json:"sessions"`
}

type DailyActivityStats struct {
	Date    string `json:"date"`
	Minutes int    `json:"minutes"`
}

// Get dashboard statistics
func (sc *StatsController) GetDashboardStats(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	stats := DashboardStats{}

	// Total sessions
	sc.DB.Model(&models.Session{}).Where("user_id = ?", userID).Count(&stats.TotalSessions)

	// Total minutes
	var totalMinutes int64
	sc.DB.Model(&models.Session{}).Where("user_id = ?", userID).
		Select("COALESCE(SUM(duration_minutes), 0)").Row().Scan(&totalMinutes)
	stats.TotalMinutes = int(totalMinutes)

	// Calculate streak
	stats.Streak = calculateStreak(sc.DB, userID.(uint))

	// Top topics
	stats.TopTopics = getTopTopics(sc.DB, userID.(uint))

	// Daily activity (last 7 days)
	stats.DailyActivity = getDailyActivity(sc.DB, userID.(uint))

	return c.JSON(http.StatusOK, stats)
}

// Helper function to calculate streak
func calculateStreak(db *gorm.DB, userID uint) int {
	// Get distinct dates when user had sessions, ordered by date DESC
	var dates []time.Time
	db.Model(&models.Session{}).
		Where("user_id = ?", userID).
		Distinct("DATE(date)").
		Order("DATE(date) DESC").
		Pluck("DATE(date)", &dates)

	if len(dates) == 0 {
		return 0
	}

	// Calculate streak
	streak := 1
	today := time.Now().Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)

	// Check if most recent day is today or yesterday
	mostRecent := dates[0].Truncate(24 * time.Hour)
	if mostRecent.After(today) || (mostRecent.Before(yesterday) && !mostRecent.Equal(yesterday)) {
		return 0 // Most recent day is not today or yesterday
	}

	// Count consecutive days
	for i := 0; i < len(dates)-1; i++ {
		curr := dates[i].Truncate(24 * time.Hour)
		next := dates[i+1].Truncate(24 * time.Hour)

		diff := curr.Sub(next).Hours() / 24

		if diff == 1 {
			// Consecutive day
			streak++
		} else {
			// Streak broken
			break
		}
	}

	return streak
}

// Helper function to get top topics
func getTopTopics(db *gorm.DB, userID uint) []TopicStats {
	var topTopics []TopicStats

	rows, err := db.Raw(`
		SELECT topic, 
			   SUM(duration_minutes) as minutes,
			   COUNT(*) as sessions
		FROM sessions 
		WHERE user_id = ?
		GROUP BY topic
		ORDER BY minutes DESC
		LIMIT 5
	`, userID).Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var t TopicStats
			rows.Scan(&t.Topic, &t.Minutes, &t.Sessions)
			topTopics = append(topTopics, t)
		}
	}

	return topTopics
}

// Helper function to get daily activity
func getDailyActivity(db *gorm.DB, userID uint) []DailyActivityStats {
	var dailyActivity []DailyActivityStats

	// Get dates for last 7 days
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		var minutes int
		db.Model(&models.Session{}).
			Where("user_id = ? AND DATE(date) = ?", userID, dateStr).
			Select("COALESCE(SUM(duration_minutes), 0)").
			Row().Scan(&minutes)

		dailyActivity = append(dailyActivity, DailyActivityStats{
			Date:    dateStr,
			Minutes: minutes,
		})
	}

	return dailyActivity
}
