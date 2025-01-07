package poll

import (
	"database/sql"
	"fmt"
)

// Choice represents a choice record in the DB.
type Choice struct {
	ID         int64  `json:"id"`
	QuestionID int64  `json:"question_id"`
	Text       string `json:"choice_text"`
}

// ListChoices fetches all choices (optionally for a specific question if needed).
func ListChoices(db *sql.DB, questionID *int64) ([]Choice, error) {
	var rows *sql.Rows
	var err error

	if questionID != nil {
		// If you want to filter by question_id
		rows, err = db.Query("SELECT id, question_id, choice_text FROM choices WHERE question_id = ?", *questionID)
	} else {
		// Otherwise, get all choices
		rows, err = db.Query("SELECT id, question_id, choice_text FROM choices")
	}
	if err != nil {
		return nil, fmt.Errorf("ListChoices: %v", err)
	}
	defer rows.Close()

	choices := []Choice{}
	for rows.Next() {
		var c Choice
		if err := rows.Scan(&c.ID, &c.QuestionID, &c.Text); err != nil {
			return nil, fmt.Errorf("ListChoices scan: %v", err)
		}
		choices = append(choices, c)
	}
	return choices, rows.Err()
}

// GetChoice returns a single Choice by ID.
func GetChoice(db *sql.DB, choiceID int64) (*Choice, error) {
	var c Choice
	err := db.QueryRow("SELECT id, question_id, choice_text FROM choices WHERE id = ?", choiceID).
		Scan(&c.ID, &c.QuestionID, &c.Text)
	if err == sql.ErrNoRows {
		// No result found
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetChoice: %v", err)
	}
	return &c, nil
}

// CreateChoice inserts a new choice into the DB.
func CreateChoice(db *sql.DB, c *Choice) error {
	result, err := db.Exec(
		"INSERT INTO choices (question_id, choice_text) VALUES (?, ?)",
		c.QuestionID, c.Text,
	)
	if err != nil {
		return fmt.Errorf("CreateChoice: %v", err)
	}
	// Retrieve the newly inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("CreateChoice LastInsertId: %v", err)
	}
	c.ID = id
	return nil
}

// UpdateChoice updates the choice text for an existing record.
func UpdateChoice(db *sql.DB, c *Choice) error {
	_, err := db.Exec(
		"UPDATE choices SET choice_text = ? WHERE id = ?",
		c.Text, c.ID,
	)
	if err != nil {
		return fmt.Errorf("UpdateChoice: %v", err)
	}
	return nil
}

// DeleteChoice deletes a choice by ID.
func DeleteChoice(db *sql.DB, choiceID int64) error {
	_, err := db.Exec("DELETE FROM choices WHERE id = ?", choiceID)
	if err != nil {
		return fmt.Errorf("DeleteChoice: %v", err)
	}
	return nil
}
