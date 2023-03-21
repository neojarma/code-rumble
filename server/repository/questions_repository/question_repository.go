package question_repository

import (
	"server/entity"
	"server/entity/join_model"
)

type QuestionRepository interface {
	CreateNewQuestion(model *entity.Question) error
	GetQuestions(limit, offset int) ([]*join_model.QuestionAswer, error)
	GetQuestionById(id string) (*join_model.QuestionAswer, error)
	GetQuestionAndTestCase(id string, limit int) (*join_model.QuestionAnswerMultiTest, error)
}
