package models

// FlashcardsResponse represents a collection of flashcards.
type FlashcardsResponse struct {
	Flashcards []Flashcard `json:"flashcards"`
}

// Flashcard represents a single flashcard.
type Flashcard struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
