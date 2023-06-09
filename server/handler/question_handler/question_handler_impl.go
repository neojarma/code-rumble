package question_handler

import (
	"log"
	"net/http"
	"server/entity"
	questions_use_case "server/use_case/questions_use_case"
	"strconv"

	"github.com/labstack/echo/v4"
)

type QuestionHandlerImpl struct {
	UseCase questions_use_case.QuestionUseCase
}

func NewQuestionHandler(uc questions_use_case.QuestionUseCase) QuestionHandler {
	return &QuestionHandlerImpl{
		UseCase: uc,
	}
}

func (h *QuestionHandlerImpl) GetQuestions(c echo.Context) error {

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

func (h *QuestionHandlerImpl) GetQuestionById(c echo.Context) error {
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

func (h *QuestionHandlerImpl) GetQuestionAndTestCase(c echo.Context) error {
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

func (h *QuestionHandlerImpl) CreateNewQuestion(c echo.Context) error {
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
