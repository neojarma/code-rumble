package answer_use_case

import "server/entity"

type AnswerUseCase interface {
	CreateAnswer(model *entity.Answer) error
}
