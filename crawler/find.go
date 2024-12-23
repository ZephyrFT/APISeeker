package crawler

import (
	"net/url"
	"regexp"
	"strings"

	"apifinder/config"
)

func urlFind(cont, host, scheme, path string) []string {
	cata := getCata(path)
	host = scheme + "://" + host

	var result []string
	var valid = url.Parse

	for _, re := range config.UrlFind {
		reg := regexp.MustCompile(re)
		urls := reg.FindAllStringSubmatch(cont, -1)
		urls = urlFilter(urls)

		result = append(result, processUrls(urls, host, cata, scheme, valid)...)
	}
	return result
}

func jsFind(cont, host, scheme, path string) []string {
	cata := getCata(path)
	host = scheme + "://" + host

	var result []string
	var valid = url.Parse

	for _, re := range config.JsFind {
		reg := regexp.MustCompile(re)
		jss := reg.FindAllStringSubmatch(cont, -1)
		jss = jsFilter(jss)

		result = append(result, processUrls(jss, host, cata, scheme, valid)...)
	}
	return result
}

// 提取公共逻辑到新函数
func processUrls(urls [][]string, host, cata, scheme string, valid func(string) (*url.URL, error)) []string {
	var result []string

	for _, u := range urls {
		if u[0] == "" {
			continue
		}
		if strings.HasPrefix(u[0], "https:") || strings.HasPrefix(u[0], "http:") { //完整的绝对路径
			if _, err := valid(u[0]); err == nil {
				result = append(result, u[0])
			}
		} else if strings.HasPrefix(u[0], "//") { // 带//但不带协议的绝对路径
			urlz := scheme + ":" + u[0]
			if _, err := valid(urlz); err == nil {
				result = append(result, urlz)
			}
		} else if strings.HasPrefix(u[0], "/") { // 带/的相对路径
			// urlz := host + url[0]
			URL, err := url.Parse(u[0])
			HOST, err2 := url.Parse(host)
			if err != nil || err2 != nil {
				continue
			}
			urlz := HOST.ResolveReference(URL).String()
			if _, err := valid(urlz); err == nil {
				result = append(result, urlz)
			}
		} else { // 包含没有“/”的相对路径
			// urlz := host + cata + u[0]
			// 排除mailto协议
			if strings.HasPrefix(u[0], "mailto:") {
				continue
			}
			// 排除MIME字符串
			if _, ok := config.BuiltinTypesLower[u[0]]; ok {
				continue
			}
			URL, err := url.Parse(cata + u[0])
			HOST, err2 := url.Parse(host)
			if err != nil || err2 != nil {
				continue
			}
			urlz := HOST.ResolveReference(URL).String()
			if _, err := valid(urlz); err == nil {
				result = append(result, urlz)
			}
		}
	}

	return result
}

// 获取cata的辅助函数
func getCata(path string) string {
	care := regexp.MustCompile("/.*/{1}|/")
	catae := care.FindAllString(path, -1)
	if len(catae) == 0 {
		return "/"
	}
	return catae[0]
}
