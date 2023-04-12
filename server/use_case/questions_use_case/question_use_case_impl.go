package question_use_case

import (
	"context"
	"encoding/json"
	fs "files_storage/entity"
	"fmt"
	"server/entity"
	"server/entity/join_model"
	"server/helper"
	event_rabbitmq "server/rabbitmq"
	question_repository "server/repository/questions_repository"
	generate_code "server/use_case/generate_code_use_case"
	"server/use_case/test_use_case"

	"github.com/rabbitmq/amqp091-go"
)

type QuestionUseCaseImpl struct {
	QuestionRepo question_repository.QuestionRepository
	TestUseCase  test_use_case.TestUseCase
	GenerateCode generate_code.GenerateCode
	RabbitConn   *amqp091.Connection
}

func NewQuestionUseCase(qR question_repository.QuestionRepository, tU test_use_case.TestUseCase, gC generate_code.GenerateCode, rC *amqp091.Connection) QuestionUseCase {
	return &QuestionUseCaseImpl{
		QuestionRepo: qR,
		TestUseCase:  tU,
		GenerateCode: gC,
		RabbitConn:   rC,
	}
}

func (useCase *QuestionUseCaseImpl) CreateNewQuestion(req *entity.QuestionPayload) error {
	newQuestionId := helper.GenerateId(15)

	err := useCase.QuestionRepo.CreateNewQuestion(&entity.Question{
		QuestionId:  newQuestionId,
		Description: req.QuestionDescription,
		Title:       req.QuestionTitle,
	})
	if err != nil {
		return err
	}

	err = useCase.TestUseCase.CreateBulkTestCase(req.TestCase, newQuestionId)
	if err != nil {
		return err
	}

	runCode, stubCode, err := useCase.GenerateCode.GenerateJavaScriptCode(req.CodeStub)
	if err != nil {
		return err
	}

	b, err := json.Marshal(fs.Question{
		QuestionId: newQuestionId,
		RunCode:    string(runCode),
		StubCode:   string(stubCode),
	})
	if err != nil {
		return err
	}

	return event_rabbitmq.EmitEvent(context.Background(), b, "new-question-file", useCase.RabbitConn)
}

func (useCase *QuestionUseCaseImpl) GetQuestions(limit int, offset int) ([]*join_model.QuestionAswer, error) {
	res, err := useCase.QuestionRepo.GetQuestions(limit, offset)
	if err != nil {
		return nil, err
	}

	questionAnswer := make([]*join_model.QuestionAswer, len(res))

	for i, v := range res {
		b, err := helper.GetStubCode(fmt.Sprintf("http://127.0.0.1:8082/file?id=%s_answer", v.QuestionId))
		if err != nil {
			return nil, err
		}

		questionAnswer[i] = &join_model.QuestionAswer{
			Question: v,
			StubCode: string(b),
		}
	}

	return questionAnswer, nil
}

func (useCase *QuestionUseCaseImpl) GetQuestionById(id string) (*join_model.QuestionAswer, error) {
	res, err := useCase.QuestionRepo.GetQuestionById(id)
	if err != nil {
		return nil, err
	}

	b, err := helper.GetStubCode(fmt.Sprintf("http://127.0.0.1:8082/file?id=%s_answer", res.QuestionId))
	if err != nil {
		return nil, err
	}

	questionAnswer := &join_model.QuestionAswer{
		Question: res,
		StubCode: string(b),
	}

	return questionAnswer, nil
}

func (useCase *QuestionUseCaseImpl) GetQuestionAndTestCase(id string, limit int) (*join_model.QuestionAnswerMultiTest, error) {
	return useCase.QuestionRepo.GetQuestionAndTestCase(id, limit)
}

func (useCase *QuestionUseCaseImpl) GetRandomQuestions(limit int) ([]*entity.RandQuestion, error) {
	res, err := useCase.QuestionRepo.GetRandomQuestions(limit)
	if err != nil {
		return nil, err
	}

	questionAnswer := make([]*entity.RandQuestion, len(res))

	for i, v := range res {
		b, err := helper.GetStubCode(fmt.Sprintf("http://127.0.0.1:8082/file?id=%s_answer", v.QuestionId))
		if err != nil {
			return nil, err
		}

		questionAnswer[i] = &entity.RandQuestion{
			Question: v,
			StubCode: string(b),
		}
	}

	return questionAnswer, nil
}
