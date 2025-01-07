package poll

import (
	"database/sql"
	"errors"
	"time"
)

type Poll struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CreatedBy   int64      `json:"created_by"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	CreatedAt   time.Time  `json:"created_at"`
	Questions   []Question `json:"questions"`
}

func GetPoll(db *sql.DB, pollID int64) (*Poll, error) {
	pollQuery := `
		SELECT id, title, description, created_by, start_date, end_date, created_at
		FROM polls
		WHERE id = ?
	`
	row := db.QueryRow(pollQuery, pollID)

	var p Poll
	if err := row.Scan(
		&p.ID,
		&p.Title,
		&p.Description,
		&p.CreatedBy,
		&p.StartDate,
		&p.EndDate,
		&p.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	qQuery := `
		SELECT id, poll_id, question_text
		FROM questions
		WHERE poll_id = ?
	`
	qRows, err := db.Query(qQuery, p.ID)
	if err != nil {
		return nil, err
	}
	defer qRows.Close()

	var questions []Question
	for qRows.Next() {
		var q Question
		if err := qRows.Scan(&q.ID, &q.PollID, &q.Text); err != nil {
			return nil, err
		}

		cQuery := `
			SELECT id, question_id, choice_text
			FROM choices
			WHERE question_id = ?
		`
		cRows, err := db.Query(cQuery, q.ID)
		if err != nil {
			return nil, err
		}

		var choices []Choice
		for cRows.Next() {
			var c Choice
			if err := cRows.Scan(&c.ID, &c.QuestionID, &c.Text); err != nil {
				cRows.Close()
				return nil, err
			}
			choices = append(choices, c)
		}
		cRows.Close()

		q.Choices = choices
		questions = append(questions, q)
	}
	p.Questions = questions

	return &p, nil
}

// CreatePoll inserts a new poll into the database.
func CreatePoll(db *sql.DB, poll *Poll) error {
	// Insert statement returning the last inserted ID
	query := `
        INSERT INTO polls (title, description, created_by, start_date, end_date)
        VALUES (?, ?, ?, ?, ?)
    `
	result, err := db.Exec(query,
		poll.Title,
		poll.Description,
		poll.CreatedBy,
		poll.StartDate,
		poll.EndDate,
	)
	if err != nil {
		return err
	}

	// Retrieve the auto-incremented ID
	newID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	poll.ID = newID
	return nil
}

// ListPolls retrieves all polls (for example).
func ListPolls(db *sql.DB) ([]Poll, error) {
	query := `
        SELECT id, title, description, created_by, start_date, end_date, created_at
        FROM polls
        ORDER BY created_at DESC
    `
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var polls []Poll
	for rows.Next() {
		var p Poll
		err := rows.Scan(
			&p.ID, &p.Title, &p.Description, &p.CreatedBy,
			&p.StartDate, &p.EndDate, &p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		polls = append(polls, p)
	}
	return polls, nil
}

// DeletePoll removes a poll by ID.
func DeletePoll(db *sql.DB, pollID int64) error {
	query := `
        DELETE FROM polls
        WHERE id = ?
    `
	result, err := db.Exec(query, pollID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no rows deleted; poll not found")
	}
	return nil
}
