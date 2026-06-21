package worker_test

import (
	"testing"
	"time"
	"todo_api/worker"
)

func TestWorkerLifecycle(t *testing.T) {
	processed := make(chan worker.TaskEvent, 1)

	w := worker.NewWorker(5)
	w.OnProcessed = func(event worker.TaskEvent) {
		processed <- event
	}

	w.Start()
	defer w.Stop()

	w.QueueEvent(worker.TaskEvent{
		Type:      worker.EventCreated,
		TaskID:    42,
		Title:     "Test Task",
		Timestamp: time.Now(),
	})

	select {
	case event := <-processed:
		if event.TaskID != 42 {
			t.Errorf("expected TaskID 42, got %d", event.TaskID)
		}
		if event.Title != "Test Task" {
			t.Errorf("expected Title 'Test Task', got %q", event.Title)
		}
	case <-time.After(1 * time.Second):
		t.Error("timeout waiting for worker to process event")
	}
}
