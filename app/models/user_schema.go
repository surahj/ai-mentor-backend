package models

// CreateUserRequest for user registration
type CreateUserRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	DailyCommitment int    `json:"daily_commitment" validate:"required"`
	LearningGoal    string `json:"learning_goal" validate:"required"`
}

// UserResponse for API responses
type UserResponse struct {
	ID              uint   `json:"id"`
	Email           string `json:"email"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	DailyCommitment int    `json:"daily_commitment"`
	LearningGoal    string `json:"learning_goal"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}
