package entity

import (
	"encoding/json"
	"strings"
)

type Submission struct {
	SubmissionId  string      `json:"submissionId"`
	QuestionId    string      `json:"questionId"`
	SubmittedCode string      `json:"submittedCode"`
	RunCode       string      `json:"runCode"`
	TestCases     []*TestCase `json:"testCases"`
}

func (s *Submission) ToString() string {
	b, err := json.Marshal(s.TestCases)
	if err != nil {
		return ""
	}

	return strings.ReplaceAll(string(b), "\n", "\\n")
}
