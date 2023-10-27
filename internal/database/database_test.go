// database_test.go
package database

import (
	"database/sql"
	"goLlamaFlash/internal/models"
	"testing"
)

// initTestDB initializes an in-memory database for testing
func initTestDB() {
	// Use SQLite in-memory database for testing
	db, _ = sql.Open("sqlite3", "file::memory:?cache=shared")
	createTable()
}

func TestInitialize(t *testing.T) {
	initTestDB()
}

func TestAddAndRetrieveNote(t *testing.T) {
	initTestDB()

	// Test AddNote
	note := models.Note{FilePath: "test/file/path"}
	addedNote, err := AddNote(note)
	if err != nil {
		t.Errorf("Failed to add note: %v", err)
	}

	// Test GetNoteByPath
	retrievedNote, err := GetNoteByPath("test/file/path")
	if err != nil {
		t.Errorf("Failed to retrieve note: %v", err)
	}

	if addedNote.ID != retrievedNote.ID {
		t.Errorf("Added and retrieved notes don't match")
	}
}

func TestAddAndRetrieveFlashcards(t *testing.T) {
	initTestDB()

	// Add a note first
	note := models.Note{FilePath: "another/test/file/path"}
	addedNote, err := AddNote(note)
	if err != nil {
		t.Errorf("Failed to add note: %v", err)
	}

	// Add flashcards for the note
	flashcards := models.FlashcardsResponse{Flashcards: []models.Flashcard{
		{Question: "What is Golang?", Answer: "A programming language."},
	}}

	err = AddFlashcards(flashcards, addedNote.ID)
	if err != nil {
		t.Errorf("Failed to add flashcards: %v", err)
	}

	// Retrieve all flashcards and verify
	allFlashcards, err := GetAllFlashcards()
	if err != nil {
		t.Errorf("Failed to retrieve flashcards: %v", err)
	}

	if len(allFlashcards) != 1 {
		t.Errorf("Unexpected number of flashcards: got %v, want 1", len(allFlashcards))
	}
}
