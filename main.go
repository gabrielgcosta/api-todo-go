package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"todo_api/database"
	"todo_api/middleware"
	"todo_api/task"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("The environment variable DB_URL is not set")
	}

	db, err := database.Init(dbURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer db.Close()

	taskRepo := task.NewRepository(db)
	taskHandler := task.NewHandler(taskRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", taskHandler.Create)
	mux.HandleFunc("GET /tasks", taskHandler.List)
	mux.HandleFunc("PUT /tasks/{id}", taskHandler.Update)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.Delete)

	loggedMux := middleware.Logger(mux)

	fmt.Println("Server running on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", loggedMux))
}
