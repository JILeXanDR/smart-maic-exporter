package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	scrapeDuration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "page_scrape_duration_seconds",
		Help: "Time taken to scrape the page in seconds.",
	})
)

// Configuration variables
var (
	baseURL = getEnv("BASE_URL", "http://192.168.10.55")
	pinCode = getEnv("PIN_CODE", "0000")
	dataURL = baseURL + "/?page=getwdata"
)

func init() {
	CustomRegistry.MustRegister(scrapeDuration)
}

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any

	zap.ReplaceGlobals(logger)

	sugar := logger.Sugar()

	path, _ := launcher.LookPath()
	if path == "" {
		logger.Panic("failed to LookPath for Chrome")
	}

	// Launch browser with headless mode disabled
	browserURL := launcher.New().Bin(path).Headless(true).Devtools(false).MustLaunch()

	sugar.Debugf("browserURL=%s", browserURL)

	browser := rod.New().ControlURL(browserURL).MustConnect()

	defer func() {
		sugar.Debug("Closing a browser...")
		browser.MustClose()
	}()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				zap.L().Error("recovered from panic", zap.Any("error", err))
				fmt.Fprintf(w, "server error: %+v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		scrapePage(browser, pinCode)

		handler := promhttp.HandlerFor(CustomRegistry, promhttp.HandlerOpts{})
		handler.ServeHTTP(w, r) // Serve updated metrics
	})

	sugar.Debug("Starting server on :8000")

	sugar.Fatal(http.ListenAndServe(":8000", nil))
}

func scrapePage(browser *rod.Browser, pincode string) {
	start := time.Now()

	zap.L().Debug("scrapePage")

	// Navigate to the page
	page := browser.Timeout(3 * time.Second).MustPage(baseURL).MustWaitLoad()
	defer page.MustClose()

	zap.S().Debugf("Current page: title=%s, url=%s", page.MustInfo().Title, page.MustInfo().URL)

	// Check if the page title indicates a login
	if page.MustInfo().Title == "Login" || page.MustInfo().Title == "MAIC Login" {
		zap.S().Debug("Do login")

		page.MustElement(".minput").MustInput(pincode)

		page.MustElement(".msbmit").MustClick()

		page = page.MustWaitLoad()
	}

	zap.S().Debugf("Current page: title=%s, url=%s", page.MustInfo().Title, page.MustInfo().URL)

	newPage := browser.MustPage(dataURL).MustWaitLoad()
	defer newPage.MustClose()

	zap.S().Debugf("Current page: title=%s, url=%s", newPage.MustInfo().Title, newPage.MustInfo().URL)

	html := newPage.MustElement("body").MustHTML()

	zap.L().Sugar().Debugf("got html: %s", html)

	var v T
	err := json.NewDecoder(strings.NewReader(extractJSON(html))).Decode(&v)
	if err != nil {
		zap.L().Error("Decoding failed: %+v", zap.String("html", html), zap.Error(err))
		return
	}

	zap.L().Debug("got result", zap.Any("v", v))

	SetMetrics(v)

	// Update scrape duration metric
	scrapeDuration.Set(time.Since(start).Seconds())
}

func extractJSON(html string) string {
	html = strings.ReplaceAll(html, "<body><pre>", "")
	html = strings.ReplaceAll(html, `</pre><div class="json-formatter-container"></div></body>`, "")
	return html
}
