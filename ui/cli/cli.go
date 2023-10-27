package cli

import (
	"flag"
	"fmt"
	"goLlamaFlash/internal/database"
	"goLlamaFlash/internal/filehandler"
	"goLlamaFlash/internal/generator"
	"goLlamaFlash/internal/models"
)

func Run() {
	var path string
	flag.StringVar(&path, "path", "", "Path to the file or directory")
	flag.Parse()

	if path == "" {
		fmt.Println("No path provided. Use --path to specify the file or directory.")
		return
	}

	files, err := filehandler.ReadFiles(path)
	if err != nil {
		fmt.Println("error reading files:", err)
		return
	}

	if err := database.Initialize(); err != nil {
		fmt.Println("error initializing database:", err)
		return
	}

	if err := processFiles(files); err != nil {
		fmt.Println("error processing files:", err)
	}

}

func processFiles(files []string) error {
	for _, file := range files {
		fmt.Println("Looking at: ", file)
		if err := processNote(file); err != nil {
			fmt.Println("error processing file:", err)
		}
	}
	return nil
}

func processNote(file string) error {
	if err := handleNoteExistence(file); err != nil {
		return fmt.Errorf("error handling note existence: %w", err)
	}

	note, err := fetchOrCreateNote(file)
	if err != nil {
		return fmt.Errorf("error fetching or creating note: %w", err)
	}

	if err := handleFlashcards(note); err != nil {
		return fmt.Errorf("error handling flashcards: %w", err)
	}

	return nil
}

func handleNoteExistence(file string) error {
	exists, err := database.NoteExistsByPath(file)
	if err != nil {
		return err
	}

	if exists {
		fmt.Println("\talready has note")
	}
	return nil
}

func fetchOrCreateNote(file string) (models.Note, error) {
	var note models.Note
	exists, err := database.NoteExistsByPath(file)
	if err != nil {
		return note, err
	}

	if exists {
		note, err = database.GetNoteByPath(file)
		if err != nil {
			return note, err
		}
	} else {
		newNote, err := models.NewNote(file)
		if err != nil {
			return note, err
		}
		note, err = database.AddNote(*newNote)
		if err != nil {
			return note, err
		}
		fmt.Println("\tadded note")
	}
	return note, nil
}

func handleFlashcards(note models.Note) error {
	exists, err := database.FlashcardsExistForNote(note.ID)
	if err != nil {
		return err
	}

	count, err := database.FlashcardsCount(note.ID)
	if err != nil {
		return fmt.Errorf("error counting flashcards: %w", err)
	}

	if exists {
		fmt.Printf("\talready has %v flashcards\n", count)
		return nil
	}

	flashcards, err := generator.GenerateFlashcards(note)
	if err != nil {
		return err
	}
	fmt.Println("\tcreated flashcards")
	return database.AddFlashcards(flashcards, note.ID)
}
