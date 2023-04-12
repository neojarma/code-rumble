package question_repository

import (
	"errors"
	"server/entity"
	"server/entity/join_model"
	"server/helper"

	"gorm.io/gorm"
)

type QuestionRepositoryImpl struct {
	DB *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &QuestionRepositoryImpl{
		DB: db,
	}
}

func (repo *QuestionRepositoryImpl) CreateNewQuestion(model *entity.Question) error {
	return repo.DB.Create(model).Error
}

func (repo *QuestionRepositoryImpl) GetQuestions(limit int, offset int) ([]*entity.Question, error) {
	result := make([]*entity.Question, 0)
	paging := helper.Pagination(limit, offset)

	err := repo.DB.Select("questions.question_id, questions.description, questions.title").Scopes(paging).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (repo *QuestionRepositoryImpl) GetQuestionById(id string) (*entity.Question, error) {
	result := new(entity.Question)

	queryResult := repo.DB.Select("questions.question_id, questions.title, questions.description").Where("questions.question_id = ?", id).Find(&result)
	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	if queryResult.RowsAffected == 0 {
		return nil, errors.New("there is no question with id: " + id)
	}

	return result, nil
}

func (repo *QuestionRepositoryImpl) GetQuestionAndTestCase(id string, limit int) (*join_model.QuestionAnswerMultiTest, error) {
	question := new(entity.Question)
	queryRes := new(join_model.QuestionAnswerTest)
	result := new(join_model.QuestionAnswerMultiTest)

	r, err := repo.DB.Model(question).Select("questions.question_id, questions.title, questions.description, test_cases.test_case_id, test_cases.input, test_cases.output").Joins("JOIN test_cases ON questions.question_id = test_cases.question_id").Where("questions.question_id = ?", id).Limit(limit).Rows()
	if err != nil {
		return nil, err
	}

	defer r.Close()

	counter := 0
	for r.Next() {
		repo.DB.ScanRows(r, queryRes)

		if counter == 0 {
			result.Question = &entity.Question{
				QuestionId:  queryRes.Question.QuestionId,
				Description: queryRes.Question.Description,
				Title:       queryRes.Question.Title,
			}
		}

		result.Test = append(result.Test, &entity.TestCase{
			TestCaseId: queryRes.Test.TestCaseId,
			Input:      queryRes.Test.Input,
			Output:     queryRes.Test.Output,
		})

		counter++
	}

	if counter == 0 {
		return nil, errors.New("there is no record with that id")
	}

	return result, nil
}

func (repo *QuestionRepositoryImpl) GetRandomQuestions(limit int) ([]*entity.Question, error) {
	result := make([]*entity.Question, 0)

	err := repo.DB.Select("questions.question_id, questions.description, questions.title").Order("rand()").Limit(limit).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
