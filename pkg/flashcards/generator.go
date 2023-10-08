package flashcards

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GenerateFlashcards(filePath string) (string, error) {

	if hasExistingFlashcards(filePath) {
		return "", fmt.Errorf("skipping %s as it already has flashcards", filePath)
	}

	println("Generating flashcards for:", filePath)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read file error: %w", err)
	}

	prompt := `Act as a writer. Generate flashcards with questions and answers
    about the notes given. Keep each flashcard concise and
    informative. Only include theoretical questions.
    Output only the text and nothing else, do
    not chat, no preamble, get to the point. Only respond with flashcards. 
    Notes: ` + string(content)

	cmd := exec.Command("ollama", "run", "llama2", prompt)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("LLM command error: %w", err)
	}

	flashcards := string(output)

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
