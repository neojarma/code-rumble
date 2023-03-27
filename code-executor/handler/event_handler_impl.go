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

	err := generateFiles(event)
	if err != nil {
		log.Println(err)
		return
	}

	cmd := fmt.Sprintf("node ./js-executor/run-code/%s.js %s", event.SubmissionId, event.SubmissionId)
	if runtime.GOOS == "windows" {
		_, err = exec.Command("cmd", "/C", cmd).Output()
	} else {
		_, err = exec.Command("bash", "-c", cmd).Output()
	}
	if err != nil {
		log.Println("error executing command", err)
		return
	}

	// cmde := exec.Command("cmd", "/C", cmd)
	// var out bytes.Buffer
	// var stderr bytes.Buffer
	// cmde.Stdout = &out
	// cmde.Stderr = &stderr
	// err = cmde.Run()
	// fmt.Println(fmt.Sprint(out) + ": " + stderr.String())
	// fmt.Println(fmt.Sprint(err) + ": " + stderr.String())

	if event.CustomTestCase {
		res, err := parsingCustomResultFile(event)
		if err != nil {
			log.Println("error parsing file", err)
			return
		}

		// write to db
		err = h.UseCase.CreateCustomTestResult(res)
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
	} else {
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

}

func generateFiles(payload *entity.Submission) error {
	err := os.WriteFile(fmt.Sprintf("./js-executor/submitted-code/%s.js", payload.SubmissionId), []byte(payload.SubmittedCode), 0644)
	if err != nil {
		return err
	}

	byteTest, err := getTestCode(payload.RunCode)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("./js-executor/run-code/%s.js", payload.SubmissionId), byteTest, 0644)
	if err != nil {
		return err
	}

	byteTestCase, err := payload.MarshallTestCase()
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("./js-executor/test-cases/%s.json", payload.SubmissionId), byteTestCase, 0644)
	if err != nil {
		return err
	}

	return nil
}

func parsingCustomResultFile(event *entity.Submission) ([]*entity.CustomTestResult, error) {
	b, err := os.ReadFile(fmt.Sprintf("./js-executor/result-code/%s.txt", event.SubmissionId))
	if err != nil {
		return nil, err
	}

	resLine := strings.Split(string(b), "\n")
	finalResult := make([]*entity.CustomTestResult, len(resLine))

	for i := 0; i < len(finalResult); i++ {
		each := strings.Split(resLine[i], "=")
		res := each[1]
		actualOutput := each[2]
		expectedOutput := each[3]

		obj := &entity.CustomTestResult{
			SubmissionId:       event.SubmissionId,
			CustomTestResultId: helper.GenerateId(15),
			ActualOutput:       actualOutput,
			Result:             res,
			Input:              event.TestCases[i].Input,
			ExpectedOutput:     expectedOutput,
		}

		finalResult[i] = obj
	}

	return finalResult, nil

}

func parsingResultFile(event *entity.Submission) ([]*entity.TestResult, error) {
	b, err := os.ReadFile(fmt.Sprintf("./js-executor/result-code/%s.txt", event.SubmissionId))
	if err != nil {
		return nil, err
	}

	resLine := strings.Split(string(b), "\n")
	finalResult := make([]*entity.TestResult, len(resLine))

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
