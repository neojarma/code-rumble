package join_model

import "server/entity"

type QuestionAnswerTest struct {
	Question *entity.Question `json:"question" gorm:"embedded"`
	StubCode string           `json:"stubCode,omitempty"`
	Test     *entity.TestCase `json:"testCases" gorm:"embedded"`
}

type QuestionAnswerMultiTest struct {
	Question *entity.Question   `json:"question"`
	StubCode string             `json:"stubCode,omitempty"`
	Test     []*entity.TestCase `json:"testCases"`
}
