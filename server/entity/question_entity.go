package entity

type Question struct {
	ID           string `json:"id"`
	Description  string `json:"description"`
	TemplateCode string `json:"templateCode"`
}
