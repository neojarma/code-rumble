package entity

type TestCase struct {
	TestCaseId string `json:"id,omitempty"`
	QuestionId string `json:"-"`
	Input      string `json:"input"`
	Output     string `json:"output"`
}
