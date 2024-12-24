package validation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

var (
	website      string
	host         string
	numOfWorkers int
)
var client = http.Client{}

// https://platform.openai.com/docs/api-reference/chat/create
func chatCompletions(model string, content string, options map[string]interface{}) ([]byte, error) {
	payload := make(map[string]interface{})
	var messages []map[string]interface{}

	sys := make(map[string]interface{})
	usr := make(map[string]interface{})

	sys["role"] = "system"
	sys["content"] = "You are a helpful assistant with extensive experience in cybersecurity.\n"
	usr["role"] = "user"
	usr["content"] = content
	messages = append(messages, sys, usr)

	payload["model"] = model
	payload["messages"] = messages

	for k, v := range options {
		payload[k] = v
	}

	bytesData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json: %v", err)
	}

	reader := bytes.NewReader(bytesData)

	request, err := http.NewRequest("POST", host, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request %v", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("accept", " application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	defer response.Body.Close()

	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return respBytes, nil
}

func extractMessageFromResponse(input []byte) (string, error) {
	var payload map[string]interface{}
	err := json.Unmarshal(input, &payload)
	if err != nil {
		return "", fmt.Errorf("failed to parse json: %v", err)
	}

	_, isExist := payload["id"]
	if !isExist {
		return "", fmt.Errorf("unknown response: %v", string(input))
	}

	choices := payload["choices"].([]interface{})
	choice := choices[0].(map[string]interface{})

	message, isExist := choice["message"].(map[string]interface{}) // 只取第一个消息
	if !isExist {
		return "", fmt.Errorf("message not in response: %v", string(input))
	}

	content, isExist := message["content"].(string)
	if !isExist {
		return "", fmt.Errorf("content not in response: %v", string(input))
	}

	return content, nil
}

func validation(response string, domain string, result chan string, failed chan string) {
	respBytes := []byte(response)
	var payload []interface{}

	err := json.Unmarshal(respBytes, &payload)
	if err != nil {
		log.Printf("[ERROR] in validation: %v", err)
		return
	}

	for i, item := range payload {
		endpoint := item.(map[string]interface{})["Endpoint"].(string)
		u := domain + endpoint

		resp, err := http.Head(u)
		if err != nil {
			log.Printf("[ERROR] in sending request: %v", err)
			continue
		}

		status := ""
		if resp.StatusCode != 404 {
			status = "ok"
			result <- fmt.Sprintf("URL: %s founded. Status code: %d.\n", u, resp.StatusCode)
		} else {
			status = "pass"
			failed <- fmt.Sprintf("URL: %s failed. Status code: %d.\n", u, resp.StatusCode)
		}

		log.Printf("[INFO] validating %v in %v (%d/%d). Status code %v, %s.", endpoint, domain, i+1, len(payload), resp.StatusCode, status)

		resp.Body.Close()
	}
}

type API struct {
	Endpoint   string
	Parameters string
}

func convertToAPISlice(input interface{}) ([]API, error) {
	if input == nil {
		return nil, fmt.Errorf("input is nil")
	}

	slice, ok := input.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected []interface{}, got %T", input)
	}

	result := make([]API, 0, len(slice))
	for _, item := range slice {
		if apiMap, ok := item.(map[string]interface{}); ok {
			endpoint := ""
			parameters := ""
			if apiMap["Endpoint"] != nil {
				endpoint = apiMap["Endpoint"].(string)
			}

			if apiMap["Parameter"] != nil {
				parameters = apiMap["Parameter"].(string)
			}

			api := API{
				Endpoint:   endpoint,
				Parameters: parameters,
			}
			result = append(result, api)
		}
	}
	return result, nil
}

