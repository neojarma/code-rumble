package test_repository

import (
	"executor/entity"
	"server/entity/join_model"

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

func (r *TestResultRepositoryImpl) CreateCustomTestResult(results []*entity.CustomTestResult) error {
	return r.DB.Create(&results).Error
}

func (r *TestResultRepositoryImpl) GetCustomTestResult(submissionId string) (*join_model.SubmissionTestResult, error) {
	res := new(join_model.SubmissionTestResult)
	model := new(entity.CustomTestResult)

	rows, err := r.DB.Where("submission_id = ?", submissionId).Find(model).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		err := r.DB.ScanRows(rows, model)
		if err != nil {
			return nil, err
		}

		res.SubmissionId = submissionId
		res.TestResult = append(res.TestResult, &join_model.SubmissionResult{
			Input:        model.Input,
			Output:       model.ExpectedOutput,
			ActualOutput: model.ActualOutput,
			Result:       model.Result,
		})
	}

	return res, nil
}
