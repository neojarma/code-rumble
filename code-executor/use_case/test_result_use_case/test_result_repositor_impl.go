package test_use_case

import (
	"executor/entity"
	test_repository "executor/repository/test_result_repository"
)

type TestResultUseCaseImpl struct {
	Repo test_repository.TestResultRepository
}

func NewTestResult(repo test_repository.TestResultRepository) TestResultUseCase {
	return &TestResultUseCaseImpl{
		Repo: repo,
	}
}

func (r *TestResultUseCaseImpl) CreateTestResult(results []*entity.TestResult) error {
	return r.Repo.CreateTestResult(results)
}
