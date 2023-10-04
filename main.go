package main

import (
	"os"
	"os/exec"
	"strings"
)

func main() {
	println("Creating flashcards!\n")
	inputPath := os.Args[1]
	inputInfo, _ := os.Stat(inputPath)

	if inputInfo.IsDir() {
		visitFolder(inputPath)
	} else if inputInfo.Mode().IsRegular() {
		generateFlashcards(inputPath)
	} else {
		println("Error")
		// error
	}
}

func visitFolder(folderPath string) {
	folder, _ := os.ReadDir(folderPath)
	for _, path := range folder {
		if path.IsDir() {
			visitFolder(folderPath + "/" + path.Name())
		} else if strings.HasSuffix(path.Name(), ".md") {
			filePath := folderPath + "/" + path.Name()
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
