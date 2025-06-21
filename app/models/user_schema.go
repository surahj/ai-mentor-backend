package models

// CreateUserRequest for user registration

// UserResponse for API responses
type UserResponse struct {
	ID              int64   `json:"id"`
	Email           string `json:"email"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	DailyCommitment int    `json:"daily_commitment"`
	LearningGoal    string `json:"learning_goal"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}
