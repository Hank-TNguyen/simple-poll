package poll

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// ChoiceRouter is the main entry point for /api/choices routes.
// Example usage:
//
//	mux.Handle("/api/choices/", http.StripPrefix("/api/choices", ChoiceRouter(db)))
func ChoiceRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// For instance, /api/choices/ or /api/choices/123
		path := strings.TrimPrefix(r.URL.Path, "/") // might be "" or "123"
		parts := strings.Split(path, "/")

		switch r.Method {
		case http.MethodGet:
			// GET /api/choices/ => list choices
			// GET /api/choices/123 => get choice by ID
			if len(parts) == 1 && parts[0] == "" {
				// e.g. GET /api/choices/ => list
				listChoicesHandler(db, w, r)
				return
			} else if len(parts) == 1 {
				// e.g. GET /api/choices/123 => get by ID
				getChoiceHandler(db, w, r, parts[0])
				return
			}
			http.NotFound(w, r)

		case http.MethodPost:
			// POST /api/choices/ => create new choice
			if len(parts) == 1 && parts[0] == "" {
				createChoiceHandler(db, w, r)
				return
			}
			http.NotFound(w, r)

		case http.MethodPut:
			// PUT /api/choices/123 => update
			if len(parts) == 1 {
				updateChoiceHandler(db, w, r, parts[0])
				return
			}
			http.NotFound(w, r)

		case http.MethodDelete:
			// DELETE /api/choices/123 => delete
			if len(parts) == 1 {
				deleteChoiceHandler(db, w, r, parts[0])
				return
			}
			http.NotFound(w, r)

		default:
			http.NotFound(w, r)
		}
	})

	return mux
}

// listChoicesHandler handles listing all choices,
// optionally you could filter by question_id if you want (using a query param).
func listChoicesHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var questionID *int64

	// Example: /api/choices?question_id=10
	queryValues := r.URL.Query()
	if v := queryValues.Get("question_id"); v != "" {
		if idVal, err := strconv.ParseInt(v, 10, 64); err == nil {
			questionID = &idVal
		}
	}

	choices, err := ListChoices(db, questionID)
	if err != nil {
		log.Printf("Error listing choices: %v", err)
		http.Error(w, "Failed to list choices", http.StatusInternalServerError)
		return
	}
	writeJSON(w, choices)
}

// getChoiceHandler handles retrieving a single choice by ID.
func getChoiceHandler(db *sql.DB, w http.ResponseWriter, r *http.Request, idParam string) {
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid choice ID", http.StatusBadRequest)
		return
	}

	choice, err := GetChoice(db, id)
	if err != nil {
		log.Printf("Error getting choice: %v", err)
		http.Error(w, "Failed to get choice", http.StatusInternalServerError)
		return
	}
	if choice == nil {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, choice)
}

// createChoiceHandler handles creating a new choice.
func createChoiceHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var c Choice
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		log.Printf("Error decoding choice: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Here you may want to validate c.QuestionID, etc.
	// For example:
	// if c.QuestionID == 0 {
	//     http.Error(w, "question_id is required", http.StatusBadRequest)
	//     return
	// }

	if err := CreateChoice(db, &c); err != nil {
		log.Printf("Error creating choice: %v", err)
		http.Error(w, "Failed to create choice", http.StatusInternalServerError)
		return
	}

	writeJSON(w, c)
}

// updateChoiceHandler handles updating an existing choice's text.
func updateChoiceHandler(db *sql.DB, w http.ResponseWriter, r *http.Request, idParam string) {
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid choice ID", http.StatusBadRequest)
		return
	}

	var c Choice
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		log.Printf("Error decoding choice: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// We must ensure the ID from the URL path matches the ID in the payload
	// or you can ignore payload ID and use only the URL's ID.
	c.ID = id

	if err := UpdateChoice(db, &c); err != nil {
		log.Printf("Error updating choice: %v", err)
		http.Error(w, "Failed to update choice", http.StatusInternalServerError)
		return
	}

	writeJSON(w, c)
}

// deleteChoiceHandler handles deleting a choice by ID.
func deleteChoiceHandler(db *sql.DB, w http.ResponseWriter, r *http.Request, idParam string) {
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid choice ID", http.StatusBadRequest)
		return
	}

	if err := DeleteChoice(db, id); err != nil {
		log.Printf("Error deleting choice: %v", err)
		http.Error(w, "Failed to delete choice", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "Choice deleted"})
}
