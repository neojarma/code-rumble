package submission_repository

import (
	"errors"
	"server/entity"
	"server/entity/join_model"

	"gorm.io/gorm"
)

type SubmissionRepositoryImpl struct {
	DB *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &SubmissionRepositoryImpl{
		DB: db,
	}
}

func (r *SubmissionRepositoryImpl) NewSubmission(req *entity.Submission) error {
	return r.DB.Create(req).Error
}

func (r *SubmissionRepositoryImpl) GetSubmission(id string) (*join_model.SubmissionTestResult, error) {
	model := new(entity.Submission)

	rows, err := r.DB.Model(model).Select("submissions.submission_id, submissions.question_id, submissions.status, test_cases.input, test_cases.output, test_results.actual_output, test_results.result").Joins("JOIN test_results on test_results.submission_id = submissions.submission_id").Joins("JOIN test_cases on test_cases.test_case_id = test_results.test_case_id").Where("submissions.submission_id = ?", id).Rows()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	queryRes := new(join_model.SubmissionResult)
	finalRes := new(join_model.SubmissionTestResult)

	counter := 0
	for rows.Next() {
		err := r.DB.ScanRows(rows, queryRes)
		if err != nil {
			return nil, err
		}

		if counter == 0 {
			finalRes.QuestionId = queryRes.QuestionId
			finalRes.SubmissionId = queryRes.SubmissionId
			finalRes.SubmissionStatus = queryRes.Status
		}

		finalRes.TestResult = append(finalRes.TestResult, &join_model.SubmissionResult{
			Input:        queryRes.Input,
			ActualOutput: queryRes.ActualOutput,
			Output:       queryRes.Output,
			Result:       queryRes.Result,
		})

		counter++
	}

	if counter == 0 {
		return nil, errors.New("there is no record with that id")
	}

	return finalRes, nil
}

func (r *SubmissionRepositoryImpl) UpdateSubmissionProgres(req *entity.Submission) error {
	return r.DB.Model(req).Where("submission_id = ?", req.SubmissionId).Updates(req).Error
}
