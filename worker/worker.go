package worker

import (
	"log"
	"time"
)

type EventType string

const (
	EventCreated EventType = "CREATED"
	EventUpdated EventType = "UPDATED"
	EventDeleted EventType = "DELETED"
)

type TaskEvent struct {
	Type      EventType
	TaskID    int
	Title     string
	Timestamp time.Time
}

type Worker struct {
	eventsChan chan TaskEvent
	quit       chan struct{}
}

func NewWorker(bufferSize int) *Worker {
	return &Worker{
		eventsChan: make(chan TaskEvent, bufferSize),
		quit:       make(chan struct{}),
	}
}

func (w *Worker) QueueEvent(event TaskEvent) {
	select {
	case w.eventsChan <- event:
	default:
		log.Printf("[Worker] [WARNING] Queue full. Dropped event: %+v", event)
	}
}

func (w *Worker) Start() {
	go func() {
		log.Println("[Worker] [INFO] Async worker started")
		for {
			select {
			case event := <-w.eventsChan:
				w.processEvent(event)
			case <-w.quit:
				log.Println("[Worker] [INFO] Async worker stopping")
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	close(w.quit)
}

func (w *Worker) processEvent(event TaskEvent) {
	// Simulate heavy background processing (e.g. sending email, auditing)
	time.Sleep(500 * time.Millisecond)
	log.Printf("[Worker] [INFO] Event processed: %s | Task ID: %d | Title: %q | At: %s",
		event.Type, event.TaskID, event.Title, event.Timestamp.Format(time.RFC3339))
}
