package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Finished  bool      `json:"finished"`
	CreatedAt time.Time `json:"created_at"`
}

var db *sql.DB

func main() {
	var err error

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("The enviroment variable DB_URL is not set")
	}

	for i := 1; i <= 5; i++ {
		db, err = sql.Open("postgres", dbURL)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		log.Printf("Awaint for database")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Enable to connect to database:", err)
	}
	defer db.Close()

	createTable()

	http.HandleFunc("POST /tasks", createTaskHandler)
	http.HandleFunc("GET /tasks", listTaskHandler)
	http.HandleFunc("PUT /tasks/{id}", updateTaskHandler)
	http.HandleFunc("DELETE /tasks/{id}", deleteTaskHandler)

	fmt.Println("Server running on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		finished BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t Task

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil || strings.TrimSpace(t.Title) == "" {
		http.Error(w, "invalid JSON or empty title", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO tasks (title) VALUES ($1) RETURNING id, finished, created_at"
	err = db.QueryRow(query, t.Title).Scan(&t.ID, &t.Finished, &t.CreatedAt)
	if err != nil {
		http.Error(w, "Error inserting task in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func listTaskHandler(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, title, finished, created_at FROM tasks ORDER BY id ASC"
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed fetching data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Finished, &t.CreatedAt); err != nil {
			http.Error(w, "Failed processing data", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusInternalServerError)
		return
	}

	var t Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid JSON", http.StatusInternalServerError)
		return
	}

	query := "UPDATE tasks SET title = $1, finished = $2 WHERE id = $3"
	result, err := db.Exec(query, t.Title, t.Finished, id)
	if err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task updated successfully"))
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM tasks WHERE id = $1"
	result, err := db.Exec(query, id)
	if err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	linhasAfetadas, _ := result.RowsAffected()
	if linhasAfetadas == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
