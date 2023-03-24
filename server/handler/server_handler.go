package handler

import "github.com/labstack/echo/v4"

type ServerHandler interface {
	GetQuestions(c echo.Context) error
	GetQuestionById(c echo.Context) error
	GetQuestionAndTestCase(c echo.Context) error
	CreateNewQuestion(c echo.Context) error
	SubmitCode(c echo.Context) error
	GetSubmissionData(c echo.Context) error
}
