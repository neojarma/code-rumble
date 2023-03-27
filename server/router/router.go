package router

import (
	custom_test_repo "executor/repository/test_result_repository"
	"server/handler"
	question_repository "server/repository/questions_repository"
	"server/repository/submission_repository"
	"server/repository/test_cases_repository"
	generate_code "server/use_case/generate_code_use_case"
	question_use_case "server/use_case/questions_use_case"
	"server/use_case/submission_use_case"
	"server/use_case/test_use_case"

	"github.com/labstack/echo/v4"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

func Router(db *gorm.DB, e *echo.Echo, rabbitConn *amqp091.Connection) {

	customTestRepo := custom_test_repo.NewTestResult(db)
	testRepo := test_cases_repository.NewTestCaseRepository(db)
	testUseCase := test_use_case.NewTestUseCase(testRepo)
	generateCode := generate_code.NewGenerateCode()
	submissionRepo := submission_repository.NewSubmissionRepository(db)
	questionRepo := question_repository.NewQuestionRepository(db)
	questionUseCase := question_use_case.NewQuestionUseCase(questionRepo, testUseCase, generateCode, rabbitConn)
	submissionUseCase := submission_use_case.NewSubmissionUseCase(submissionRepo, questionUseCase, rabbitConn, customTestRepo)

	handler := handler.NewServerHandler(questionUseCase, submissionUseCase)

	e.GET("/questions", handler.GetQuestions)
	e.GET("/question", handler.GetQuestionById)
	e.POST("/question", handler.CreateNewQuestion)
	e.GET("/question/case", handler.GetQuestionAndTestCase)
	e.POST("/submission", handler.SubmitCode)
	e.GET("/submission", handler.GetSubmissionData)
}
