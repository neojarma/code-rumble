package test_repository

import (
	"executor/entity"

	"gorm.io/gorm"
)

type TestResultRepositoryImpl struct {
	DB *gorm.DB
}

func NewTestResult(db *gorm.DB) TestResultRepository {
	return &TestResultRepositoryImpl{
		DB: db,
	}
}

func (r *TestResultRepositoryImpl) CreateTestResult(results []*entity.TestResult) error {
	return r.DB.Create(&results).Error
}
