package models

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	DailyCommitment int    `json:"daily_commitment"`
	LearningGoal    string `json:"learning_goal"`
	Age             int    `json:"age"`
	Level           string `json:"level"`
}

type CreateUserRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	DailyCommitment int    `json:"daily_commitment" validate:"required"`
	LearningGoal    string `json:"learning_goal" validate:"required"`
}

type UpdateProfileRequest struct {
	Age               *int    `json:"age"`
	Level             *string `json:"level"`
	Background        *string `json:"background"`
	PreferredLanguage *string `json:"preferred_language"`
	Interests         *string `json:"interests"`
	Country           *string `json:"country"`
}
