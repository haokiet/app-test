package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Vantuan1606/app-buff/config"
	"github.com/Vantuan1606/app-buff/service/database"
	userHTTPHandler "github.com/Vantuan1606/app-buff/user/delivery/http"
	userRepo "github.com/Vantuan1606/app-buff/user/repo"
	userUsecase "github.com/Vantuan1606/app-buff/user/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	config := config.GetConfig()

	logrus.SetReportCaller(true)

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet},
	}))

	database.NewMongoService()

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	pr := userRepo.NewMongoUserRepo()

	pu := userUsecase.NewUserUsecase(pr, 60*time.Second)
	userHTTPHandler.NewUserHTTPHandler(e, pu)

	privateAPI := e.Group("/private")

	privateAPI.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", config.App.Port)))
}
