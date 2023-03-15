package use_case

import "server/entity"

type ServerUseCase interface {
	GetQuestions(limit, offset int) ([]*entity.Question, error)
	GetQuestionById(id string) (*entity.Question, error)
}
