package answer_repository

import (
	"server/entity"

	"gorm.io/gorm"
)

type AnswerRepositoryImpl struct {
	DB *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) AnswerRepository {
	return &AnswerRepositoryImpl{
		DB: db,
	}
}

func (repo *AnswerRepositoryImpl) CreateAnswer(model *entity.Answer) error {
	return repo.DB.Create(model).Error
}
