package router

import (
	custom_test_repo "executor/repository/test_result_repository"
	leaderboard_handler "server/handler/leaderboard_game_handler"
	"server/handler/question_handler"
	question_repository "server/repository/questions_repository"
	"server/repository/submission_repository"
	"server/repository/test_cases_repository"
	generate_code "server/use_case/generate_code_use_case"
	question_use_case "server/use_case/questions_use_case"
	"server/use_case/submission_use_case"
	"server/use_case/test_use_case"

	"server/handler/submission_handler"

	"github.com/labstack/echo/v4"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Router(db *gorm.DB, e *echo.Echo, rabbitConn *amqp091.Connection, redisConn *redis.Client) {

	customTestRepo := custom_test_repo.NewTestResult(db)
	testRepo := test_cases_repository.NewTestCaseRepository(db)
	testUseCase := test_use_case.NewTestUseCase(testRepo)
	generateCode := generate_code.NewGenerateCode()
	submissionRepo := submission_repository.NewSubmissionRepository(db)
	questionRepo := question_repository.NewQuestionRepository(db)
	questionUseCase := question_use_case.NewQuestionUseCase(questionRepo, testUseCase, generateCode, rabbitConn)
	submissionUseCase := submission_use_case.NewSubmissionUseCase(submissionRepo, questionUseCase, rabbitConn, customTestRepo)

	questionHandler := question_handler.NewQuestionHandler(questionUseCase)
	submissionHandler := submission_handler.NewSubmissionHandler(submissionUseCase)

	leaderboardHandler := leaderboard_handler.NewLeaderboardHandler(questionUseCase, submissionUseCase, redisConn)

	e.GET("/questions", questionHandler.GetQuestions)
	e.GET("/question", questionHandler.GetQuestionById)
	e.POST("/question", questionHandler.CreateNewQuestion)
	e.GET("/question/case", questionHandler.GetQuestionAndTestCase)
	e.POST("/submission", submissionHandler.SubmitCode)
	e.GET("/submission", submissionHandler.GetSubmissionData)
	e.GET("/create-room", leaderboardHandler.CreateRoom)
	e.GET("/join-room", leaderboardHandler.JoinRoom)
	e.GET("/start", leaderboardHandler.StartGame)
	e.GET("/submit-code", leaderboardHandler.SubmitCode)
	e.GET("/leaderboard", leaderboardHandler.GetLeaderboard)
}
