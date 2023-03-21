package router

import (
	"server/handler"
	"server/repository/answer_repository"
	question_repository "server/repository/questions_repository"
	"server/repository/test_cases_repository"
	"server/use_case/answer_use_case"
	question_use_case "server/use_case/questions_use_case"
	"server/use_case/test_use_case"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Router(db *gorm.DB, e *echo.Echo) {

	answerRepo := answer_repository.NewAnswerRepository(db)
	answerUseCase := answer_use_case.NewAnswerUseCase(answerRepo)

	testRepo := test_cases_repository.NewTestCaseRepository(db)
	testUseCase := test_use_case.NewTestUseCase(testRepo)

	questionRepo := question_repository.NewQuestionRepository(db)
	questionUseCase := question_use_case.NewQuestionUseCase(questionRepo, answerUseCase, testUseCase)
	handler := handler.NewServerHandler(questionUseCase)

	e.GET("/questions", handler.GetQuestions)
	e.GET("/question", handler.GetQuestionById)
	e.POST("/question", handler.CreateNewQuestion)
	e.GET("/question/case", handler.GetQuestionAndTestCase)
	e.POST("/submission", handler.SubmitCode)
}
