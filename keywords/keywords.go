package keywords

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model    string                 `json:"model"`
	Messages []Message              `json:"messages"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

var prompt string = `You will extract a moderate number of concise and relevant keywords that reflect the core functionality or description of the website from the following text. The text need to be processed will be wrapped by '<text></text>'. Ensure no keywords are repeated, and separate them with spaces.

Output format: Keyword1 Keyword2 Keyword3

Do not provide any extra commentary, explanation, or text. Only output the keywords. Do not duplicate keywords.
<text> %s </text>`

func ChatCompletions(url, model, content string, options map[string]interface{}) (string, error) {
	if options == nil {
		options = make(map[string]interface{})
	}

	systemMessage := Message{
		Role:    "system",
		Content: "You are a keyword extraction model that can process natural language text. Your task is to identify and extract key terms or phrases that are central to the content of the text. You will then output these keywords in the format specified by the user, following the user’s input requirements. For this task, you will focus on extracting relevant, meaningful terms or concepts and presenting them in the specified format.",
	}

	userMessage := Message{
		Role:    "user",
		Content: fmt.Sprintf(prompt, content),
	}

	requestBody := ChatCompletionRequest{
		Model:    model,
		Messages: []Message{systemMessage, userMessage},
		Options:  options,
	}

	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(respBody))
	}

	msg, err := extractMessageFromResponse(respBody)
	if err != nil {
		return "", err
	}

	return msg, nil
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
