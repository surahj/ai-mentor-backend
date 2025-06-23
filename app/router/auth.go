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

// @Summary Verify OTP
// @Description This API will attempt to verify a user's OTP
// @Tags Authentication
// @Param request body controllers.VerifyOTPRequest true "Email and OTP"
// @Accept json
// @Produce json
// @Success      200  {object}  models.UserResponse "Status 200 will be returned if the verification was successful"
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /verify-otp [post]
func (a *App) VerifyOTP(c echo.Context) error {
	return a.Controller.VerifyOTP(c)
}

// @Summary Resend OTP
// @Description This API will attempt to resend OTP to user's email
// @Tags Authentication
// @Param request body controllers.ResendOTPRequest true "Email"
// @Accept json
// @Produce json
// @Success      200  {object}  models.SuccessResponse "Status 200 will be returned if the OTP was resent successfully"
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /resend-otp [post]
func (a *App) ResendOTP(c echo.Context) error {
	return a.Controller.ResendOTP(c)
}

// @Summary Forgot Password
// @Description This API will attempt to send password reset OTP to user's email
// @Tags Authentication
// @Param request body controllers.ForgotPasswordRequest true "Email"
// @Accept json
// @Produce json
// @Success      200  {object}  models.SuccessResponse "Status 200 will be returned if the password reset OTP was sent successfully"
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /forgot-password [post]
func (a *App) ForgotPassword(c echo.Context) error {
	return a.Controller.ForgotPassword(c)
}

// @Summary Reset Password
// @Description This API will attempt to reset user's password with OTP verification
// @Tags Authentication
// @Param request body controllers.ResetPasswordRequest true "Email, OTP, and new password"
// @Accept json
// @Produce json
// @Success      200  {object}  models.SuccessResponse "Status 200 will be returned if the password was reset successfully"
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /reset-password [post]
func (a *App) ResetPassword(c echo.Context) error {
	return a.Controller.ResetPassword(c)
}

// @Summary Google Login
// @Description This API will authenticate a user with a Google ID token
// @Tags Authentication
// @Param request body controllers.GoogleLoginRequest true "Google ID Token"
// @Accept json
// @Produce json
// @Success      200  {object}  models.SuccessResponse "Status 200 will be returned if the login was successful"
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /auth/google/login [post]
func (a *App) GoogleLogin(c echo.Context) error {
	return a.Controller.GoogleLogin(c)
}
