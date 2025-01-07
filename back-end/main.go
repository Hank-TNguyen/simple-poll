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

	// ROUTES
	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"message": "Hello from Go backend!"}`)
	})

	// Attach poll routes at /api/polls
	http.Handle("/api/polls/", http.StripPrefix("/api/polls", poll.PollRouter(db)))
	http.Handle("/api/questions/", http.StripPrefix("/api/questions", poll.QuestionRouter(db)))
	http.Handle("/api/choices/", http.StripPrefix("/api/choices", poll.ChoiceRouter(db)))

	log.Println("Backend running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
