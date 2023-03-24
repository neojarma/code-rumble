package entity

type TestResult struct {
	SubmissionId string `json:"submissionId"`
	TestResultId string `json:"testResultId"`
	TestCaseId   string `json:"testCaseId"`
	ActualOutput string `json:"actualOutput"`
	Result       string `json:"result"`
}
