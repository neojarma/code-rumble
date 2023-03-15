package repository

import (
	"server/entity"

	"server/helper"

	"gorm.io/gorm"
)

type ServerRepositoryImpl struct {
	DB *gorm.DB
}

func NewServerRepository(db *gorm.DB) ServerRepository {
	return &ServerRepositoryImpl{
		DB: db,
	}
}

func (r *ServerRepositoryImpl) GetQuestions(limit int, offset int) ([]*entity.Question, error) {
	model := new(entity.Question)
	paging := helper.Pagination(limit, offset)
	res := make([]*entity.Question, 0)

	if err := r.DB.Model(&model).Scopes(paging).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (r *ServerRepositoryImpl) GetQuestionById(id string) (*entity.Question, error) {
	model := &entity.Question{
		ID: id,
	}

	if err := r.DB.Find(model).Error; err != nil {
		return nil, err
	}

	return model, nil
}
