package worker

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"todo_api/rabbitmq"
)

type Consumer struct {
	rbClient *rabbitmq.Client
	email    string
}

func NewConsumer(rbClient *rabbitmq.Client) *Consumer {
	email := os.Getenv("NOTIFICATION_EMAIL")
	if email == "" {
		log.Println("[Consumer] [WARNING] NOTIFICATION_EMAIL not set, falling back to default")
		email = "teste@gmail.com"
	}
	return &Consumer{
		rbClient: rbClient,
		email:    email,
	}
}

func (c *Consumer) Start() {
	if c.rbClient == nil || c.rbClient.Channel == nil {
		log.Fatal("[Consumer] [ERROR] RabbitMQ client is not connected")
	}

	msgs, err := c.rbClient.Channel.Consume(
		"task_events", // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		log.Fatalf("[Consumer] [ERROR] Failed to register a consumer: %v", err)
	}

	log.Printf("[Consumer] [INFO] Waiting for messages. To exit press CTRL+C")

	go func() {
		for d := range msgs {
			var event TaskEvent
			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				log.Printf("[Consumer] [ERROR] Error decoding message: %v", err)
				continue
			}

			// Simulate processing / sending email
			time.Sleep(500 * time.Millisecond)

			action := ""
			switch event.Type {
			case EventCreated:
				action = "criada"
			case EventUpdated:
				action = "atualizada"
			case EventDeleted:
				action = "excluída"
			default:
				action = "modificada"
			}

			log.Printf("=====================================================")
			log.Printf("[EMAIL SIMULATOR] Enviando e-mail para: %s", c.email)
			log.Printf("[EMAIL SIMULATOR] Assunto: Notificação de Tarefa (%s)", action)
			log.Printf("[EMAIL SIMULATOR] Corpo do E-mail:")
			log.Printf("  A tarefa com ID %d foi %s com sucesso no sistema.", event.TaskID, action)
			if event.Title != "" {
				log.Printf("  Título: %s", event.Title)
			}
			log.Printf("  Data/Hora: %s", event.Timestamp.Format(time.RFC1123))
			log.Printf("=====================================================")
		}
	}()
}
