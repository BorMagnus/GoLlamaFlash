package main

import (
	"flag"
	"os"
	"os/exec"
	"strings"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "", "Path to the file or directory")
	flag.Parse()

	if path == "" {
		println("No path provided. Run: go run your_program.go --path your_file_name")
		return
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		println("Error:", err)
		return
	}

	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		visitFolder(path)
	case mode.IsRegular():
		generateFlashcards(path)
	}
}

// Recursive function to find files in folders
func visitFolder(folderPath string) {
	folder, err := os.ReadDir(folderPath)
	if err != nil {
		println("Read directory error:", err)
		return
	}
	for _, path := range folder {
		if path.IsDir() {
			visitFolder(folderPath + "/" + path.Name())
		} else if strings.HasSuffix(path.Name(), ".md") {
			filePath := folderPath + "/" + path.Name()
			generateFlashcards(filePath)
		}
	}
}

// Generate flashcards from the given file.
func generateFlashcards(filePath string) {
	println("Generating flashcards for:", filePath)
	content, err := os.ReadFile(filePath)
	if err != nil {
		println("Read file error: ", err)
		return
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
		println("LLM command error:", err)
		return
	}

	flashcards := string(output)

	// Append the generated flashcards to the notes.
	appendToNotes(filePath, flashcards)

}

// Add the flashcards to the notes.
func appendToNotes(filePath string, flashcards string) {
	// Open the file in append mode.
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		println("File open error:", err)
		return
	}
	defer file.Close()

	// Add a separator and the flashcard content.
	_, err = file.WriteString("\n---\n## Flashcards\n" + flashcards + "\n")
	if err != nil {
		println("File write error:", err)
		return
	}
}

// TODO: Add that it can't generate new flashcards if there already exist flashcards for the notes
// TODO: Change from using the model with cmd to the model API
// TODO: Improve on promt or change settings for the model
// TODO: Restructure the app