func RunValidation(hst string, site string, h string, tmpDir string, workers int, resultDir string) {
	website = site
	host = h
	numOfWorkers = workers

	queue := make(chan bool, numOfWorkers)
	result := make(chan string, numOfWorkers)
	failed := make(chan string, numOfWorkers)

	currentTime := time.Now()
	fileName := "Result_" + hst + "_" + currentTime.Format("20060102_150405")
	fileNameSucceed := fileName + ".txt"
	fileNameFailed := fileName + "_failed.txt"

	// 不存在则创建
	if _, err := os.Stat(resultDir); os.IsNotExist(err) {
		err = os.MkdirAll(resultDir, 0755)
		if err != nil {
			log.Fatal("can't create the result directory: ", err)
		}
	}

	fileSucceed, err := os.OpenFile(path.Join(resultDir, fileNameSucceed), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("can't open the success result file: ", err)
	}
	defer fileSucceed.Close()

	fileFailed, err := os.OpenFile(path.Join(resultDir, fileNameFailed), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("can't open the success result file: ", err)
	}
	defer fileFailed.Close()

	inputTemplate := "Given the following website information keywords and some visible apis in the website, please respond with a JSON for any hidden API endpoints that may exist in this website.\n\nResponse Format:\n{\n  \"Shadow_APIs\": [\n    {\n      \"Endpoint\": \"API endpoint\",\n      \"Parameters\": [\"parameter1\", \"parameter2\", ...] or \"\"\n    }\n  ]\n}\n\nGiven the keywords \"RoleAdd adds role name authentication system\" and the visible API \"/auth/user/login\", we can deduce a potential hidden API endpoint. First, \"RoleAdd\" indicates adding a role, \"role name\" is the resource, and \"authentication system\" provides the context. Since the visible API \"/auth/user/login\" handles user authentication, it’s likely that other user-related actions follow a similar structure. Therefore, we can infer that an API for role management might be under the same \"/auth\" path, using a similar naming convention, such as \"/auth/role/add\" or \"/auth/role/create\". This reasoning, based on the naming patterns and structure of the visible API, leads to the conclusion that a hidden endpoint could be \"/auth/role/add\", responsible for adding roles within the authentication system.\n\nFocus more on specific, hidden, or less obvious endpoints based on provided keywords. DO NOT ANSWER THE GIVEN VISIBLE APIS! THE JSON FORMAT MUST BE CORRECT!\n\nkeywords: "
	filename := path.Join(tmpDir, "extracted-"+hst+".json")

	jsonByte, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("failed to read json file: ", err)
	}

	var payload []map[string]interface{}
	err = json.Unmarshal(jsonByte, &payload)
	if err != nil {
		log.Fatal("failed to parse json file: ", err)
	}

	// 成功结果处理
	go func() {
		for res := range result {
			_, err = fileSucceed.WriteString(res)
			if err != nil {
				log.Printf("[ERROR] in write succeed result: %v", err)
			}
		}
	}()

	// 失败结果处理
	go func() {
		for res := range failed {
			_, err = fileFailed.WriteString(res)
			if err != nil {
				log.Printf("[ERROR] in write failed result: %v", err)
			}
		}

	}()

	for _, element := range payload {
		content := element["Content"].(string)
		u := element["URL"].(string)

		visibleAPI, err := convertToAPISlice(element["VisibleAPI"])
		if err != nil {
			log.Printf("[ERROR] converting VisibleAPI: %v", err)
			continue
		}

		parsedURL, err := url.Parse(u)
		if err != nil {
			log.Printf("[ERROR] invalid url format: %v", u)
			continue
		}

		domain := parsedURL.Scheme + "://" + parsedURL.Host

		content = inputTemplate + content + "\n\nVisible APIs: "

		for _, api := range visibleAPI {
			content += api.Endpoint + " "
		}

		queue <- true
		go func(content string, domain string) {
			defer func() { <-queue }()

			respBytes, err := chatCompletions("Meta-Llama-3.1-8B-Instruct", content, nil)
			if err != nil {
				log.Printf("[ERROR] in receive response: %v", err)
				return
			}

			message, err := extractMessageFromResponse(respBytes)
			if err != nil {
				log.Printf("[ERROR] in extract message: %v", err)
				return
			}

			go validation(message, domain, result, failed)
		}(content, domain)
	}
	close(result)
	close(failed)
}
