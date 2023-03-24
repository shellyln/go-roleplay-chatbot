package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	apiKey      = ""
	apiEndpoint = ""
)

// App request
type ClientPrompt struct {
	IsBot        bool   `json:"isBot"`
	IsDirective  bool   `json:"isDirective"`
	MyCharName   string `json:"myCharName"`
	YourCharName string `json:"yourCharName"`
	Prompt       string `json:"prompt"`
}

// App request
type ClientPromptReqPayload struct {
	History []ClientPrompt `json:"history"`
}

// App response
type ClientPromptResPayload struct {
	Text string `json:"text"`
}

// OpenAI request
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAI request
type CompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

// OpenAI response
type Choice struct {
	Index        int     `json:"index"`
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
}

// OpenAI response
type CompletionResponse struct {
	Error   string   `json:"error"`
	Choices []Choice `json:"choices"`
}

// OpenAI エンドポイントにリクエストを送る
func sendChatRequest(requestBody CompletionRequest) (*CompletionResponse, error) {
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{Timeout: time.Second * 60}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response CompletionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response.Error != "" || len(response.Choices) == 0 {
		return nil, errors.New("No choices or has error: " + response.Error)
	}

	return &response, nil
}
