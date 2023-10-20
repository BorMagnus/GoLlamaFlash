package flashcards

import (
	"encoding/json"
	"fmt"
	"goLlamaFlash/internal/api"
	"os"
)

func GenerateFlashcards(filePath string) (string, error) {
	// Initialize the HTTP client
	client := api.NewClient()

	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read file error: %w", err)
	}

	// Prepare the request data
	jsonData, err := api.PrepareRequestData(string(content))
	if err != nil {
		return "", err
	}

	// Fetch the API response
	body, err := api.FetchAPIResponse(client, "http://localhost:11434/api/generate", jsonData)
	if err != nil {
		return "", err
	}

	// Parse the API response to get flashcards
	flashcardsResponse, err := api.ParseAPIResponse(string(body))
	if err != nil {
		return "", err
	}

	// Convert the parsed flashcards back to a JSON string
	finalJSON, err := json.Marshal(flashcardsResponse)
	if err != nil {
		return "", err
	}

	return string(finalJSON), nil
}
