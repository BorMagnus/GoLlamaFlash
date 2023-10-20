package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

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

type FlashcardsResponse struct {
	Flashcards []Flashcard `json:"flashcards"`
}

type Flashcard struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// NewClient returns a configured HTTP client.
func NewClient() *http.Client { // TODO: Change time
	return &http.Client{
		Timeout: time.Second * 20,
	}
}

// ParseAPIResponse extracts relevant data from the API response.
func ParseAPIResponse(apiResponse string) (FlashcardsResponse, error) {
	var fullAPIResponse string
	var flashcardsResponse FlashcardsResponse

	responses := strings.Split(apiResponse, "\n")
	for _, response := range responses {
		if response == "" {
			continue
		}
		var apiResponseObj Response
		err := json.Unmarshal([]byte(response), &apiResponseObj)
		if err != nil {
			return flashcardsResponse, fmt.Errorf("error unmarshaling response: %w", err)
		} else {
			fullAPIResponse += apiResponseObj.APIResponse
		}
	}

	startIdx := strings.Index(fullAPIResponse, "{") // TODO: Better way to find JSON in output
	endIdx := strings.LastIndex(fullAPIResponse, "}") + 1

	if startIdx == -1 || endIdx == -1 {
		return flashcardsResponse, fmt.Errorf("valid JSON not found in the API response")
	}

	jsonPart := fullAPIResponse[startIdx:endIdx]

	err := json.Unmarshal([]byte(jsonPart), &flashcardsResponse)
	if err != nil {
		return flashcardsResponse, fmt.Errorf("error parsing flashcards JSON: %w", err)
	}

	return flashcardsResponse, nil
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
