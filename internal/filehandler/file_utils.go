package filehandler

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"goLlamaFlash/internal/flashcards"
)

// Processes either a file or a directory path for flashcard generation
func ProcessPath(path string) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		visitFolder(path)
	case mode.IsRegular():
		handleFile(path)
	}
}

// Function to find all files within directories
func visitFolder(folderPath string) {
	folder, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Println("Read directory error:", err)
		return
	}

	for _, entry := range folder {
		newPath := filepath.Join(folderPath, entry.Name())
		if entry.IsDir() {
			visitFolder(newPath)
		} else if strings.HasSuffix(entry.Name(), ".md") {
			handleFile(newPath)
		}
	}
}

// Handle a single markdown file
func handleFile(filePath string) error {
	if hasExistingFlashcards(filePath) {
		fmt.Println("Skipping, as it already has flashcards:", filePath)
		return nil
	}
	flashcards, err := flashcards.GenerateFlashcards(filePath)
	if err != nil {
		fmt.Println(flashcards, err)
		return err
	} else {
		appendToFile(filePath, flashcards)
		return err
	}
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

// Add the flashcards to the file
func appendToFile(filePath string, flashcards string) {
	// Open the file in append mode
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("File open error:", err)
		return
	}
	defer file.Close()

	// Add a separator and the flashcard content
	_, err = file.WriteString("\n\n---\n## Flashcards\n```Json\n" + flashcards + "\n```")
	if err != nil {
		fmt.Println("File write error:", err)
		return
	}
}
