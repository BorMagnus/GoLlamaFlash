package models

// Note represents a single note with flashcards.
type Note struct { // TODO: Save content of note?
	ID         int
	FilePath   string
	Flashcards []Flashcard
}

// NewNote initializes a new Note object.
func NewNote(filePath string) (*Note, error) {
	note := &Note{
		FilePath:   filePath,
		Flashcards: []Flashcard{},
	}

	return note, nil
}
