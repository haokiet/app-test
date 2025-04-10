package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"github.com/Vantuan1606/app-test/config"
	hashtagHTTPHandler "github.com/Vantuan1606/app-test/hashtag/delivery/http"
	hashtagRepo "github.com/Vantuan1606/app-test/hashtag/repo"
	hashtagUsecase "github.com/Vantuan1606/app-test/hashtag/usecase"
	"github.com/Vantuan1606/app-test/service/database"
	userHTTPHandler "github.com/Vantuan1606/app-test/user/delivery/http"
	userRepo "github.com/Vantuan1606/app-test/user/repo"
	userUsecase "github.com/Vantuan1606/app-test/user/usecase"
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

	http.Handle("/", http.FileServer(http.Dir(".")))
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
	hg := hashtagRepo.NewMongoHashtagRepo()

	pu := userUsecase.NewUserUsecase(pr, 60*time.Second)
	hgs := hashtagUsecase.NewHashtagUsecase(hg, 60*time.Second)
	userHTTPHandler.NewUserHTTPHandler(e, pu)
	hashtagHTTPHandler.NewHashtagHTTPHandler(e, hgs)

	privateAPI := e.Group("/private")

	privateAPI.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", config.App.Port)))
}
