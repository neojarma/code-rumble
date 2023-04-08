package submission_handler

import (
	"net/http"
	"server/entity"
	"server/use_case/submission_use_case"
	"strconv"

	"github.com/labstack/echo/v4"
)

type SubmissionHandlerImpl struct {
	SubmissionUseCase submission_use_case.SubmissionUseCase
}

func NewSubmissionHandler(uc submission_use_case.SubmissionUseCase) SubmissionHandler {
	return &SubmissionHandlerImpl{
		SubmissionUseCase: uc,
	}
}

func (h *SubmissionHandlerImpl) SubmitCode(c echo.Context) error {
	request := new(entity.SubmissionPayload)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, entity.Response{
			Message: "invalid body request",
		})
	}

	id, err := h.SubmissionUseCase.NewSubmission(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.Response{
			Message: "failed to submit code",
		})
	}

	return c.JSON(http.StatusCreated, entity.Response{
		Message: "success submit code",
		Data:    id,
	})
}

func (h *SubmissionHandlerImpl) GetSubmissionData(c echo.Context) error {
	id := c.QueryParam("id")
	isCustomTest, err := strconv.ParseBool(c.QueryParam("custom-test"))
	if id == "" {
		return c.JSON(http.StatusBadRequest, entity.Response{
			Message: "invalid question id",
		})
	}

	if err != nil {
		isCustomTest = false
	}

	res, err := h.SubmissionUseCase.GetSubmission(id, isCustomTest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.Response{
			Message: "failed to get data",
		})
	}

	return c.JSON(http.StatusOK, entity.Response{
		Message: "success to get data",
		Data:    res,
	})
}
