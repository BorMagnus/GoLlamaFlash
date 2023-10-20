package main

import (
	"flag"
	"fmt"
	"goLlamaFlash/internal/filehandler"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "", "Path to the file or directory")
	flag.Parse()

	if path == "" {
		fmt.Println("No path provided. Run: go run main.go --path your_file_name")
		return
	}

	filehandler.ProcessPath(path)
}

// TODO: Improve on promt or change settings for the model
// TODO: Fix error handeling
// TODO: Create a folder scanner that checks if a there are some changes on the markdown files in the folder
// TODO: Change the prompt to see if there is any flashcards that can be made
// TODO: Improve HTTP Client
// TODO: Check diff between md files
// TODO: Save md file with SQLite?
// TODO: Make model return a JSON with flashcards
