package flashcards

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

func GenerateFlashcards(filePath string) (string, error) {

	url := "http://localhost:11434/api/generate"

	if hasExistingFlashcards(filePath) {
		return "", fmt.Errorf("skipping %s as it already has flashcards", filePath)
	}

	println("Generating flashcards for:", filePath)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read file error: %w", err)
	}

	data := map[string]interface{}{
		"model": "flashcards",
		"prompt": `Act as a writer. Generate flashcards with questions and answers
		about the notes given. Keep each flashcard concise and
		informative. Only include theoretical questions.
		Output only the text and nothing else, do
		not chat, no preamble, get to the point. Only respond with flashcards. 
		Notes: ` + string(content),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling data:", err)
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}
	var fullAPIResponse string
	responses := strings.Split(string(body), "\n")
	for _, response := range responses {
		if response == "" {
			continue
		}
		var apiResponse Response
		err = json.Unmarshal([]byte(response), &apiResponse)
		if err != nil {
			fmt.Println("Error unmarshaling response:", err)
		} else {
			fullAPIResponse += apiResponse.APIResponse

		}
	}

	flashcards := string(fullAPIResponse)

	fmt.Print(flashcards)

	return flashcards, nil
}

func hasExistingFlashcards(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		println("File open error:", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "## Flashcards") {
			return true
		}
	}

	return false
}
