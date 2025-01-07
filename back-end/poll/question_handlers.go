package poll

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// QuestionRouter is the main entry point for /api/questions routes.
// Example usage:
//
//	mux.Handle("/api/questions/", http.StripPrefix("/api/questions", QuestionRouter(db)))
func QuestionRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// For instance, /api/questions/ or /api/questions/123
		path := strings.TrimPrefix(r.URL.Path, "/") // might be "" or "123"
		parts := strings.Split(path, "/")

		switch r.Method {
		case http.MethodGet:
			// GET /api/questions/ => list questions
			// GET /api/questions/123 => get question by ID
			if len(parts) == 1 && parts[0] == "" {
				// e.g. GET /api/questions/ => list
				listQuestionsHandler(db, w, r)
				return
			} else if len(parts) == 1 {
				// e.g. GET /api/questions/123 => get by ID
				getQuestionHandler(db, w, r, parts[0])
				return
			}
			http.NotFound(w, r)

		case http.MethodPost:
			// POST /api/questions/ => create new question
			if len(parts) == 1 && parts[0] == "" {
				createQuestionHandler(db, w, r)
				return
			}
			http.NotFound(w, r)

		case http.MethodPut:
			// PUT /api/questions/123 => update
			if len(parts) == 1 {
				updateQuestionHandler(db, w, r, parts[0])
				return
			}
			http.NotFound(w, r)

		case http.MethodDelete:
			// DELETE /api/questions/123 => delete
			if len(parts) == 1 {
				deleteQuestionHandler(db, w, r, parts[0])
				return
			}
			http.NotFound(w, r)

		default:
			http.NotFound(w, r)
		}
	})

	return mux
}

// listQuestionsHandler handles listing all questions,
// optionally you could filter by poll_id if you want (using a query param).
func listQuestionsHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var pollID *int64

	// Example: /api/questions?poll_id=5
	queryValues := r.URL.Query()
	if v := queryValues.Get("poll_id"); v != "" {
		if idVal, err := strconv.ParseInt(v, 10, 64); err == nil {
			pollID = &idVal
		}
	}

	questions, err := ListQuestions(db, pollID)
	if err != nil {
		log.Printf("Error listing questions: %v", err)
		http.Error(w, "Failed to list questions", http.StatusInternalServerError)
		return
	}
	writeJSON(w, questions)
}

// getQuestionHandler handles retrieving a single question by ID.
func getQuestionHandler(db *sql.DB, w http.ResponseWriter, r *http.Request, idParam string) {
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	question, err := GetQuestion(db, id)
	if err != nil {
		log.Printf("Error getting question: %v", err)
		http.Error(w, "Failed to get question", http.StatusInternalServerError)
		return
	}
	if question == nil {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, question)
}

// createQuestionHandler handles creating a new question.
func createQuestionHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var q Question
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		log.Printf("Error decoding question: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Here you may want to validate PollID, etc.
	// For example:
	// if q.PollID == 0 {
	//     http.Error(w, "poll_id is required", http.StatusBadRequest)
	//     return
	// }

	if err := CreateQuestion(db, &q); err != nil {
		log.Printf("Error creating question: %v", err)
		http.Error(w, "Failed to create question", http.StatusInternalServerError)
		return
	}

	writeJSON(w, q)
}

// updateQuestionHandler handles updating an existing question's text.
func updateQuestionHandler(db *sql.DB, w http.ResponseWriter, r *http.Request, idParam string) {
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	var q Question
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		log.Printf("Error decoding question: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// We must ensure the ID from the URL path matches the ID in the payload
	// or you can ignore payload ID and use only the URL's ID.
	q.ID = id

	if err := UpdateQuestion(db, &q); err != nil {
		log.Printf("Error updating question: %v", err)
		http.Error(w, "Failed to update question", http.StatusInternalServerError)
		return
	}

	writeJSON(w, q)
}

// deleteQuestionHandler handles deleting a question by ID.
func deleteQuestionHandler(db *sql.DB, w http.ResponseWriter, r *http.Request, idParam string) {
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	if err := DeleteQuestion(db, id); err != nil {
		log.Printf("Error deleting question: %v", err)
		http.Error(w, "Failed to delete question", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "Question deleted"})
}
