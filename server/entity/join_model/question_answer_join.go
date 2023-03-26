package join_model

import "server/entity"

type QuestionAswer struct {
	Question *entity.Question `json:"question" gorm:"embedded"`
	StubCode string           `json:"stubCode,omitempty"`
}
