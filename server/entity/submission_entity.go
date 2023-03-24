package entity

type Submission struct {
	SubmissionId  string `json:"submissionId"`
	SubmittedCode string `json:"submittedCode"`
	QuestionId    string `json:"questionId"`
	Status        string `json:"status" gorm:"default:PROCESS"`
}
