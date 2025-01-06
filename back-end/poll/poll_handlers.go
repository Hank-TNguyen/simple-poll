package poll

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// PollRouter is the main entry point for /api/polls routes.
func PollRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/") // might be "" or "123"
		parts := strings.Split(path, "/")

		// Example logic:
		if r.Method == http.MethodGet {
			switch len(parts) {
			case 1:
				// GET /api/polls/ => list all polls
				if parts[0] == "" {
					listPollsHandler(db, w, r)
				} else {
					// GET /api/polls/123 => get that poll
					getPollHandler(db, w, r, parts[0])
				}
			default:
				http.NotFound(w, r)
			}
		} else if r.Method == http.MethodPost && len(parts) == 1 && parts[0] == "" {
			// POST /api/polls/
			createPollHandler(db, w, r)
		} else if r.Method == http.MethodDelete && len(parts) == 1 {
			// DELETE /api/polls/123
			deletePollHandler(db, w, r, parts[0])
		} else {
			http.NotFound(w, r)
		}
	})

	return mux
}

func listPollsHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	polls, err := ListPolls(db)
	if err != nil {
		log.Printf("Error listing polls: %v", err)
		http.Error(w, "Failed to list polls", http.StatusInternalServerError)
		return
	}
	writeJSON(w, polls)
}

func getPollHandler(db *sql.DB, w http.ResponseWriter, r *http.Request, idParam string) {
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid poll ID", http.StatusBadRequest)
		return
	}

	poll, err := GetPoll(db, id)
	if err != nil {
		log.Printf("Error getting poll: %v", err)
		http.Error(w, "Failed to get poll", http.StatusInternalServerError)
		return
	}
	if poll == nil {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, poll)
}

func createPollHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var p Poll
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("Error decoding poll: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// (Optional) parse or set StartDate/EndDate. Example:
	now := time.Now()
	if p.StartDate == nil {
		p.StartDate = &now
	}
	future := now.Add(24 * time.Hour)
	if p.EndDate == nil {
		p.EndDate = &future
	}

	err := CreatePoll(db, &p)
	if err != nil {
		log.Printf("Error creating poll: %v", err)
		http.Error(w, "Failed to create poll", http.StatusInternalServerError)
		return
	}

	writeJSON(w, p)
}

func deletePollHandler(db *sql.DB, w http.ResponseWriter, r *http.Request, idParam string) {
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid poll ID", http.StatusBadRequest)
		return
	}

	err = DeletePoll(db, id)
	if err != nil {
		log.Printf("Error deleting poll: %v", err)
		http.Error(w, "Failed to delete poll", http.StatusInternalServerError)
		return
	}

	// Return a simple success message
	writeJSON(w, map[string]string{"message": "Poll deleted"})
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
