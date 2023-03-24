package submission_use_case

import (
	"context"
	"encoding/json"
	"fmt"
	"server/entity"
	"server/entity/join_model"
	"server/helper"
	event_rabbitmq "server/rabbitmq"
	"server/repository/submission_repository"
	question_use_case "server/use_case/questions_use_case"

	"github.com/rabbitmq/amqp091-go"
)

type SubmissionUseCaseImpl struct {
	Repo            submission_repository.SubmissionRepository
	QuestionUseCase question_use_case.QuestionUseCase
	RabbitConn      *amqp091.Connection
}

func NewSubmissionUseCase(r submission_repository.SubmissionRepository, q question_use_case.QuestionUseCase, c *amqp091.Connection) SubmissionUseCase {
	return &SubmissionUseCaseImpl{
		Repo:            r,
		QuestionUseCase: q,
		RabbitConn:      c,
	}
}

func (u *SubmissionUseCaseImpl) NewSubmission(req *entity.SubmissionPayload) (string, error) {
	id := helper.GenerateId(15)
	req.SubmissionId = id

	req.RunCode = fmt.Sprintf("http://127.0.0.1:8082/file?id=%s", req.QuestionId)

	var bodyByte []byte
	var err error

	// post to rabbit mq
	ctx := context.Background()
	if req.CustomTestCase {
		bodyByte, err = json.Marshal(req)

	} else {
		// get test case from db
		res, err := u.QuestionUseCase.GetQuestionAndTestCase(req.QuestionId, -1)
		if err != nil {
			return "", err
		}

		req.TestCases = res.Test
		bodyByte, err = json.Marshal(req)
		if err != nil {
			return "", err
		}
	}

	if err != nil {
		return "", err
	}

	err = event_rabbitmq.EmitEvent(ctx, bodyByte, "submission", u.RabbitConn)
	if err != nil {
		return "", err
	}

	// save to db
	err = u.Repo.NewSubmission(&entity.Submission{
		SubmissionId:  req.SubmissionId,
		SubmittedCode: req.SubmittedCode,
		QuestionId:    req.QuestionId,
	})

	return id, err
}

func (u *SubmissionUseCaseImpl) GetSubmission(id string) (*join_model.SubmissionTestResult, error) {
	return u.Repo.GetSubmission(id)
}

func (u *SubmissionUseCaseImpl) UpdateSubmissionProgres(req *entity.Submission) error {
	return u.Repo.UpdateSubmissionProgres(req)
}
