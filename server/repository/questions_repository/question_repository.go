package question_repository

import (
	"server/entity"
	"server/entity/join_model"
)

type QuestionRepository interface {
	CreateNewQuestion(model *entity.Question) error
	GetQuestions(limit, offset int) ([]*entity.Question, error)
	GetQuestionById(id string) (*entity.Question, error)
	GetQuestionAndTestCase(id string, limit int) (*join_model.QuestionAnswerMultiTest, error)
	GetRandomQuestions(limit int) ([]*entity.Question, error)
}
