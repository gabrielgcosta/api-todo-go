package worker_test

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"todo_api/worker"
)

func TestWorker_QueueEvent_NoRabbitMQ(t *testing.T) {
	// Captura logs para verificar se ele loga o erro e não panica ao receber client nil
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	w := worker.NewWorker(nil)
	w.QueueEvent(worker.TaskEvent{
		Type:      worker.EventCreated,
		TaskID:    1,
		Title:     "Test",
		Timestamp: time.Now(),
	})

	output := buf.String()
	if !strings.Contains(output, "RabbitMQ client is not connected") {
		t.Errorf("Expected log indicating RabbitMQ is not connected, got: %s", output)
	}
}
