package answer_repository

import "server/entity"

type AnswerRepository interface {
	CreateAnswer(model *entity.Answer) error
}
