package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Vantuan1606/app-test/domain"
	"github.com/Vantuan1606/app-test/hashtag"
)

type hashtagHTTPHandler struct {
	hashtagUsecase domain.IHashtagUsecase
}

type responseErr struct {
	Error domain.Error `json:"error"`
}

type responseHashtag struct {
	Hashtag interface{} `json:"hashtag"`
}

type responseHashtags struct {
	Hashtags interface{} `json:"hashtags"`
}

func NewHashtagHTTPHandler(e *echo.Echo, hg domain.IHashtagUsecase) {
	handler := &hashtagHTTPHandler{
		hashtagUsecase: hg,
	}

	e.GET("/hashtag", handler.Lists)

}

func (hg *hashtagHTTPHandler) Lists(c echo.Context) error {
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

	input := hashtag.ListHashtagInput{}

	if sort != "" {
		input.SetSort(sort)
	}
	input.SetAscending(ascending)
	input.SetLimit(limitInt)
	input.SetOffset(offetInt)

	hashtags, err := hg.hashtagUsecase.List(ctx, &input)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, &responseErr{
				Error: domain.Error{
					Code:    http.StatusNotFound,
					Message: err.Error(),
					Type:    "HashtagException",
				},
			})
		}

		return c.JSON(http.StatusInternalServerError, &responseErr{
			Error: domain.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Type:    "HashtagException",
			},
		})
	}

	return c.JSON(http.StatusOK, &responseHashtags{Hashtags: hashtags})
}
