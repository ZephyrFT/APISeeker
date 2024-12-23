package ocr

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"math"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	browserWidth = 1440
	port         = 9515
)

var (
	chromeDriverPath = "/usr/local/bin/chromedriver"
	jpgOptions       = &jpeg.Options{Quality: 100}
)

func TakeScreenShot(url string, host string, tmpDir string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	// create dir named with host
	dir := path.Join(tmpDir, host)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.Mkdir(dir, 0755)
	}

	contenType := resp.Header.Get("Content-Type")
	if !strings.Contains(contenType, "text/html") {
		return "", errors.New("not a HTML page")
	}

	// Start a WebDriver server instance
	var opts []selenium.ServiceOption
	selenium.SetDebug(false)

	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, opts...)
	if err != nil {
		return "", err
	}
	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "Chrome"}
	chromeCaps := chrome.Capabilities{
		Path: "/usr/bin/google-chrome",
		Args: []string{
			"--headless",
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
		},
	}
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		return "", err
	}
	defer wd.Quit()
	if err := wd.Get(url); err != nil {
		return "", err
	}

	_, _ = wd.ExecuteScript("window.scrollTo(0, document.body.scrollHeight)", nil)
	time.Sleep(5 * time.Second)

	// var width, _ = wd.ExecuteScript("return document.body.parentNode.scrollWidth", nil)
	var height, _ = wd.ExecuteScript("return document.body.parentNode.scrollHeight", nil)
	var realHeight = int(height.(float64))
	// var realWidth = int(width.(float64))

	// lazy loading
	const step = 1000
	var rounds = int(math.Ceil(float64(realHeight) / float64(step)))
	for i := 1; i <= rounds; i++ {
		_, _ = wd.ExecuteScript(fmt.Sprintf("window.scrollTo(0, %d)", i*step), nil)
		time.Sleep(300 * time.Millisecond)
	}

	_ = wd.ResizeWindow("", browserWidth, realHeight)
	_, _ = wd.ExecuteScript("window.scrollTo(0, 0)", nil)

	screenshot, _ := wd.Screenshot()

	// named file with MD5 hash
	filename := fmt.Sprintf("%x.jpg", md5.Sum([]byte(url)))

	pngReader := bytes.NewReader(screenshot)
	img, _ := png.Decode(pngReader)
	outputFile, _ := os.Create(path.Join(dir, filename))
	defer outputFile.Close()
	_ = jpeg.Encode(outputFile, img, jpgOptions)

	return path.Join(dir, filename), nil
}
