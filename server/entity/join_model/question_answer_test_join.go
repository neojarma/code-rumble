package join_model

import "server/entity"

type QuestionAnswerTest struct {
	Question *entity.Question `json:"question" gorm:"embedded"`
	Answer   *entity.Answer   `json:"answer" gorm:"embedded"`
	Test     *entity.TestCase `json:"testCases" gorm:"embedded"`
}

type QuestionAnswerMultiTest struct {
	Question *entity.Question   `json:"question"`
	Answer   *entity.Answer     `json:"answer"`
	Test     []*entity.TestCase `json:"testCases"`
}
