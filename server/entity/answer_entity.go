package entity

type Answer struct {
	AnswerID     string `json:"answerId"`
	QuestionId   string `json:"-"`
	TemplateCode string `json:"function"`
}
