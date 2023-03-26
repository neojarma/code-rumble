package handler

import (
	"log"
	"net/http"
	"server/entity"
	questions_use_case "server/use_case/questions_use_case"
	"server/use_case/submission_use_case"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ServerHandlerImpl struct {
	UseCase           questions_use_case.QuestionUseCase
	SubmissionUseCase submission_use_case.SubmissionUseCase
}

func NewServerHandler(useCase questions_use_case.QuestionUseCase, s submission_use_case.SubmissionUseCase) ServerHandler {
	return &ServerHandlerImpl{
		UseCase:           useCase,
		SubmissionUseCase: s,
	}
}

func (h *ServerHandlerImpl) GetQuestions(c echo.Context) error {

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 5
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}

	res, err := h.UseCase.GetQuestions(limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.Response{
			Message: "failed to get data",
		})
	}

	return c.JSON(http.StatusOK, entity.Response{
		Message: "succes get data",
		Data:    res,
	})

}

func (h *ServerHandlerImpl) GetQuestionById(c echo.Context) error {

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, entity.Response{
			Message: "please specify question id",
		})
	}

	res, err := h.UseCase.GetQuestionById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.Response{
			Message: "failed to get data",
		})
	}

	return c.JSON(http.StatusOK, entity.Response{
		Message: "succes get data",
		Data:    res,
	})
}

func (h *ServerHandlerImpl) GetQuestionAndTestCase(c echo.Context) error {
	id := c.QueryParam("id")
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if id == "" {
		return c.JSON(http.StatusBadRequest, entity.Response{
			Message: "invalid question id",
		})
	}

	if err != nil {
		limit = -1
	}

	res, err := h.UseCase.GetQuestionAndTestCase(id, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.Response{
			Message: "failed to get data",
		})
	}

	return c.JSON(http.StatusOK, entity.Response{
		Message: "succes get data",
		Data:    res,
	})
}

func (h *ServerHandlerImpl) CreateNewQuestion(c echo.Context) error {
	request := new(entity.QuestionPayload)
	if err := c.Bind(request); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, entity.Response{
			Message: "invalid body request",
		})
	}

	err := h.UseCase.CreateNewQuestion(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.Response{
			Message: "failed to create data",
			Data:    err.Error,
		})
	}

	return c.JSON(http.StatusCreated, entity.Response{
		Message: "succes create data",
	})
}

func (h *ServerHandlerImpl) SubmitCode(c echo.Context) error {
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

func (h *ServerHandlerImpl) GetSubmissionData(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, entity.Response{
			Message: "invalid question id",
		})
	}

	res, err := h.SubmissionUseCase.GetSubmission(id)
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
