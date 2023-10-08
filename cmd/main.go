package main

import (
	"flag"
	"fmt"

	"github.com/BorMagnus/goLlamaFlash/pkg/utils"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "", "Path to the file or directory")
	flag.Parse()

	if path == "" {
		fmt.Println("No path provided. Run: go run your_program.go --path your_file_name")
		return
	}

	utils.ProcessPath(path)
}

// TODO: Change from using the model with cmd to the model API
// TODO: Improve on promt or change settings for the model
// TODO: Implement Goroutines for concurrent flashcard generation.
// TODO: Fix error handeling after adding new error checks
