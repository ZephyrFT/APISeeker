package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/TheCoderAlex/APISeer/run"
	"github.com/TheCoderAlex/APISeer/validation"
)

var (
	website           string
	host              string
	tmpDir            string
	resultDir         string
	chromeDriverPath  string
	chromePath        string
	cookie            string
	keywordsLLMURL    string
	validationsLLMURL string
	maxCrawledLimit   int
	useOCR            bool
)

func main() {
	flag.StringVar(&website, "url", "https://www.baidu.com", "Website URL")
	flag.StringVar(&host, "host", "baidu.com", "Host name")
	flag.StringVar(&tmpDir, "tmp", ".tmp", "Temporary directory")
	flag.StringVar(&resultDir, "result", "result", "Result directory")
	flag.StringVar(&keywordsLLMURL, "keywords", "http://192.168.3.6:8000/v1/chat/completions", "Keywords LLM URL")
	flag.StringVar(&validationsLLMURL, "validations", "http://192.168.3.5:8000/v1/chat/completions", "Validations LLM URL")
	flag.StringVar(&chromeDriverPath, "chromedriver", "/usr/local/bin/chromedriver", "Chrome driver path")
	flag.StringVar(&chromePath, "chrome", "/usr/bin/google-chrome", "Chrome path")
	flag.StringVar(&cookie, "cookie", "", "Cookie")
	flag.IntVar(&maxCrawledLimit, "limit", 500, "Max crawled limitation")
	flag.BoolVar(&useOCR, "ocr", false, "Use OCR to extract text from images")

	flag.Parse()

	if len(host) == 0 {
		fmt.Println("error: host cannot be empty.")
		flag.Usage()
		os.Exit(1)
	}

	if len(website) == 0 {
		fmt.Println("error: website cannot be empty.")
		flag.Usage()
		os.Exit(1)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if !filepath.IsAbs(tmpDir) {
		tmpDir = filepath.Join(currentDir, tmpDir)
	}

	if !filepath.IsAbs(resultDir) {
		resultDir = filepath.Join(currentDir, resultDir)
	}

	_, err = os.Stat(tmpDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(tmpDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// JSCrawler 爬取
	log.Printf("Start crawling %v", website)
	run.CrawlerRun(website, host, maxCrawledLimit, tmpDir, cookie)

	// 使用selenium爬取第二遍正文
	if useOCR {
		log.Printf("Start crawling %v again with OCR", website)
		run.RunOCR(host, tmpDir, chromePath, chromeDriverPath)
	}

	// 关键词提取
	log.Printf("Start extracting keywords from %v", website)
	run.ExtractKeywords(host, tmpDir, keywordsLLMURL)

	// 大模型推理 + validation 验证
	log.Printf("Start validating %v", website)
	validation.RunValidation(host, website, validationsLLMURL, tmpDir, 10, resultDir)
}
