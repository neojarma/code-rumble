package router

import (
	"server/handler"
	"server/repository"
	"server/use_case"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Router(db *gorm.DB, e *echo.Echo) {

	repo := repository.NewServerRepository(db)
	useCase := use_case.NewServerUseCase(repo)
	handler := handler.NewServerHandler(useCase)

	e.GET("/questions", handler.GetQuestions)
	e.GET("/question", handler.GetQuestionById)
}
