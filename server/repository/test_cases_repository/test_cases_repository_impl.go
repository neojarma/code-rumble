package test_cases_repository

import (
	"server/entity"

	"gorm.io/gorm"
)

type TestCaseRepositoryImpl struct {
	DB *gorm.DB
}

func NewTestCaseRepository(db *gorm.DB) TestCaseRepository {
	return &TestCaseRepositoryImpl{
		DB: db,
	}
}

func (repo *TestCaseRepositoryImpl) CreateNewTestCase(model *entity.TestCase) error {
	return repo.DB.Create(model).Error
}

func (repo *TestCaseRepositoryImpl) CreateBulkTestCase(req []*entity.TestCase) error {
	return repo.DB.Create(&req).Error
}
