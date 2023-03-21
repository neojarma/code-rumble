package answer_use_case

import (
	"server/entity"
	"server/helper"
	"server/repository/answer_repository"
)

type AnswerUseCaseImpl struct {
	Repo answer_repository.AnswerRepository
}

func NewAnswerUseCase(repo answer_repository.AnswerRepository) AnswerUseCase {
	return &AnswerUseCaseImpl{
		Repo: repo,
	}
}

func (useCase *AnswerUseCaseImpl) CreateAnswer(model *entity.Answer) error {
	id := helper.GenerateId(15)
	model.AnswerID = id
	return useCase.Repo.CreateAnswer(model)
}
