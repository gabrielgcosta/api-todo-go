package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// Init conecta ao RabbitMQ e retorna um Client contendo a conexão e o canal
func Init(url string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declarar a fila de eventos
	_, err = ch.QueueDeclare(
		"task_events", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	log.Println("[RabbitMQ] [INFO] Connected and queue declared")

	return &Client{
		Conn:    conn,
		Channel: ch,
	}, nil
}

// Close fecha o canal e a conexão
func (c *Client) Close() {
	if c.Channel != nil {
		c.Channel.Close()
	}
	if c.Conn != nil {
		c.Conn.Close()
	}
}
