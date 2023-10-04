package main

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

func main() {
	println("Creating flashcards!\n")

	usr, _ := user.Current()
	dir := usr.HomeDir
	folderPath := dir + "/Documents/Obsidian/Notes"

	files, err := os.ReadDir(folderPath)
	if err != nil {
		println("Directory error: ", err)
		return
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			filePath := filepath.Join(folderPath, file.Name())

			generateFlashcards(filePath)
			appendToNotes()
		}
	}
}

// Generate flashcards from the given file.
func generateFlashcards(filePath string) {
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

	flashcard := string(output)
	println("Generated Flashcard: ", flashcard)

}

// TODO: Create an function to add the flashcards to the notes.
func appendToNotes() {

}
