package crawler

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/TheCoderAlex/APISeer/config"
	"github.com/TheCoderAlex/APISeer/util"
	"golang.org/x/net/html"
)

type Spider struct {
	URL        string
	Cookie     string
	UserAgent  string
	Host       string
	ParentHTML *Resp

	UrlList    chan []string
	APIList    chan *sync.Map
	SpiderList chan []*Spider

	WG *sync.WaitGroup
}

type API struct {
	Endpoint   string
	Parameters string
}

func (s *Spider) Do() {
	content, err := s.fetch()
	if err != nil {
		return
	}
	// htmlInfo := trimHtml("URL: " + content[5] + " " + "information: " + html.UnescapeString(content[0]))
	htmlInfo, err := contentExtract(content[0], content[5])
	if err != nil {
		return
	}
	if s.ParentHTML != nil {
		htmlInfo = s.ParentHTML
	}
	// 包含的URL
	urls := urlFind(content[0], content[1], content[2], content[3])

	// 包含的JS
	jss := jsFind(content[0], content[1], content[2], content[3])

	// API特征
	re := regexp.MustCompile(`(https?:)?\/\/[\w\.-]+(\/[\w\.-]+\/[\w\./?=&%-]*)|(/[\w\.-]+\/[\w\./?=&%-]*)`)

	parse := url.Parse

	var urlList []string
	apiList := &sync.Map{}
	spiderList := []*Spider{}

	var f bool

	//去个重
	for _, url := range urls {
		// 只爬取单独的网站
		val, err := util.CompareHost(s.Host, url)
		if err != nil || !val {
			continue
		}

		isCrawled := make(chan bool)
		urlcopy := url // 当心闭包
		mapRequests <- func(m map[string]bool) {
			isCrawled <- m[urlcopy]
		}
		vis := <-isCrawled
		// 重复了就不加入结果列表
		if vis {
			continue
		}
		// 打上标记
		mapRequests <- func(m map[string]bool) {
			m[urlcopy] = true
		}
		//“/static/css”之类的前端端点
		for _, p := range config.FrontendPoint {
			if strings.Contains(url, p) {
				f = true
				break
			}
		}

		// 跳过危险接口
		for _, r := range config.Risks {
			if strings.Contains(url, r) {
				f = true
				break
			}
		}

		if !f {
			urlList = append(urlList, url)
			//无需把url传到main中，只需要根据url在这里创建spider，再发送到main中进行爬取
			spiderList = append(spiderList, &Spider{
				URL:        url,
				Cookie:     s.Cookie,
				UserAgent:  s.UserAgent,
				Host:       s.Host,
				ParentHTML: nil,
				UrlList:    s.UrlList,
				APIList:    s.APIList,
				SpiderList: s.SpiderList,
				WG:         s.WG,
			})
		}

		// 如果是根目录，不放入结果
		u, err := parse(url)
		if err != nil {
			continue
		}
		if u.Path == "/" {
			continue
		}
		// 去除前端页面
		if re.MatchString(url) {
			f := false
			for _, suf := range config.FrontendExtensions {
				if strings.HasSuffix(url, suf) {
					f = true
					break
				}
			}
			if !f {
				// apiList = append(apiList, url)
				apiList.Store(url, htmlInfo)
			}
		}

	}

	f = false
	//发现的JS也会被爬虫爬，会在被爬到的时候按照上面的逻辑发现API
	for _, js := range jss {
		val, err := util.CompareHost(s.Host, js)
		if err != nil || !val {
			continue
		}

		isCrawled := make(chan bool)
		jscopy := js
		mapRequests <- func(m map[string]bool) {
			isCrawled <- m[jscopy]
		}
		vis := <-isCrawled
		if vis {
			continue
		}
		mapRequests <- func(m map[string]bool) {
			m[jscopy] = true
		}

		//“/static/css”之类的前端端点
		for _, p := range config.FrontendPoint {
			if strings.Contains(js, p) {
				f = true
				break
			}
		}

		// 跳过危险接口
		for _, r := range config.Risks {
			if strings.Contains(js, r) {
				f = true
				break
			}
		}

		if !f {
			urlList = append(urlList, js)
			spiderList = append(spiderList, &Spider{
				URL:        js,
				Cookie:     s.Cookie,
				UserAgent:  s.UserAgent,
				Host:       s.Host,
				ParentHTML: htmlInfo,
				UrlList:    s.UrlList,
				APIList:    s.APIList,
				SpiderList: s.SpiderList,
				WG:         s.WG,
			})
		}

		// 如果是根目录，不放入结果
		u, err := parse(js)
		if err != nil {
			continue
		}
		if u.Path == "/" {
			continue
		}
		// 去除前端页面
		if re.MatchString(js) {
			f := false
			for _, suf := range config.FrontendExtensions {
				if strings.HasSuffix(js, suf) {
					f = true
					break
				}
			}
			if !f {
				// apiList = append(apiList, js)
				apiList.Store(js, htmlInfo)
			}
		}
	}

	//返回给主线程，同时添加到WaitGroup
	// s.UrlList <- urlList
	s.APIList <- apiList
	s.SpiderList <- spiderList
	s.WG.Add(len(urlList))
}

var client *http.Client

