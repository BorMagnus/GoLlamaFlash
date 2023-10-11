package utils

import (
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
	flashcards, err := flashcards.GenerateFlashcards(filePath)
	if err != nil {
		fmt.Println(flashcards, err)
		return err
	} else {
		appendToFile(filePath, flashcards)
		return err
	}
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
	_, err = file.WriteString("\n---\n## Flashcards\n" + flashcards + "\n")
	if err != nil {
		fmt.Println("File write error:", err)
		return
	}
}
