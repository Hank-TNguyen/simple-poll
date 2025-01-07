package poll

import (
	"database/sql"
	"fmt"
)

// Question represents a question record in the DB.
type Question struct {
	ID      int64    `json:"id"`
	PollID  int64    `json:"poll_id"`
	Text    string   `json:"text"`
	Choices []Choice `json:"choices"`
}

// ListQuestions fetches all questions (optionally for a specific poll if needed).
func ListQuestions(db *sql.DB, pollID *int64) ([]Question, error) {
	var rows *sql.Rows
	var err error

	if pollID != nil {
		// If you want to filter by poll_id
		rows, err = db.Query("SELECT id, poll_id, question_text FROM questions WHERE poll_id = ?", *pollID)
	} else {
		// Otherwise, get all questions
		rows, err = db.Query("SELECT id, poll_id, question_text FROM questions")
	}
	if err != nil {
		return nil, fmt.Errorf("ListQuestions: %v", err)
	}
	defer rows.Close()

	questions := []Question{}
	for rows.Next() {
		var q Question
		if err := rows.Scan(&q.ID, &q.PollID, &q.Text); err != nil {
			return nil, fmt.Errorf("ListQuestions scan: %v", err)
		}
		questions = append(questions, q)
	}
	return questions, rows.Err()
}

// GetQuestion returns a single Question by ID.
func GetQuestion(db *sql.DB, questionID int64) (*Question, error) {
	var q Question
	err := db.QueryRow("SELECT id, poll_id, question_text FROM questions WHERE id = ?", questionID).
		Scan(&q.ID, &q.PollID, &q.Text)
	if err == sql.ErrNoRows {
		// No result found
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetQuestion: %v", err)
	}
	return &q, nil
}

// CreateQuestion inserts a new question into the DB.
func CreateQuestion(db *sql.DB, q *Question) error {
	result, err := db.Exec(
		"INSERT INTO questions (poll_id, question_text) VALUES (?, ?)",
		q.PollID, q.Text,
	)
	if err != nil {
		return fmt.Errorf("CreateQuestion: %v", err)
	}
	// Retrieve the newly inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("CreateQuestion LastInsertId: %v", err)
	}
	q.ID = id
	return nil
}

// UpdateQuestion updates the question text for an existing record.
func UpdateQuestion(db *sql.DB, q *Question) error {
	_, err := db.Exec(
		"UPDATE questions SET question_text = ? WHERE id = ?",
		q.Text, q.ID,
	)
	if err != nil {
		return fmt.Errorf("UpdateQuestion: %v", err)
	}
	return nil
}

// DeleteQuestion deletes a question by ID.
func DeleteQuestion(db *sql.DB, questionID int64) error {
	_, err := db.Exec("DELETE FROM questions WHERE id = ?", questionID)
	if err != nil {
		return fmt.Errorf("DeleteQuestion: %v", err)
	}
	return nil
}
