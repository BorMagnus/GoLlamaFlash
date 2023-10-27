package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// The Ollama model response
type Response struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	APIResponse        string `json:"response"`
	Done               bool   `json:"done"`
	TotalDuration      int64  `json:"total_duration"`
	LoadDuration       int64  `json:"load_duration"`
	SampleCount        int    `json:"sample_count"`
	SampleDuration     int64  `json:"sample_duration"`
	PromptEvalCount    int    `json:"prompt_eval_count"`
	PromptEvalDuration int64  `json:"prompt_eval_duration"`
	EvalCount          int    `json:"eval_count"`
	EvalDuration       int64  `json:"eval_duration"`
	Context            []int  `json:"context"`
}

// NewClient returns a configured HTTP client.
func NewClient() *http.Client { // TODO: Change time
	return &http.Client{
		Timeout: time.Second * 200,
	}
}

// FetchAPIResponse makes the API call and returns the body.
func FetchAPIResponse(client *http.Client, url string, jsonData []byte) ([]byte, error) {
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// PrepareRequestData prepares JSON data for API call.
func PrepareRequestData(content string) ([]byte, error) {
	data := map[string]interface{}{
		"model":  "flashcards",
		"prompt": content,
	}

	return json.Marshal(data)
}