// 发送操作map的请求
var mapRequests = make(chan func(map[string]bool), 10)

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Second * 30,
			KeepAlive: time.Second * 30,
		}).DialContext,
		MaxIdleConns:          2 * 30,
		MaxIdleConnsPerHost:   30 + 10,
		IdleConnTimeout:       time.Second * 90,
		TLSHandshakeTimeout:   time.Second * 90,
		ExpectContinueTimeout: time.Second * 10,
	}

	//TODO: 设置代理
	// tr.Proxy = http.ProxyURL(&url.URL{Host: "127.0.0.1:7890"})
	client = &http.Client{
		Timeout:   time.Duration(5) * time.Second,
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			// 是否需要检测重定向？
			return nil
		},
	}

	// map去重协程
	go func() {
		m := make(map[string]bool)
		for req := range mapRequests {
			req(m)
		}
	}()
}

func (s *Spider) fetch() ([]string, error) {
	link, _ := url.QueryUnescape(s.URL)

	request, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept-Encoding", "gzip")
	request.Header.Set("User-Agent", util.GetUserAgent(s.UserAgent))
	request.Header.Set("Accept", "*/*")

	if s.Cookie != "" {
		request.Header.Set("Cookie", s.Cookie)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 404 {
		return nil, fmt.Errorf("page not found: %s", link)
	}

	defer response.Body.Close()

	var result string

	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		con, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		result = string(con)
	} else {
		dataBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		result = string(dataBytes)
	}

	// contentType := response.Header.Get("Content-Type")
	// if contentType == "" {
	// 	return nil, fmt.Errorf("NO CONTENT-TYPE: %s", link)
	// }
	// // 判断 Content-Type 是否包含 "application/json"
	// if !strings.Contains(contentType, "application/json") {
	// 	return nil, fmt.Errorf("NO JSON: %s", link)
	// }

	// 针对seu的pass
	if strings.Contains(result, "访问地址无效") {
		return nil, fmt.Errorf("page not found: %s", link)
	}

	//针对人类认证
	if strings.Contains(result, "you need to verify that you're not a robot") {
		return nil, fmt.Errorf("captcha found: %s", link)
	}

	//error page
	if strings.HasPrefix(result, "error") {
		return nil, fmt.Errorf("error page: %s", link)
	}

	//去除文件
	for _, fileSignature := range config.FileSignatures {
		if strings.HasPrefix(result, fileSignature) {
			return nil, fmt.Errorf("file found: %s", link)
		}
	}

	// 为相对地址做准备
	path := response.Request.URL.Path
	host := response.Request.URL.Host
	scheme := response.Request.URL.Scheme
	source := scheme + "://" + host + path

	// 如果有base标签，更改地址
	re := regexp.MustCompile("base.{1,5}href.{1,5}(http.+?//[^\\s]+?)[\"'‘“]")
	base := re.FindAllStringSubmatch(result, -1)
	if len(base) > 0 {
		host = regexp.MustCompile("http.*?//([^/]+)").FindAllStringSubmatch(base[0][1], -1)[0][1]
		scheme = regexp.MustCompile("(http.*?)://").FindAllStringSubmatch(base[0][1], -1)[0][1]
		paths := regexp.MustCompile("http.*?//.*?(/.*)").FindAllStringSubmatch(base[0][1], -1)
		if len(paths) > 0 {
			path = paths[0][1]
		} else {
			path = "/"
		}
	}

	return []string{result, host, scheme, path, source, link}, nil
}

func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile(`\\<[\\S\\s]+?\\>`)
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile(`\\<style[\\S\\s]+?\\</style\\>`)
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile(`\\<script[\\S\\s]+?\\</script\\>`)
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile(`\\<[\\S\\s]+?\\>`)
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile(`\\s{2,}`)
	src = re.ReplaceAllString(src, "\n")

	// 替换所有可能的换行符
	src = strings.ReplaceAll(src, "\r\n", " ")
	src = strings.ReplaceAll(src, "\n", " ")
	src = strings.ReplaceAll(src, "\r", " ")

	// 替换 HTML 中的换行和空格
	src = strings.ReplaceAll(src, "<br>", " ")
	src = strings.ReplaceAll(src, "&nbsp;", " ")

	// 正则表达式（可选，用于更复杂的情况）
	ree := regexp.MustCompile(`\r?\n`)
	src = ree.ReplaceAllString(src, " ")

	return strings.TrimSpace(src)
}

// 正文提取
var cleanRe = regexp.MustCompile(`\s{2,}`)

type Resp struct {
	URL     string
	Content string
	APIS    []API
}

func RemoveHTMLTags(input string) (string, error) {
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return "", err
	}

	var result strings.Builder
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
			return
		}

		if n.Type == html.CommentNode {
			return
		}

		if n.Type == html.TextNode {
			result.WriteString(html.UnescapeString(n.Data))
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}

	extractText(doc)

	return result.String(), nil
}

func contentExtract(input string, u string) (*Resp, error) {
	rmHTML, err := RemoveHTMLTags(input)
	if err != nil {
		return nil, err
	}

	// 简单清洗：去除首尾和多余空格
	content := rmHTML
	content = strings.TrimSpace(content)
	content = cleanRe.ReplaceAllString(content, " ")

	response := Resp{
		URL:     u,
		Content: content,
		APIS:    []API{},
	}

	return &response, nil
}
