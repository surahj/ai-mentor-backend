package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/surahj/ai-mentor-backend/app/library"
	"github.com/surahj/ai-mentor-backend/app/models"
	"github.com/surahj/ai-mentor-backend/app/utils"
	"gorm.io/datatypes"
)

type GeneratePlanRequest struct {
	CategoryID      uint `json:"category_id"`
	DailyCommitment int  `json:"daily_commitment"`
}

type StructureRequest struct {
	Goal            string `json:"goal"`
	TotalWeeks      int    `json:"total_weeks"`
	DailyCommitment int    `json:"daily_commitment"`
}

type ContentRequest struct {
	PlanID       int64                  `json:"plan_id"`
	WeekNumber   int                    `json:"week_number"`
	UserProgress map[string]interface{} `json:"user_progress"`
}

func (c *Controller) GeneratePlanStructure(ctx echo.Context) error {
	userID, err := library.GetUserIDFronContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	if userID == 0 {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	log.Printf("User ID in GeneratePlanStructure: %v", userID)

	var req StructureRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Goal == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Goal is required"})
	}

	if req.TotalWeeks == 0 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Total weeks is required"})
	}

	if req.DailyCommitment == 0 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Daily commitment is required"})
	}

	// check if the goal is already in the database
	var existingPlan models.LearningPlanStructure
	if err := c.DB.Where("goal = ?", req.Goal).First(&existingPlan).Error; err == nil {
		return ctx.JSON(http.StatusOK, models.SuccessResponse{
			Status:  http.StatusOK,
			Message: "Goal retrieved successfully",
			Data:    existingPlan,
		})
	}

	// Generate the learning plan structure using OpenAI
	plan, err := utils.GenerateLearningPlanStructure(req.Goal, req.TotalWeeks, req.DailyCommitment)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate structure: " + err.Error()})
	}

	// Convert the plan to JSON for storage
	planJSON, err := json.Marshal(plan)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed tou serialize plan"})
	}

	log.Printf("Plan: %v", planJSON)

	// Save to database
	learningPlan := models.LearningPlanStructure{
		UserID:     userID,
		Goal:       req.Goal,
		TotalWeeks: req.TotalWeeks,
		Structure:  datatypes.JSON(planJSON),
	}

	if err := c.DB.Create(&learningPlan).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save structure"})
	}

	return ctx.JSON(http.StatusOK, models.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Learning plan structure generated successfully",
		Data: map[string]interface{}{
			"id":   learningPlan.ID,
			"plan": plan,
		},
	})
}

func (c *Controller) GetPlanStructure(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var plan models.LearningPlanStructure
	if err := c.DB.First(&plan, id).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Plan structure not found"})
	}

	// Parse the JSON structure back to the complete plan
	var completePlan models.CompleteLearningPlan
	if err := json.Unmarshal(plan.Structure, &completePlan); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse plan structure"})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"id":   plan.ID,
		"plan": completePlan,
	})
}

func (c *Controller) GenerateWeekContent(ctx echo.Context) error {
	var req ContentRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.PlanID == 0 {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Plan ID is required",
		})
	}

	if req.WeekNumber == 0 {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Week number is required",
		})
	}

	// Get the plan structure to extract the goal
	userID, err := library.GetUserIDFronContext(ctx)
	if err != nil || userID == 0 {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// check if the plan exists

	var plan models.LearningPlanStructure
	if err := c.DB.Where("id = ? AND user_id = ?", req.PlanID, userID).First(&plan).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Plan structure not found"})
	}

	// Generate weekly content using OpenAI
	content, err := utils.GenerateWeeklyContent(plan.Goal, req.WeekNumber, req.UserProgress)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate content: " + err.Error()})
	}

	// Convert content to JSON for storage
	contentJSON, err := json.Marshal(content)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to serialize content"})
	}

	// Convert user progress to JSON for storage
	progressJSON, _ := json.Marshal(req.UserProgress)

	// Save to database
	generatedContent := models.GeneratedWeeklyContent{
		PlanID:           req.PlanID,
		WeekNumber:       req.WeekNumber,
		ContentData:      datatypes.JSON(contentJSON),
		GeneratedBasedOn: datatypes.JSON(progressJSON),
		CreatedAt:        time.Now(),
	}

	if err := c.DB.Create(&generatedContent).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save content"})
	}

	return ctx.JSON(http.StatusOK, models.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Content generated successfully",
		Data: map[string]interface{}{
			"id":      generatedContent.ID,
			"content": content,
		},
	})
}

