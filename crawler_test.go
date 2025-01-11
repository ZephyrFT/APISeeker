package main

import (
	"log"
	"os"
	"testing"

	"github.com/TheCoderAlex/APISeeker/run"
	"github.com/TheCoderAlex/APISeeker/validation"
)

func BenchmarkCrawlerRun(b *testing.B) {
	website = "http://www.seu.edu.cn"
	host = "seu.edu.cn"
	tmpDir = "/tmp"
	chromePath = "/path/to/chrome"
	chromeDriverPath = "/path/to/chromedriver"
	useOCR = false
	maxCrawledLimit = 1000

	nullFile, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0644)
	defer nullFile.Close()

	log.SetOutput(nullFile)

	b.ReportAllocs()

	b.ResetTimer() // 在测试开始前重置计时器

	b.Run("JSCrawler", func(b *testing.B) {
		// 模拟 JSCrawler 爬取
		for i := 0; i < b.N; i++ {
			log.Printf("Start crawling %v", website)
			run.CrawlerRun(website, host, maxCrawledLimit, tmpDir, cookie)
		}
	})
}

// func BenchmarkRunOCR(b *testing.B) {
// 	if useOCR {
// 		b.Run("SeleniumOCR", func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				log.Printf("Start crawling %v again with OCR", website)
// 				run.RunOCR(host, tmpDir, chromePath, chromeDriverPath)
// 			}
// 		})
// 	}
// }

func BenchmarkExtractKeywords(b *testing.B) {
	website = "http://www.seu.edu.cn"
	host = "seu.edu.cn"
	tmpDir = "/tmp"
	chromePath = "/path/to/chrome"
	chromeDriverPath = "/path/to/chromedriver"
	useOCR = false

	nullFile, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0644)
	defer nullFile.Close()

	log.SetOutput(nullFile)

	b.ReportAllocs()

	b.Run("KeywordExtraction", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			log.Printf("Start extracting keywords from %v", website)
			run.ExtractKeywords(host, tmpDir, keywordsLLMURL)
		}
	})
}

func BenchmarkRunValidation(b *testing.B) {
	website = "http://www.seu.edu.cn"
	host = "seu.edu.cn"
	tmpDir = "/tmp"
	chromePath = "/path/to/chrome"
	chromeDriverPath = "/path/to/chromedriver"
	useOCR = false

	nullFile, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0644)
	defer nullFile.Close()

	log.SetOutput(nullFile)

	b.ReportAllocs()

	b.Run("Validation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			log.Printf("Start validating %v", website)
			validation.RunValidation(host, website, "http://192.168.3.5:8000/v1/chat/completions", ".tmp", 10, "result")
		}
	})
}
