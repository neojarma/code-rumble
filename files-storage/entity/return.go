package entity

type Return struct {
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
}
