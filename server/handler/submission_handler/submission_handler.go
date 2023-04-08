package submission_handler

import "github.com/labstack/echo/v4"

type SubmissionHandler interface {
	SubmitCode(c echo.Context) error
	GetSubmissionData(c echo.Context) error
}
