package handler

import (
	"net/http"
	"server/entity"
	"server/use_case"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ServerHandlerImpl struct {
	UseCase use_case.ServerUseCase
}

func NewServerHandler(useCase use_case.ServerUseCase) ServerHandler {
	return &ServerHandlerImpl{
		UseCase: useCase,
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
