package entity

type CustomTestResult struct {
	CustomTestResultId string `json:"testResultId"`
	SubmissionId       string `json:"submissionId"`
	Input              string `json:"input"`
	ExpectedOutput     string `json:"expectedOutput"`
	ActualOutput       string `json:"actualOutput"`
	Result             string `json:"result"`
}
