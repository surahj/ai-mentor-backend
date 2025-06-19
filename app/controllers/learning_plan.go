package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/surahj/ai-mentor-backend/app/models"
	"github.com/surahj/ai-mentor-backend/app/utils"
	"gorm.io/datatypes"
)

type GeneratePlanRequest struct {
	CategoryID      uint `json:"category_id"`
	DailyCommitment int  `json:"daily_commitment"`
}

type StructureRequest struct {
	Goal       string `json:"goal"`
	TotalWeeks int    `json:"total_weeks"`
}

type ContentRequest struct {
	PlanID       uint           `json:"plan_id"`
	WeekNumber   int            `json:"week_number"`
	UserProgress map[string]any `json:"user_progress"`
}

func (c *Controller) GenerateLearningPlan(ctx echo.Context) error {
	// Get user ID from JWT claims
	userToken := ctx.Get("user")
	var userID uint
	if userToken != nil {
		if id, ok := userToken.(float64); ok {
			userID = uint(id)
		}
	}
	if userID == 0 {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Fetch user from DB
	var user models.User
	if err := c.DB.First(&user, userID).Error; err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "User not found"})
	}

	var req GeneratePlanRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Fetch category for prompt enrichment
	var category models.Category
	if err := c.DB.First(&category, req.CategoryID).Error; err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Category not found"})
	}

	// Build prompt
	prompt := "Create a detailed JSON learning plan for the following user profile: "
	prompt += "Goal: " + category.Name + ". "
	prompt += "Daily commitment: " + strconv.Itoa(req.DailyCommitment) + " minutes. "
	if user.Age != nil {
		prompt += "Age: " + strconv.Itoa(*user.Age) + ". "
	}
	if user.Level != nil {
		prompt += "Level: " + *user.Level + ". "
	}
	if user.Background != nil {
		prompt += "Background: " + *user.Background + ". "
	}
	if user.PreferredLanguage != nil {
		prompt += "Preferred language: " + *user.PreferredLanguage + ". "
	}
	if user.Interests != nil {
		prompt += "Interests: " + *user.Interests + ". "
	}
	prompt += "Return a JSON object with a 4-week curriculum, each week split into daily milestones, with topics, explanations, quizzes, and resource links."

	// Call OpenAI
	planJSON, err := utils.GenerateLearningPlan(prompt)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate plan: " + err.Error()})
	}

	// Save to DB
	learningPath := models.LearningPath{
		Name:            category.Name + " - " + strconv.Itoa(req.DailyCommitment) + "min",
		Description:     "AI-generated learning plan for " + category.Name,
		CategoryID:      req.CategoryID,
		DailyCommitment: req.DailyCommitment,
		PlanJSON:        planJSON,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := c.DB.Create(&learningPath).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save learning plan"})
	}

	return ctx.JSON(http.StatusOK, learningPath)
}

func (c *Controller) GeneratePlanStructure(ctx echo.Context) error {
	var req StructureRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID := ctx.Get("user_id").(int64) // adjust as per your auth middleware

	log.Printf("Generating structure for user %d with goal %s", userID, req.Goal)
	prompt := "Create a learning roadmap structure for: " + req.Goal + `\nInclude:\n- Weekly themes and objectives\n- Key milestones\n- Prerequisite mapping\n- Flexible checkpoints for adaptation\nDO NOT generate detailed lesson content yet. Return as JSON.`
	planJSON, err := utils.GenerateLearningPlan(prompt)

	log.Printf("Plan JSON: %s", planJSON)

	if err != nil {
		log.Printf("Error generating structure: %v", err)
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: "Failed to generate structure: " + err.Error(),
		})
	}

	var structure datatypes.JSON
	if err := json.Unmarshal([]byte(planJSON), &structure); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "OpenAI did not return valid JSON"})
	}

	plan := models.LearningPlanStructure{
		UserID:     userID,
		Goal:       req.Goal,
		TotalWeeks: req.TotalWeeks,
		Structure:  structure,
	}
	if err := c.DB.Create(&plan).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save structure"})
	}

	return ctx.JSON(http.StatusOK, plan)
}

// GET /learning-plan/structure/:id
func (c *Controller) GetPlanStructure(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var plan models.LearningPlanStructure
	if err := c.DB.First(&plan, id).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Plan structure not found"})
	}
	return ctx.JSON(http.StatusOK, plan)
}

// POST /learning-plan/content
func (c *Controller) GenerateWeekContent(ctx echo.Context) error {
	var req ContentRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	var plan models.LearningPlanStructure
	if err := c.DB.First(&plan, req.PlanID).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Plan structure not found"})
	}

	progressJSON, _ := json.Marshal(req.UserProgress)
	prompt := "Generate Week " + strconv.Itoa(req.WeekNumber) + " content for learning goal: " + plan.Goal + ".\nUser's current progress: " + string(progressJSON) + ".\nAdapt difficulty and focus areas accordingly. Return as JSON."
	contentJSON, err := utils.GenerateLearningPlan(prompt)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate content: " + err.Error()})
	}

	var contentData datatypes.JSON
	if err := json.Unmarshal([]byte(contentJSON), &contentData); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "OpenAI did not return valid JSON"})
	}

	content := models.GeneratedContent{
		PlanID:           plan.ID,
		WeekNumber:       req.WeekNumber,
		ContentData:      contentData,
		GeneratedBasedOn: datatypes.JSON(progressJSON),
		CreatedAt:        time.Now(),
	}
	if err := c.DB.Create(&content).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save content"})
	}
	return ctx.JSON(http.StatusOK, content)
}

// GET /learning-plan/content/:plan_id/:week_number
func (c *Controller) GetWeekContent(ctx echo.Context) error {
	planID, _ := strconv.Atoi(ctx.Param("plan_id"))
	week, _ := strconv.Atoi(ctx.Param("week_number"))
	var content models.GeneratedContent
	if err := c.DB.Where("plan_id = ? AND week_number = ?", planID, week).First(&content).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Content not found"})
	}
	return ctx.JSON(http.StatusOK, content)
}
