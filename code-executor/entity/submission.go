package entity

import "encoding/json"

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

	return string(b)
}

type TestCase struct {
	TestCaseId string `json:"id,omitempty"`
	QuestionId string `json:"-"`
	Input      string `json:"input"`
	Output     string `json:"output"`
}

type TestResult struct {
	CaseNumber     string `json:"caseNumber"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expectedOutput"`
	ActualOutput   string `json:"actualOutput"`
	Status         string `json:"status"`
}
