package model

// Generic response with message and optional data
type APIResponse struct {
	Message string      `json:"message" example:"User logged in successfully !"`
	Data    interface{} `json:"data,omitempty"`
}
