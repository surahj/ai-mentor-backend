package router

import "github.com/labstack/echo/v4"

// @Summary Sign Up
// @Description This API will attempt to create a new user
// @Tags Authentication
// @Param request body models.CreateUserRequest true "User Details"
// @Accept json
// @Produce json
// @Success      201  {object}  models.SuccessResponse "User created, OTP was send, show the verification page"
// @Success      202  {object}  models.UserResponse "Status 202 will be returned if the signup was successfully, DONT show verification page"
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /signup [post]
func (a *App) SignUp(c echo.Context) error {
	return a.Controller.Register(c)
}

// @Summary Login
// @Description This API will attempt to login a user
// @Tags Authentication
// @Param request body models.LoginRequest true "User Details"
// @Accept json
// @Produce json
// @Success      200  {object}  models.UserResponse "Status 200 will be returned if the login was successfully"
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /login [post]
func (a *App) Login(c echo.Context) error {
	return a.Controller.Login(c)
}
