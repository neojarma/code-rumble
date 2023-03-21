package test_use_case

import (
	"server/entity"
	"server/helper"
	"server/repository/test_cases_repository"
)

type TestUseCaseImpl struct {
	Repo test_cases_repository.TestCaseRepository
}

func NewTestUseCase(repo test_cases_repository.TestCaseRepository) TestUseCase {
	return &TestUseCaseImpl{
		Repo: repo,
	}
}

func (useCase *TestUseCaseImpl) CreateNewTestCase(model *entity.TestCase) error {
	id := helper.GenerateId(15)
	model.TestCaseId = id
	return useCase.Repo.CreateNewTestCase(model)
}
