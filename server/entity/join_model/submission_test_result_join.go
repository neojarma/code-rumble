package join_model

type SubmissionTestResult struct {
	SubmissionId     string              `json:"submissionId"`
	SubmissionStatus string              `json:"submissionStatus"`
	QuestionId       string              `json:"questionId"`
	TestResult       []*SubmissionResult `json:"testResult"`
}

type SubmissionResult struct {
	SubmissionId string `json:"submissionId,omitempty"`
	QuestionId   string `json:"questionId,omitempty"`
	Status       string `json:"submissionStatus,omitempty"`
	Input        string `json:"input"`
	Output       string `json:"expectedOutput"`
	ActualOutput string `json:"actualOutput"`
	Result       string `json:"result"`
}
