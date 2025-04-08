package http

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Vantuan1606/app-test/domain"
	"github.com/Vantuan1606/app-test/user"
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

	e.GET("/user", handler.Lists)
	e.POST("watch", handler.WatchVideoWithAccount)

}

func (uss *UserHTTPHandler) Lists(c echo.Context) error {
	limit := c.QueryParam("limit")
	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil || limitInt <= 0 {
		return c.JSON(http.StatusBadRequest, &responseErr{
			Error: domain.Error{
				Code:    http.StatusBadRequest,
				Message: "Invalid limit value",
				Type:    "UserException",
			},
		})
	}

	ctx := c.Request().Context()

	input := user.ListUserInput{}
	input.SetLimit(limitInt)

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

	// Trả về danh sách user
	return c.JSON(http.StatusOK, &responseUsers{Users: users})
}

func (uss *UserHTTPHandler) WatchVideoWithAccount(c echo.Context) error {
	// Lấy dữ liệu từ body của yêu cầu
	req := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		VideoURL string `json:"video_url"`
	}{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &responseErr{
			Error: domain.Error{
				Code:    http.StatusBadRequest,
				Message: "Invalid request body",
				Type:    "UserException",
			},
		})
	}

	// Kiểm tra URL video
	if _, err := url.ParseRequestURI(req.VideoURL); err != nil {
		return c.JSON(http.StatusBadRequest, &responseErr{
			Error: domain.Error{
				Code:    http.StatusBadRequest,
				Message: "Invalid video URL",
				Type:    "UserException",
			},
		})
	}

	// Tạo tài khoản người dùng từ thông tin yêu cầu
	account := &domain.User{
		Username: req.Username,
		Password: req.Password,
	}

	// Gọi logic xem video
	var wg sync.WaitGroup
	wg.Add(1)
	go uss.watchVideo(account, req.VideoURL, &wg)
	wg.Wait()

	// Phản hồi thành công
	return c.JSON(http.StatusOK, &responseUser{
		User: "Video watched successfully",
	})
}

func (uss *UserHTTPHandler) watchVideo(account *domain.User, videoURL string, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	loginURL := "https://www.tiktok.com/login/phone-or-email/email"
	account.Username = "pikakun53"
	account.Password = "Kiet2001!"
	log.Println("🔄 videoUrl:", videoURL)
	log.Println("🔄 Starting WatchVideoWithAccount for account:", account.Username)
	defer log.Println("✅ Finished WatchVideoWithAccount for account:", account.Username)

	// ctx, cancel := chromedp.NewContext(context.Background())
	// defer cancel()
	// 1. Cấu hình trình duyệt nâng cao
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("start-maximized", true),
		chromedp.UserAgent(getRandomUserAgent()),
		chromedp.Flag("window-size", "1920,1080"),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("disable-notifications", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("profile-directory", "Default"),
		chromedp.Flag("remote-debugging-port", "9222"),
	)

	// 2. Thêm profile người dùng thật
	userDataDir := filepath.Join(os.TempDir(), "chrome_profile_"+strconv.Itoa(rand.Intn(10000)))
	opts = append(opts, chromedp.UserDataDir(userDataDir))

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// 3. Thiết lập headers và fingerprint
	err := chromedp.Run(ctx,
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers{
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
			"Accept-Language": "en-US,en;q=0.9",
			"Referer":         "https://www.tiktok.com/",
			"DNT":             "1",
			"Connection":      "keep-alive",
		}),
		emulation.SetUserAgentOverride(getRandomUserAgent()),
	)
	if err != nil {
		log.Println("❌ Lỗi khi thiết lập trình duyệt:", err)
		return
	}

	/// Thêm hành vi người dùng tự nhiên
	err = chromedp.Run(ctx,
		chromedp.Navigate(loginURL),
		chromedp.Sleep(randomDuration(2, 4)),

		chromedp.WaitVisible(`input[name="username"]`, chromedp.ByQuery),
		chromedp.Click(`input[name="username"]`),
		chromedp.Sleep(randomDuration(1, 2)),
		sendKeysWithDelay(`input[name="username"]`, account.Username, 200*time.Millisecond),
		chromedp.Sleep(randomDuration(1, 3)),
		chromedp.SendKeys(`input[name="username"]`, kb.Tab),

		chromedp.WaitVisible(`input[placeholder="Password"]`, chromedp.ByQuery),
		chromedp.Click(`input[placeholder="Password"]`),
		chromedp.Sleep(randomDuration(1, 2)),
		sendKeysWithDelay(`input[placeholder="Password"]`, account.Password, 200*time.Millisecond),
		chromedp.Sleep(randomDuration(1, 3)),

		chromedp.Click(`button[type="submit"]`),
		chromedp.Sleep(randomDuration(5, 8)),
	)

	if err != nil {
		log.Println("❌ Lỗi khi đăng nhập:", err)
		return
	}

	// Kiểm tra CAPTCHA
	var hasCaptcha bool
	err = chromedp.Run(ctx, chromedp.Evaluate(`document.querySelector('.captcha-verify-container') !== null`, &hasCaptcha))
	if err != nil || hasCaptcha {
		log.Println("⚠️ CAPTCHA detected. Unable to proceed.")
		return
	}

	// Xem video
	log.Println("✅ Login successful. Watching video...")
	err = chromedp.Run(ctx,
		chromedp.Navigate(videoURL),
		chromedp.Sleep(randomDuration(3, 5)),
	)
	if err != nil {
		log.Println("❌ Error watching video:", err)
		return
	}

	log.Println("✅ Video watched successfully.")
}

// Hàm hỗ trợ - lấy user-agent ngẫu nhiên
func getRandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
	}
	return userAgents[rand.Intn(len(userAgents))]
}

func randomDuration(min, max int) time.Duration {
	return time.Duration(rand.Intn(max-min)+min) * time.Second
}

func sendKeysWithDelay(sel, text string, delay time.Duration) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		for _, c := range text {
			if err := chromedp.SendKeys(sel, string(c), chromedp.ByQuery).Do(ctx); err != nil {
				return err
			}
			if err := chromedp.Sleep(delay).Do(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}
