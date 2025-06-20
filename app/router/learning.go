package router

import "github.com/labstack/echo/v4"

// @Summary Generate Plan Structure
// @Description Generate and store a high-level learning plan structure for a user
// @Tags LearningPlan
// @Param request body controllers.StructureRequest true "Learning Plan Structure Request"
// @Accept json
// @Produce json
// @Success 200 {object} models.LearningPlanStructure
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /learnings/structure [post]
func (a *App) GeneratePlanStructure(c echo.Context) error {
	return a.Controller.GeneratePlanStructure(c)
}

// @Summary Get Plan Structure
// @Description Retrieve a learning plan structure by ID
// @Tags LearningPlan
// @Param id path int true "Plan Structure ID"
// @Produce json
// @Success 200 {object} models.LearningPlanStructure
// @Failure 404 {object} models.ErrorResponse
// @Router /learnings/structure/{id} [get]
func (a *App) GetPlanStructure(c echo.Context) error {
	return a.Controller.GetPlanStructure(c)
}

// @Summary Generate Week Content
// @Description Generate and store detailed weekly content for a learning plan
// @Tags LearningPlan
// @Param request body controllers.ContentRequest true "Weekly Content Request"
// @Accept json
// @Produce json
// @Success 200 {object} models.GeneratedContent
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /learnings/weekly-content [post]
func (a *App) GenerateWeekContent(c echo.Context) error {
	return a.Controller.GenerateWeekContent(c)
}

// @Summary Get Week Content
// @Description Retrieve generated content for a specific week of a learning plan
// @Tags LearningPlan
// @Param plan_id path int true "Plan ID"
// @Param week_number path int true "Week Number"
// @Produce json
// @Success 200 {object} models.GeneratedContent
// @Failure 404 {object} models.ErrorResponse
// @Router /learnings/weekly-content/{week_number}/{plan_id} [get]
func (a *App) GetWeekContent(c echo.Context) error {
	return a.Controller.GetWeekContent(c)
}

// @Summary Get My Learning Plans
// @Description Retrieve all learning plans for the authenticated user
// @Tags LearningPlan
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /learnings [get]
func (a *App) GetLearnings(c echo.Context) error {
	return a.Controller.GetMyLearnings(c)
}

// @Summary Get Daily Content
// @Description Retrieve daily content for a specific day of a learning plan
// @Tags LearningPlan
// @Param plan_id path int true "Plan ID"
// @Param week_number path int true "Week Number"
// @Param day_number path int true "Day Number"
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse

// @Router /learnings/daily-content/{day_number}/{week_number}/{plan_id} [get]
func (a *App) GetDailyContent(c echo.Context) error {
	return a.Controller.GetDailyContent(c)
}

// @Summary Generate Exercises for Daily Content
// @Description Generate exercises for a specific day of a learning plan
// @Tags LearningPlan
// @Param plan_id path int true "Plan ID"
// @Param week_number path int true "Week Number"
// @Param day_number path int true "Day Number"
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse

// @Router /learnings/daily-content/{day_number}/{week_number}/{plan_id}/exercises [get]
func (a *App) GenerateDailyExercises(c echo.Context) error {
	return a.Controller.GenerateDailyExercises(c)
}
