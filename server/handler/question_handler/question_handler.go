package question_handler

import "github.com/labstack/echo/v4"

type QuestionHandler interface {
	GetQuestions(c echo.Context) error
	GetQuestionById(c echo.Context) error
	GetQuestionAndTestCase(c echo.Context) error
	CreateNewQuestion(c echo.Context) error
}
