package question_use_case

import (
	"server/entity/join_model"
	"server/helper"
	question_repository "server/repository/questions_repository"
	"server/use_case/answer_use_case"
	"server/use_case/test_use_case"
)

type QuestionUseCaseImpl struct {
	QuestionRepo  question_repository.QuestionRepository
	AnswerUseCase answer_use_case.AnswerUseCase
	TestUseCase   test_use_case.TestUseCase
}

func NewQuestionUseCase(qR question_repository.QuestionRepository, aU answer_use_case.AnswerUseCase, tU test_use_case.TestUseCase) QuestionUseCase {
	return &QuestionUseCaseImpl{
		QuestionRepo:  qR,
		AnswerUseCase: aU,
		TestUseCase:   tU,
	}
}

func (useCase *QuestionUseCaseImpl) CreateNewQuestion(model *join_model.QuestionAnswerTest) error {
	// insert data to questions table
	id := helper.GenerateId(15)
	model.Question.QuestionId = id
	if err := useCase.QuestionRepo.CreateNewQuestion(model.Question); err != nil {
		return err
	}

	// insert data to answer table
	if err := useCase.AnswerUseCase.CreateAnswer(model.Answer); err != nil {
		return err
	}

	// insert data to testcases table
	// for _, v := range model.Test {
	// 	if err := useCase.TestUseCase.CreateNewTestCase(v); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (useCase *QuestionUseCaseImpl) GetQuestions(limit int, offset int) ([]*join_model.QuestionAswer, error) {
	return useCase.QuestionRepo.GetQuestions(limit, offset)
}

func (useCase *QuestionUseCaseImpl) GetQuestionById(id string) (*join_model.QuestionAswer, error) {
	return useCase.QuestionRepo.GetQuestionById(id)
}

func (useCase *QuestionUseCaseImpl) GetQuestionAndTestCase(id string, limit int) (*join_model.QuestionAnswerMultiTest, error) {
	return useCase.QuestionRepo.GetQuestionAndTestCase(id, limit)
}
