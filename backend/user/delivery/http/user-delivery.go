package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/Vantuan1606/app-test/domain"
	"github.com/Vantuan1606/app-test/user"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"

	// "github.com/chromedp/chromedp/kb"
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
	// Lấy proxy ngẫu nhiên từ API
	proxy, erro := getRandomProxy("3uisbtk31gnbde9bhzd0kbwbqpe0ftph4hajrsxt")
	if erro != nil {
		fmt.Printf("Không thể lấy proxy: %v\n", erro)
	}
	// Gọi logic xem video
	var wg sync.WaitGroup
	wg.Add(1)
	go uss.watchVideo(account, req.VideoURL, &wg, proxy)
	wg.Wait()

	// Phản hồi thành công
	return c.JSON(http.StatusOK, &responseUser{
		User: "Video watched successfully",
	})
}

// ProxyResponse cấu trúc dữ liệu trả về từ API Webshare
type ProxyResponse struct {
	Count   int `json:"count"`
	Results []struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		ProxyAddress string `json:"proxy_address"`
		Ports        struct {
			HTTP   int `json:"http"`
			Socks5 int `json:"socks5"`
		} `json:"ports"`
		CountryCode string `json:"country_code"`
		CityName    string `json:"city_name"`
	} `json:"results"`
}

// Hàm lấy proxy ngẫu nhiên từ Webshare API
// Hàm lấy proxy ngẫu nhiên từ Webshare API
func getRandomProxy(apiKey string) (string, error) {
	// URL của API Webshare để lấy danh sách proxy
	apiURL := "https://proxy.webshare.io/api/proxy/list/"

	// Gửi yêu cầu GET với API Key
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Token "+apiKey)

	// Gửi yêu cầu HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Đọc phản hồi body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Kiểm tra mã trạng thái
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("lỗi khi lấy proxy từ Webshare: %s", resp.Status)
	}

	// Giải mã phản hồi JSON vào struct ProxyResponse
	var proxyResp ProxyResponse
	if err := json.Unmarshal(body, &proxyResp); err != nil {
		return "", err
	}

	// Kiểm tra nếu không có proxy nào trong kết quả
	if len(proxyResp.Results) == 0 {
		return "", fmt.Errorf("không có proxy nào trong kết quả")
	}

	rand.Seed(time.Now().UnixNano()) // Đảm bảo mỗi lần chạy sẽ chọn ngẫu nhiên
	randomIndex := rand.Intn(len(proxyResp.Results))

	// Lấy địa chỉ proxy và cổng
	proxyAddress := proxyResp.Results[randomIndex].ProxyAddress
	httpPort := proxyResp.Results[randomIndex].Ports.HTTP // Cổng HTTP
	// socksPort := proxyResp.Results[randomIndex].Ports.Socks5 // Cổng SOCKS5

	// Bạn có thể thay đổi nếu muốn dùng cổng SOCKS5
	proxyWithPort := fmt.Sprintf("%s:%s@%s:%d", proxyAddress, httpPort) // Sử dụng cổng HTTP

	// Trả về proxy đầy đủ (địa chỉ + cổng)
	return proxyWithPort, nil
}

func (uss *UserHTTPHandler) watchVideo(account *domain.User, videoURL string, wg *sync.WaitGroup, proxy string) {
	if wg != nil {
		defer wg.Done()
	}

	loginURL := "https://www.tiktok.com/login/phone-or-email/email"
	// account.Username = "pikakun53"
	// account.Password = "Kiet2001!"
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

	// Thêm cấu hình proxy nếu được cung cấp
	if proxy != "" {
		opts = append(opts, chromedp.ProxyServer(proxy))
	}
	// 2. Thêm profile người dùng thật
	userDataDir := filepath.Join(os.TempDir(), "new_chrome_profile_"+strconv.Itoa(rand.Intn(10000)))
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
		fmt.Printf("Lỗi thiết lập")
		return
	}

	/// Thêm hành vi người dùng tự nhiên
	err = chromedp.Run(ctx,
		chromedp.Navigate(loginURL),
		chromedp.Sleep(randomDuration(2, 3)),
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
		return
	}

	//Kiểm tra CAPTCHA
	var hasCaptcha bool
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`
	        // Kiểm tra cả 2 loại CAPTCHA phổ biến của TikTok
	        document.querySelector('.captcha-verify-container') !== null ||
	        document.querySelector('iframe[src*="captcha"]') !== null ||
	        document.querySelector('div[id*="verify"]') !== null
	    `, &hasCaptcha),
	)
	if err != nil {
		return
	}

	if hasCaptcha {
		time.Sleep(20 * time.Second)
	}

	// Kiểm tra đăng nhập thành công
	var cookies []*network.Cookie
	err = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		cookies, err = network.GetCookies().Do(ctx)
		return err
	}))
	var loginCheck bool

	// Tìm giá trị sessionid
	for _, cookie := range cookies {
		if cookie.Name == "sessionid" {
			loginCheck = true
		}
	}

	if err != nil || loginCheck == false {

		return
	}

	// Xem video với hành vi tự nhiên
	err = chromedp.Run(ctx,
		chromedp.Navigate(videoURL),
		chromedp.Sleep(randomDuration(3, 5)), // Chờ trang tải

		// Bắt đầu gửi bình luận tự động
		chromedp.ActionFunc(func(ctx context.Context) error {
			comments := []string{
				"Video hay quá!",
				"Tôi rất thích nội dung này",
				"Cảm ơn bạn đã chia sẻ",
				"❤️❤️❤️",
				"Quá tuyệt vời!",
				"Tôi sẽ chia sẻ video này",
				"Nội dung chất lượng",
				"Bạn thật tài năng",
			}

			// Chọn ngẫu nhiên một bình luận từ danh sách
			randomIndex := rand.Intn(len(comments))
			randomComment := comments[randomIndex]

			// Tìm ô nhập bình luận và gửi
			err := chromedp.Run(ctx,
				chromedp.WaitVisible(`.tiktok-1772j3i[contenteditable="plaintext-only"]`, chromedp.ByQuery),
				chromedp.Click(`.tiktok-1772j3i[contenteditable="plaintext-only"]`, chromedp.ByQuery),
				chromedp.Sleep(2*time.Second),
				chromedp.SendKeys(`.tiktok-1772j3i[contenteditable="plaintext-only"]`, randomComment, chromedp.ByQuery),
				chromedp.Sleep(1*time.Second),
				chromedp.Click(`.tiktok-mortok.e2lzvyu9`, chromedp.ByQuery),        // Nút gửi
				chromedp.WaitVisible(`.tiktok-fa6jvh.e1tv929b2`, chromedp.ByQuery), // Đảm bảo bình luận đã được gửi (có thể cần điều chỉnh selector)
				chromedp.Click(`.tiktok-fa6jvh.e1tv929b2`, chromedp.ByQuery),       // Có thể cần thêm một hành động sau khi gửi
				chromedp.Sleep(5*time.Second),
			)

			if err != nil {
				return err
			}
			return nil
		}),
	)

	if err != nil {
		return
	}

	time.Sleep(30 * time.Minute)

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
