package question_use_case

import (
	"server/entity"
	"server/entity/join_model"
)

type QuestionUseCase interface {
	CreateNewQuestion(req *entity.QuestionPayload) error
	GetQuestions(limit, offset int) ([]*join_model.QuestionAswer, error)
	GetQuestionById(id string) (*join_model.QuestionAswer, error)
	GetQuestionAndTestCase(id string, limit int) (*join_model.QuestionAnswerMultiTest, error)
}
