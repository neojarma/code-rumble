package handler

import "github.com/labstack/echo/v4"

type ServerHandler interface {
	GetQuestions(c echo.Context) error
	GetQuestionById(c echo.Context) error
}
