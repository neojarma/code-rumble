package entity

type Submission struct {
	SubmissionId   string      `json:"submissionId"`
	QuestionId     string      `json:"questionId"`
	SourceCode     string      `json:"sourceCode"`
	CustomTestCase bool        `json:"customTestCase"`
	TestCase       []*TestCase `json:"testCase,omitempty"`
}
