package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"todo_api/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventType string

const (
	EventCreated EventType = "CREATED"
	EventUpdated EventType = "UPDATED"
	EventDeleted EventType = "DELETED"
)

type TaskEvent struct {
	Type      EventType `json:"type"`
	TaskID    int       `json:"task_id"`
	Title     string    `json:"title,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type Worker struct {
	rbClient *rabbitmq.Client
}

// NewWorker cria um novo produtor de eventos
func NewWorker(rbClient *rabbitmq.Client) *Worker {
	return &Worker{
		rbClient: rbClient,
	}
}

// QueueEvent publica o evento no RabbitMQ
func (w *Worker) QueueEvent(event TaskEvent) {
	if w.rbClient == nil || w.rbClient.Channel == nil {
		log.Println("[Worker] [ERROR] RabbitMQ client is not connected")
		return
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("[Worker] [ERROR] Failed to marshal event: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = w.rbClient.Channel.PublishWithContext(ctx,
		"",            // exchange
		"task_events", // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		log.Printf("[Worker] [ERROR] Failed to publish a message: %v", err)
		return
	}

	log.Printf("[Worker] [INFO] Published event %s for Task %d", event.Type, event.TaskID)
}

// Start e Stop não são mais loops contínuos aqui, mas os mantemos
// para não quebrar o handler, ou podemos apenas removê-los e ajustar o main.go.
// Para manter a compatibilidade com a struct atual:
func (w *Worker) Start() {
	// não faz nada agora, já que publicamos diretamente no QueueEvent
}

func (w *Worker) Stop() {
	// não faz nada
}
