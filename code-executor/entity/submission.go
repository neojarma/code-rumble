package entity

import (
	"encoding/json"
)

type Submission struct {
	SubmissionId   string      `json:"submissionId"`
	QuestionId     string      `json:"questionId"`
	SubmittedCode  string      `json:"submittedCode"`
	RunCode        string      `json:"runCode"`
	CustomTestCase bool        `json:"customTestCase"`
	TestCases      []*TestCase `json:"testCases"`
}

func (s *Submission) MarshallTestCase() ([]byte, error) {
	b, err := json.Marshal(s.TestCases)
	if err != nil {
		return nil, err
	}

	return b, nil
}
