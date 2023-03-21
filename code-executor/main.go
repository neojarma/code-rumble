package main

import (
	"executor/entity"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v4"
)

func main_code() {
	e := echo.New()

	e.POST("/submission", func(c echo.Context) error {
		req := new(entity.Submission)
		if err := c.Bind(req); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		err := os.WriteFile(fmt.Sprintf("./js-executor/submitted-code/%s.js", req.SubmissionId), []byte(req.SubmittedCode), 0644)
		if err != nil {
			log.Println("error 1")
			return c.String(http.StatusInternalServerError, err.Error())
		}

		err = os.WriteFile(fmt.Sprintf("./js-executor/run-code/%s.js", req.SubmissionId), []byte(req.RunCode), 0644)
		if err != nil {
			log.Println("error 1")
			return c.String(http.StatusInternalServerError, err.Error())
		}

		cmd := fmt.Sprintf("node ./js-executor/run-code/%s.js %s %v", req.SubmissionId, req.SubmissionId, req.ToString())
		_, err = exec.Command("cmd", "/C", cmd).Output()
		if err != nil {
			fmt.Println(fmt.Sprint(err))
		}

		// read file result and return json
		b, err := os.ReadFile(fmt.Sprintf("./js-executor/result-code/%s.txt", req.SubmissionId))
		if err != nil {
			fmt.Println(fmt.Sprint(err))
		}

		resLine := strings.Split(string(b), "\n")
		finalResult := make([]*entity.TestResult, len(resLine))

		for i := 0; i < len(resLine); i++ {
			each := strings.Split(resLine[i], "=")
			obj := &entity.TestResult{
				Input:          req.TestCases[i].Input,
				ExpectedOutput: req.TestCases[i].Output,
				ActualOutput:   each[2],
				Status:         each[1],
				CaseNumber:     each[0],
			}

			log.Println(obj)

			finalResult[i] = obj
		}

		return c.JSON(http.StatusOK, finalResult)
	})

	e.Logger.Fatal(e.Start(":8083"))
}

func main() {
	main_code()
}
