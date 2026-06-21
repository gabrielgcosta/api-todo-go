package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	
	"todo_api/database"
	"todo_api/middleware"
	"todo_api/rabbitmq"
	"todo_api/task"
	"todo_api/worker"
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

	rbURL := os.Getenv("RABBITMQ_URL")
	if rbURL == "" {
		log.Fatal("The environment variable RABBITMQ_URL is not set")
	}

	var rbClient *rabbitmq.Client
	for i := 0; i < 10; i++ {
		rbClient, err = rabbitmq.Init(rbURL)
		if err == nil {
			break
		}
		log.Printf("Waiting for RabbitMQ... (%v)", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Unable to connect to RabbitMQ after retries:", err)
	}
	defer rbClient.Close()

	// Produtor (usado pelo handler)
	asyncWorker := worker.NewWorker(rbClient)

	// Consumidor (roda em background)
	consumer := worker.NewConsumer(rbClient)
	consumer.Start()

	taskRepo := task.NewRepository(db)
	taskHandler := task.NewHandler(taskRepo, asyncWorker)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", taskHandler.Create)
	mux.HandleFunc("GET /tasks", taskHandler.List)
	mux.HandleFunc("PUT /tasks/{id}", taskHandler.Update)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.Delete)

	loggedMux := middleware.Logger(mux)

	fmt.Println("Server running on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", loggedMux))
}
