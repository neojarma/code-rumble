package entity

type Question struct {
	QuestionId string `json:"questionId"`
	RunCode    string `json:"runCode"`
	StubCode   string `json:"stubCode"`
}
