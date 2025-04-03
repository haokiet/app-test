package http

import (
	"net/http"
	"strconv"

	"github.com/Vantuan1606/app-buff/domain"
	"github.com/Vantuan1606/app-buff/user"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHTTPHandler struct {
	userUsecase domain.IUserUsecase
}

type responseErr struct {
	Error domain.Error `json:"error"`
}

type responseUser struct {
	User interface{} `json:"user"`
}

type responseUsers struct {
	Users interface{} `json:"users"`
}

func NewUserHTTPHandler(e *echo.Echo, us domain.IUserUsecase) {
	handler := &UserHTTPHandler{
		userUsecase: us,
	}

	e.GET("/user", handler.List)

}

func (uss *UserHTTPHandler) List(c echo.Context) error {
	offset := c.QueryParam("offset")
	offetInt, _ := strconv.ParseInt(offset, 10, 64)
	limit := c.QueryParam("limit")
	limitInt, _ := strconv.ParseInt(limit, 10, 64)
	sort := c.QueryParam("sort")

	ascendingStr := c.QueryParam("ascending")
	ascending := false
	if ascendingStr == "true" {
		ascending = true
	}

	ctx := c.Request().Context()

	input := user.ListUserInput{}

	if sort != "" {
		input.SetSort(sort)
	}
	input.SetAscending(ascending)
	input.SetLimit(limitInt)
	input.SetOffset(offetInt)

	users, err := uss.userUsecase.List(ctx, &input)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, &responseErr{
				Error: domain.Error{
					Code:    http.StatusNotFound,
					Message: err.Error(),
					Type:    "UserException",
				},
			})
		}

		return c.JSON(http.StatusInternalServerError, &responseErr{
			Error: domain.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Type:    "UserException",
			},
		})
	}

	return c.JSON(http.StatusOK, &responseUsers{Users: users})
}
