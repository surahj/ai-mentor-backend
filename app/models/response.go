package models

type ErrorResponse struct {
	ErrorCode    int    `json:"error_code"  validate:"required"`
	ErrorMessage string `json:"error_message"  validate:"required"`
}

type SuccessResponse struct {
	Status  int         `json:"status" validate:"required"`
	Message string      `json:"message" validate:"required"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseMessage struct {
	Status  int         `json:"status"  validate:"required"`
	Message interface{} `json:"message"  validate:"required"`
}
