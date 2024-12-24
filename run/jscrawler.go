package run

import (
	"apifinder/crawler"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
)

type data struct {
	URL        string
	Content    string
	VisibleAPI []crawler.API
}

func processAPI(API string, apiChan chan string, paramChan chan string) {
	//依次处理api的情况
	/*
		1. 对于普通URL：直接去除域名后保留
		2. 带GET参数的：用url提取参数
		3. 路径中带参数的："\\/\\d+\\/", "\\/[0-9a-fA-F_]+\\/",
		4. 带后缀的：去除后缀后当成普通URL使用
	*/
	var newAPI string
	var params string

	//先把URL分成协议+域名、端点、GET参数三个部分
	u, err := url.Parse(API)
	if err != nil {
		logrus.Error("Error parsing URL:", err)
		return
	}

	//处理GET参数
	p := u.Query()
	if p != nil {
		params = p.Encode()
	}

	//端点
	//首次去除后缀包括.jsp 和 #home-t这种
	if u.Path != "" && strings.Contains(u.Path, ".") {
		newAPI = strings.Split(u.Path, ".")[0] + "/"
	} else if u.Path != "" {
		newAPI = u.Path + "/"
	}
	//前缀为“//”的情况
	if u.Path != "" && strings.HasPrefix(u.Path, "//") {
		newAPI = strings.Replace(u.Path, "//", "/", -1)
	}

	//判断是否有路径中参数，有的话用{info}占位，否则直接保留
	//只考虑path中： 存在十六进制、十进制数字则认为是路径参数
	// for _, re := range config.ParamsInPath {
	// 	reg := regexp.MustCompile(re)
	// 	for {
	// 		pip := reg.FindAllStringSubmatch(newAPI, -1)
	// 		if len(pip) == 0 {
	// 			break
	// 		}
	// 		for p := range pip {
	// 			//对每一个匹配上的字符，都进行格式替换，并增加param
	// 			newS := strings.TrimLeft(pip[p][0], "/")
	// 			// fmt.Println("====================================", pip[p])
	// 			newAPI = strings.ReplaceAll(newAPI, newS, "{info}/")
	// 			// fmt.Println("************************************", newAPI)
	// 			if len(params) != 0 {
	// 				params = params + "&info=" + newS
	// 			} else {
	// 				params = "info=" + newS
	// 			}
	// 			params = strings.TrimRight(params, "/")
	// 			// fmt.Println(params)

	// 		}
	// 	}
	// }

	if newAPI == "" {
		logrus.Error("Error processing API:", API)
		return
	}
	newAPI = strings.TrimRight(newAPI, " ")
	newAPI = strings.TrimRight(newAPI, "/")

	apiChan <- newAPI
	paramChan <- params

}

func CrawlerRun(website string, host string, crawledLimit int, tmpDir string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	APIsReceived := make(chan *sync.Map, 10)
	UrlsReceived := make(chan []string, 10)
	SpiderReceived := make(chan []*crawler.Spider, 10)
	sigs := make(chan os.Signal, 1)
	var newAPIs sync.Map
	htmlMap := &sync.Map{}

	out := make(chan string, 10)
	newAPIChan := make(chan string)
	newParamsChan := make(chan string)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	jsonFileName := path.Join(tmpDir, "APIInfo-"+host+".json")

	// 输出处理
	go func() {
		cnt := 0
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-out:
				cnt++
				log.Printf("The %d pages visited: %s\n", cnt, msg)
				if cnt >= crawledLimit {
					log.Printf("Crawler limit has been reached: %d, ready to return.\n", crawledLimit)
					cancel()
					return
				}
			}
		}
	}()

	// 接收spider
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case spiders := <-SpiderReceived:
				for _, spider := range spiders {
					out <- spider.URL
					go func(sp *crawler.Spider) {
						sp.Do()
						sp.WG.Done()
					}(spider)
				}
			}
		}
	}()

	// API去重+写入csv
	go func() {
		for apis := range APIsReceived {
			//处理API信息，提取参数，存为csv文件
			apis.Range(func(key, value interface{}) bool {
				api := key.(string)
				url := value.(*crawler.Resp)

				go processAPI(api, newAPIChan, newParamsChan)
				//在写之前处理
				newAPI := <-newAPIChan
				params := <-newParamsChan
				if newAPI == "" {
					return false
				}

				// 使用 sync.Map 的 LoadOrStore 方法
				if _, loaded := newAPIs.LoadOrStore(newAPI, true); !loaded {
					if strings.HasPrefix(newAPI, "/_vsl") {
						return false
					}
					if strings.Contains(newAPI, "/js") {
						return false
					}

					pre, isExist := htmlMap.Load(url)
					if isExist {
						pre = append(pre.([]crawler.API),
							crawler.API{
								Endpoint:   newAPI,
								Parameters: params,
							},
						)
					} else {
						pre = append([]crawler.API{},
							crawler.API{
								Endpoint:   newAPI,
								Parameters: params,
							},
						)
					}

					htmlMap.Store(url, pre)
				}

				return true
			})

		}
	}()

	// 添加第一个任务
	wg.Add(1)
	job := &crawler.Spider{
		URL:        website,
		Cookie:     "",
		UserAgent:  "",
		Host:       host,
		ParentHTML: nil,
		UrlList:    UrlsReceived,
		APIList:    APIsReceived,
		SpiderList: SpiderReceived,
		WG:         &wg,
	}

	go func() {
		job.Do()
		wg.Done()
	}()

	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-ctx.Done():
		saveToJSON(htmlMap, jsonFileName)
		return
	case <-done:
		saveToJSON(htmlMap, jsonFileName)
		cancel()
		return
	case <-sigs:
		saveToJSON(htmlMap, jsonFileName)
		cancel()
		return
	}
}

func saveToJSON(htmlMap *sync.Map, filename string) {
	jsonFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error creating jsonfile:", err)
		return
	}
	defer jsonFile.Close()
	htmlEncoder := json.NewEncoder(jsonFile)

	var resultList []data
	htmlMap.Range(func(key, value interface{}) bool {

		inputList := key.(*crawler.Resp)

		URL := inputList.URL
		Content := inputList.Content
		APIs := value.([]crawler.API)

		res := data{
			URL:        URL,
			Content:    Content,
			VisibleAPI: APIs,
		}

		resultList = append(resultList, res)
		return true
	})

	error := htmlEncoder.Encode(resultList)
	if error != nil {
		fmt.Println("Error writing to jsonfile:", error)
		return
	}
}
