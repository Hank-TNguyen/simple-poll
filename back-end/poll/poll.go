package poll

import (
	"database/sql"
	"errors"
	"time"
)

// Poll struct matches your 'polls' table columns
type Poll struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CreatedBy   int64      `json:"created_by"`
	StartDate   *time.Time `json:"start_date"` // optional pointer
	EndDate     *time.Time `json:"end_date"`
	CreatedAt   time.Time  `json:"created_at"`
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

// GetPoll retrieves a single poll by ID.
func GetPoll(db *sql.DB, pollID int64) (*Poll, error) {
	query := `
        SELECT id, title, description, created_by, start_date, end_date, created_at
        FROM polls
        WHERE id = ?
    `
	row := db.QueryRow(query, pollID)

	var p Poll
	err := row.Scan(
		&p.ID, &p.Title, &p.Description, &p.CreatedBy,
		&p.StartDate, &p.EndDate, &p.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // no poll found, return nil
	} else if err != nil {
		return nil, err
	}

	return &p, nil
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
