package submission_repository

import (
	"server/entity"

	"server/entity/join_model"
)

type SubmissionRepository interface {
	NewSubmission(req *entity.Submission) error
	GetSubmission(id string) (*join_model.SubmissionTestResult, error)
	UpdateSubmissionProgres(req *entity.Submission) error
}
