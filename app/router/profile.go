package router

import "github.com/labstack/echo/v4"

// @Summary Get Profile
// @Description This API will retrieve user profile information
// @Tags Profile
// @Accept json
// @Produce json
// @Success      200  {object}  models.SuccessResponse "Status 200 will be returned if the profile was retrieved successfully"
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router /profile [get]
func (a *App) GetProfile(c echo.Context) error {
	return a.Controller.GetProfile(c)
}

// @Summary Update Profile
// @Description This API will update user profile information
// @Tags Profile
// @Param request body models.UpdateProfileRequest true "Profile Update Data"
// @Accept json
// @Produce json
// @Success      200  {object}  models.SuccessResponse "Status 200 will be returned if the profile was updated successfully"
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /profile [put]
func (a *App) UpdateProfile(c echo.Context) error {
	return a.Controller.UpdateProfile(c)
}
