package main

import (
	"apifinder/run"
	"apifinder/validation"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	website           string
	host              string
	tmpDir            string
	keywordsLLMURL    string
	validationsLLMURL string
	maxCrawledLimit   int
	useOCR            bool
)

func main() {
	flag.StringVar(&website, "url", "https://www.seu.edu.cn", "Website URL")
	flag.StringVar(&host, "host", "seu.edu.cn", "Host name")
	flag.StringVar(&tmpDir, "tmp", ".tmp", "Temporary directory")
	flag.StringVar(&keywordsLLMURL, "keywords", "http://121.248.52.64:8000/v1/chat/completions", "Keywords LLM URL")
	flag.StringVar(&validationsLLMURL, "validations", "http://121.248.52.64:8001/v1/chat/completions", "Validations LLM URL")
	flag.IntVar(&maxCrawledLimit, "limit", 1000, "Max crawled limitation")
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

	_, err = os.Stat(tmpDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(tmpDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// JSCrawler 爬取
	log.Printf("Start crawling %v", website)
	run.CrawlerRun(website, host, maxCrawledLimit, tmpDir)

	// 使用selenium爬取第二遍正文
	if useOCR {
		log.Printf("Start crawling %v again with OCR", website)
		run.RunOCR(host, tmpDir)
	}

	// 关键词提取
	log.Printf("Start extracting keywords from %v", website)
	run.ExtractKeywords(host, tmpDir, keywordsLLMURL)

	// 大模型推理 + validation 验证
	log.Printf("Start validating %v", website)
	validation.RunValidation(website, validationsLLMURL, tmpDir, 10)
}
