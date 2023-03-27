package test_repository

import (
	"executor/entity"
	"server/entity/join_model"
)

type TestResultRepository interface {
	CreateTestResult(results []*entity.TestResult) error
	CreateCustomTestResult(results []*entity.CustomTestResult) error
	GetCustomTestResult(submissionId string) (*join_model.SubmissionTestResult, error)
}
