package handler

import (
	"executor/entity"
	test_use_case "executor/use_case/test_result_use_case"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	serv "server/entity"
	"server/helper"
	"server/repository/submission_repository"
	"strings"
)

type EventHandlerImpl struct {
	UseCase        test_use_case.TestResultUseCase
	SubmissionRepo submission_repository.SubmissionRepository
}

func NewEventHandler(useCase test_use_case.TestResultUseCase, sr submission_repository.SubmissionRepository) EventHandler {
	return &EventHandlerImpl{
		UseCase:        useCase,
		SubmissionRepo: sr,
	}
}

func (h *EventHandlerImpl) HandleEvent(event *entity.Submission) {
	err := os.WriteFile(fmt.Sprintf("./js-executor/submitted-code/%s.js", event.SubmissionId), []byte(event.SubmittedCode), 0644)
	if err != nil {
		log.Println("failed writing submitted code", err)
		return
	}

	byteTest, err := getTestCode(event.RunCode)
	if err != nil {
		log.Println("error getting file", err)
		return
	}

	err = os.WriteFile(fmt.Sprintf("./js-executor/run-code/%s.js", event.SubmissionId), byteTest, 0644)
	if err != nil {
		log.Println("failed writing run code", err)
		return
	}

	cmd := fmt.Sprintf("node ./js-executor/run-code/%s.js %s %v", event.SubmissionId, event.SubmissionId, event.ToString())

	if runtime.GOOS == "windows" {
		_, err = exec.Command("cmd", "/C", cmd).Output()
	} else {
		_, err = exec.Command("bash", "-c", cmd).Output()
	}
	if err != nil {
		log.Println("error executing command", err)
		return
	}

	res, err := parsingResultFile(event)
	if err != nil {
		log.Println("error parsing file", err)
		return
	}

	// write to db
	err = h.UseCase.CreateTestResult(res)
	if err != nil {
		log.Println("error writing test result to db", err)
		return
	}

	// update submission table
	err = h.SubmissionRepo.UpdateSubmissionProgres(&serv.Submission{
		SubmissionId: event.SubmissionId,
		Status:       "FINISHED",
	})
	if err != nil {
		log.Println("error writing test result to db", err)
		return
	}
}

func parsingResultFile(event *entity.Submission) ([]*entity.TestResult, error) {
	b, err := os.ReadFile(fmt.Sprintf("./js-executor/result-code/%s.txt", event.SubmissionId))
	if err != nil {
		return nil, err
	}

	resLine := strings.Split(string(b), "\n")
	finalResult := make([]*entity.TestResult, len(resLine)-1)

	for i := 0; i < len(finalResult); i++ {
		each := strings.Split(resLine[i], "=")
		id := each[0]
		res := each[1]
		out := each[2]

		obj := &entity.TestResult{
			SubmissionId: event.SubmissionId,
			TestCaseId:   id,
			ActualOutput: out,
			Result:       res,
			TestResultId: helper.GenerateId(15),
		}

		finalResult[i] = obj
	}

	return finalResult, nil
}

func getTestCode(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	content := make([]byte, resp.ContentLength)
	_, err = resp.Body.Read(content)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return content, nil

}
