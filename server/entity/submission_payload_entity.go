package entity

type SubmissionPayload struct {
	SubmissionId   string      `json:"submissionId"`
	QuestionId     string      `json:"questionId"`
	SubmittedCode  string      `json:"submittedCode"`
	CustomTestCase bool        `json:"customTestCase"`
	RunCode        string      `json:"runCode"`
	TestCases      []*TestCase `json:"testCases,omitempty"`
}
