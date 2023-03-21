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

func (repo *QuestionRepositoryImpl) GetQuestions(limit int, offset int) ([]*join_model.QuestionAswer, error) {
	question := new(entity.Question)
	result := make([]*join_model.QuestionAswer, 0)
	paging := helper.Pagination(limit, offset)

	err := repo.DB.Model(question).Select("questions.question_id, questions.description, answers.answer_id, answers.template_code").Joins("JOIN answers ON questions.question_id = answers.question_id").Scopes(paging).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (repo *QuestionRepositoryImpl) GetQuestionById(id string) (*join_model.QuestionAswer, error) {
	question := new(entity.Question)
	result := new(join_model.QuestionAswer)

	queryResult := repo.DB.Model(question).Select("questions.question_id, questions.description, answers.answer_id, answers.template_code").Joins("JOIN answers ON questions.question_id = answers.question_id").Where("questions.question_id = ?", id).Find(&result)
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

	r, err := repo.DB.Model(question).Select("questions.question_id, questions.description, answers.answer_id, answers.template_code, test_cases.test_case_id, test_cases.input, test_cases.output").Joins("JOIN answers ON questions.question_id = answers.question_id").Joins("JOIN test_cases ON questions.question_id = test_cases.question_id").Where("questions.question_id = ?", id).Limit(limit).Rows()
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
			}

			result.Answer = &entity.Answer{
				AnswerID:     queryRes.Answer.AnswerID,
				TemplateCode: queryRes.Answer.TemplateCode,
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
