package repository

import "server/entity"

type ServerRepository interface {
	GetQuestions(limit, offset int) ([]*entity.Question, error)
	GetQuestionById(id string) (*entity.Question, error)
}
