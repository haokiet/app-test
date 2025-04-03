package main

import (
	"encoding/json"
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

type RequestData struct {
	Link   string `json:"link"`
	Number int    `json:"number"`
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data RequestData

		// Giải mã dữ liệu JSON từ request
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Lỗi đọc dữ liệu", http.StatusBadRequest)
			return
		}

		// In dữ liệu ra console
		fmt.Printf("Nhận được dữ liệu: Link = %s, Number = %d\n", data.Link, data.Number)

		// Phản hồi lại client
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Dữ liệu nhận được: Link = %s, Number = %d", data.Link, data.Number)))
	} else {
		http.Error(w, "Chỉ hỗ trợ POST", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("."))) // Phục vụ file tĩnh (index.html)
	http.HandleFunc("/user", apiHandler)
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
