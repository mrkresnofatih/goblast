package goblast

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Description: for constructing a base model for http response body
type BaseResponseModel[T interface{}] struct {
	Data         T      `json:"data"`
	ErrorMessage string `json:"errorMessage"`
}

// Description: a base method for returning an OK-200 Response
func OkResponse[T interface{}](c echo.Context, data T) error {
	return c.JSON(http.StatusOK, BaseResponseModel[T]{
		Data:         data,
		ErrorMessage: "",
	})
}

// Description: a base method for returning a 400-Bad Response
func NotOkResponse(c echo.Context, errorMessage string) error {
	return c.JSON(http.StatusBadRequest, BaseResponseModel[interface{}]{
		Data:         new(struct{}),
		ErrorMessage: errorMessage,
	})
}
