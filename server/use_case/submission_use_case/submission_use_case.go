package submission_use_case

import (
	"server/entity"
	"server/entity/join_model"
)

type SubmissionUseCase interface {
	NewSubmission(req *entity.SubmissionPayload) (string, error)
	GetSubmission(id string, isCustomTest bool) (*join_model.SubmissionTestResult, error)
	UpdateSubmissionProgres(req *entity.Submission) error
}
