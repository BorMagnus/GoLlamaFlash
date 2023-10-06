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
