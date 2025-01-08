package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	poll "simple-poll/poll"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Read environment variables for DB connection
	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbHost := os.Getenv("DATABASE_HOST")
	dbName := os.Getenv("DATABASE_DB")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Could not open DB connection: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Could not ping DB: %v", err)
	}

	// Use a ServeMux to handle all routes
	mux := http.NewServeMux()

	// Example route
	mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"message": "Hello from Go backend!"}`)
	})

	// Attach poll routes at /api/polls
	mux.Handle("/api/polls/", http.StripPrefix("/api/polls", poll.PollRouter(db)))
	mux.Handle("/api/questions/", http.StripPrefix("/api/questions", poll.QuestionRouter(db)))
	mux.Handle("/api/choices/", http.StripPrefix("/api/choices", poll.ChoiceRouter(db)))

	// Wrap the mux with our CORS middleware
	handlerWithCORS := corsMiddleware(mux)

	log.Println("Backend running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", handlerWithCORS))
}

// corsMiddleware sets CORS headers on every request
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from your React app’s URL/port
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3001")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Allowed methods and headers
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If this is a preflight request, return 200 directly
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Otherwise, pass the request along
		next.ServeHTTP(w, r)
	})
}
