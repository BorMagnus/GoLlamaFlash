package database

import (
	"database/sql"
	"goLlamaFlash/internal/models"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Initialize initializes the database
func Initialize() error {
	var err error
	db, err = sql.Open("sqlite3", "./flashcards.db")
	if err != nil {
		return err
	}

	createTable()
	return nil
}

// createTable creates the necessary tables if they don't already exist
func createTable() {
	// Create Notes table
	createNotesQuery := `CREATE TABLE IF NOT EXISTS notes (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"filePath" TEXT
	);`

	_, err := db.Exec(createNotesQuery)
	if err != nil {
		log.Fatal(err)
	}

	// Create Flashcards table
	createFlashcardsQuery := `CREATE TABLE IF NOT EXISTS flashcards (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"noteID" INTEGER,
		"question" TEXT,
		"answer" TEXT,
		FOREIGN KEY (noteID) REFERENCES notes (id)
	);`

	_, err = db.Exec(createFlashcardsQuery)
	if err != nil {
		log.Fatal(err)
	}
}

// AddNote adds a new Note to the database and returns its ID
func AddNote(n models.Note) (models.Note, error) {
	exists, _ := NoteExistsByPath(n.FilePath)
	if exists {
		return GetNoteByPath(n.FilePath)
	} else {
		statement, err := db.Prepare("INSERT INTO notes (filePath) VALUES (?)")
		if err != nil {
			return n, err
		}
		res, err := statement.Exec(n.FilePath)
		if err != nil {
			return n, err
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return n, err
		}
		n.ID = int(lastID)
		return n, nil
	}
}

// AddFlashcards adds flashcards for a note in the database
func AddFlashcards(flashcards models.FlashcardsResponse, noteID int) error {
	exists, err := FlashcardsExistForNote(noteID)
	if err != nil {
		return err
	}
	if !exists {
		statement, err := db.Prepare("INSERT INTO flashcards (noteID, question, answer) VALUES (?, ?, ?)")
		if err != nil {
			return err
		}

		for _, f := range flashcards.Flashcards {
			_, err := statement.Exec(noteID, f.Question, f.Answer)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetNoteByPath retrieves a Note from the database by its file path
func GetNoteByPath(filePath string) (models.Note, error) {
	var note models.Note
	err := db.QueryRow("SELECT id, filePath FROM notes WHERE filePath = ?", filePath).
		Scan(&note.ID, &note.FilePath)
	if err != nil {
		if err == sql.ErrNoRows {
			return note, nil
		}
		return note, err
	}
	return note, nil
}

// GetAllNotes retrieves all the notes from the database
func GetAllNotes() ([]string, error) {
	var filePaths []string
	rows, err := db.Query("SELECT filePath FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var filePath string
		if err := rows.Scan(&filePath); err != nil {
			return nil, err
		}
		filePaths = append(filePaths, filePath)
	}
	return filePaths, nil
}

// GetAllFlashcards retrieves all the flashcards from the database
func GetAllFlashcards() ([]models.Flashcard, error) {
	var flashcards []models.Flashcard
	rows, err := db.Query("SELECT question, answer FROM flashcards")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var f models.Flashcard
		if err := rows.Scan(&f.Question, &f.Answer); err != nil {
			return nil, err
		}
		flashcards = append(flashcards, f)
	}
	return flashcards, nil
}

// FlashcardsExistForNote checks if flashcards already exist for a given note ID
func FlashcardsExistForNote(noteID int) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM flashcards WHERE noteID = ?", noteID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// NoteExistsByID checks if a note exists based on its ID.
func NoteExistsByID(noteID int) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM notes WHERE id = ?", noteID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// NoteExistsByPath checks if a note exists based on its file path.
func NoteExistsByPath(filePath string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM notes WHERE filePath = ?", filePath).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FlashcardsExistForNote checks if flashcards already exist for a given note ID
func FlashcardsCount(noteID int) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM flashcards WHERE noteID = ?", noteID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, err
}

// Other CRUD operations
