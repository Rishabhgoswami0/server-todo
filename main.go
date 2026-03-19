package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Rishabhgoswami0/shared-go/database"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// Authentication removed, user details are hardcoded directly in handlers

func main() {
	_ = godotenv.Load() // Ignore error if .env doesn't exist

	// 1. Connect to Database using shared package
	connStr := os.Getenv("POSTGRES_URI")
	if connStr == "" {
		log.Fatal("POSTGRES_URI environment variable is required")
	}
	db, err := database.ConnectPostgres(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 2. Routing with Go 1.22+ ServeMux
	mux := http.NewServeMux()

	mux.HandleFunc("POST /todos", createTodoHandler(db))
	mux.HandleFunc("GET /todos", getTodosHandler(db))
	mux.HandleFunc("PUT /todos/{id}", updateTodoHandler(db))
	mux.HandleFunc("DELETE /todos/{id}", deleteTodoHandler(db))

	// 3. Mux directly handles requests with Auth removed
	handler := mux

	port := ":3002"
	log.Printf("server-todo running on port %s", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}


// CRUD Handlers

func createTodoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := "default-user"

		var t Todo
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		var id int
		// Insert into DB. The handler implicitly assigns the extracted UserID
		err := db.QueryRow("INSERT INTO todos (user_id, title, completed) VALUES ($1, $2, $3) RETURNING id",
			userID, t.Title, t.Completed).Scan(&id)
		
		if err != nil {
			http.Error(w, "failed to create todo", http.StatusInternalServerError)
			log.Printf("DB Insert Error: %v", err)
			return
		}

		t.ID = id
		t.UserID = userID
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)
	}
}

func getTodosHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := "default-user"

		// Filter purely by user_id
		rows, err := db.Query("SELECT id, user_id, title, completed FROM todos WHERE user_id = $1", userID)
		if err != nil {
			http.Error(w, "failed to query todos", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		todos := []Todo{}
		for rows.Next() {
			var t Todo
			if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Completed); err != nil {
				continue
			}
			todos = append(todos, t)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todos)
	}
}

func updateTodoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := "default-user"
		todoID := r.PathValue("id")
		
		var t Todo
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		// Crucial: AND user_id = $4 ensures users can only edit their own
		result, err := db.Exec("UPDATE todos SET title = $1, completed = $2 WHERE id = $3 AND user_id = $4",
			t.Title, t.Completed, todoID, userID)
		if err != nil {
			http.Error(w, "failed to update todo", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			http.Error(w, "todo not found or not owned by user", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status": "updated"}`)
	}
}

func deleteTodoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := "default-user"
		todoID := r.PathValue("id")

		// Crucial: AND user_id = $2 ensures users can only delete their own
		result, err := db.Exec("DELETE FROM todos WHERE id = $1 AND user_id = $2", todoID, userID)
		if err != nil {
			http.Error(w, "failed to delete todo", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			http.Error(w, "todo not found or not owned by user", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status": "deleted"}`)
	}
}
