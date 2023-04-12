package entity

type WSResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
