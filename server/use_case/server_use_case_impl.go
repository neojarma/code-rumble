package use_case

import (
	"server/entity"
	"server/repository"
)

type ServerUseCaseImpl struct {
	Repository repository.ServerRepository
}

func NewServerUseCase(r repository.ServerRepository) ServerUseCase {
	return &ServerUseCaseImpl{
		Repository: r,
	}
}

func (u *ServerUseCaseImpl) GetQuestions(limit int, offset int) ([]*entity.Question, error) {
	return u.Repository.GetQuestions(limit, offset)
}

func (u *ServerUseCaseImpl) GetQuestionById(id string) (*entity.Question, error) {
	return u.Repository.GetQuestionById(id)
}