func (c *Controller) GetWeekContent(ctx echo.Context) error {
	planID, _ := strconv.ParseInt(ctx.Param("plan_id"), 10, 64)
	week, _ := strconv.Atoi(ctx.Param("week_number"))

	if planID == 0 {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Plan ID is required",
		})
	}

	if week == 0 {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Week number is required",
		})
	}

	userID, err := library.GetUserIDFronContext(ctx)
	if err != nil || userID == 0 {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var content models.GeneratedWeeklyContent
	if err := c.DB.Where("plan_id = ? AND week_number = ? AND user_id = ?", planID, week, userID).First(&content).Error; err != nil {
		return ctx.JSON(http.StatusOK, models.ErrorResponse{
			ErrorCode:    http.StatusNotFound,
			ErrorMessage: "Content not found",
		})
	}

	// Parse the JSON content back to the weekly content structure
	var weeklyContent models.WeeklyContent
	if err := json.Unmarshal(content.ContentData, &weeklyContent); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse content"})
	}

	return ctx.JSON(http.StatusOK, models.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Content fetched successfully",
		Data:    weeklyContent,
	})
}

func (c *Controller) GetMyLearnings(ctx echo.Context) error {
	userID, err := library.GetUserIDFronContext(ctx)
	if err != nil || userID == 0 {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var plans []models.LearningPlanStructure
	if err := c.DB.Where("user_id = ?", userID).Find(&plans).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch user plans"})
	}

	return ctx.JSON(http.StatusOK, models.SuccessResponse{
		Status:  http.StatusOK,
		Message: "User learning plans fetched successfully",
		Data:    plans,
	})
}

// GET /learning-plan/daily-content/:plan_id/:week_number/:day_number
func (c *Controller) GetDailyContent(ctx echo.Context) error {
	userID, err := library.GetUserIDFronContext(ctx)
	if err != nil || userID == 0 {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	planID, _ := strconv.ParseInt(ctx.Param("plan_id"), 10, 64)
	week, _ := strconv.Atoi(ctx.Param("week_number"))
	day, _ := strconv.Atoi(ctx.Param("day_number"))

	// get the week content
	var weekContent models.GeneratedWeeklyContent
	err = c.DB.Where("plan_id = ? AND week_number = ? AND user_id = ?", planID, week, userID).First(&weekContent).Error
	if err != nil {
		return ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			ErrorCode:    http.StatusNotFound,
			ErrorMessage: "Week content not found",
		})
	}

	var daily models.DailyContent
	err = c.DB.Where("plan_id = ? AND user_id = ? AND week_number = ? AND day_number = ?", planID, userID, week, day).First(&daily).Error
	if err == nil {
		return ctx.JSON(http.StatusOK, models.SuccessResponse{
			Message: "Daily content fetched successfully",
			Data:    daily,
		})
	}

	// Optionally, fetch user progress for this day/plan
	userProgress := map[string]interface{}{} // TODO: fetch from progress table if available

	dailyStructure := string(weekContent.ContentData)
	// Generate content
	plan := models.LearningPlanStructure{}
	c.DB.Where("id = ? AND user_id = ?", planID, userID).First(&plan)
	lesson, resources, genErr := utils.GenerateDailyContent(plan.Goal, dailyStructure, week, day, userProgress)
	if genErr != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate daily content"})
	}
	daily = models.DailyContent{
		PlanID:     planID,
		UserID:     userID,
		WeekNumber: week,
		DayNumber:  day,
		Content:    lesson,
		Resources:  resources,
	}
	c.DB.Create(&daily)

	return ctx.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Daily content generated successfully",
		Data:    daily,
	})
}

func (c *Controller) GenerateDailyExercises(ctx echo.Context) error {
	userID, err := library.GetUserIDFronContext(ctx)
	if err != nil || userID == 0 {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	planID, _ := strconv.ParseInt(ctx.Param("plan_id"), 10, 64)
	week, _ := strconv.Atoi(ctx.Param("week_number"))
	day, _ := strconv.Atoi(ctx.Param("day_number"))

	if planID == 0 {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Plan ID is required",
		})
	}

	if week == 0 {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Week number is required",
		})
	}

	if day == 0 {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Day number is required",
		})
	}

	var daily models.DailyContent
	err = c.DB.Where("plan_id = ? AND user_id = ? AND week_number = ? AND day_number = ?", planID, userID, week, day).First(&daily).Error
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Daily content not found. Please generate the daily lesson first."})
	}

	userProgress := map[string]interface{}{} // TODO: fetch from progress table if available

	exercises, err := utils.GenerateExercisesForLesson(string(daily.Content), userProgress)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate exercises"})
	}

	daily.Exercises = exercises
	if err := c.DB.Save(&daily).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save exercises"})
	}

	return ctx.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Exercises generated and saved successfully",
		Data:    daily.Exercises,
	})
}
