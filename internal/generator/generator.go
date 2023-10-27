package generator

import (
	"encoding/json"
	"fmt"
	"goLlamaFlash/internal/api"
	"goLlamaFlash/internal/filehandler"
	"goLlamaFlash/internal/models"
	"strings"
)

// GenerateFlashcards generates flashcards using the Ollama API
func GenerateFlashcards(n models.Note) (models.FlashcardsResponse, error) {
	// Initialize the HTTP client
	client := api.NewClient()

	// Read the file content
	notes, err := filehandler.GetNotes(n.FilePath)
	if err != nil {
		return models.FlashcardsResponse{}, err
	}

	// Prepare the request data
	jsonData, err := api.PrepareRequestData(notes)
	if err != nil {
		return models.FlashcardsResponse{}, err
	}

	// Fetch the API response
	body, err := api.FetchAPIResponse(client, "http://localhost:11434/api/generate", jsonData)
	if err != nil {
		return models.FlashcardsResponse{}, err
	}

	// Parse the API response to get flashcards
	flashcardsResponse, err := ParseAPIResponse(string(body))
	if err != nil {
		return models.FlashcardsResponse{}, err
	}

	return flashcardsResponse, nil
}

// ParseAPIResponse extracts relevant data from the API response.
func ParseAPIResponse(apiResponse string) (models.FlashcardsResponse, error) {
	var fullAPIResponse string
	var flashcardsResponse models.FlashcardsResponse

	responses := strings.Split(apiResponse, "\n")
	for _, response := range responses {
		if response == "" {
			continue
		}
		var apiResponseObj api.Response
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
