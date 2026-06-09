package task

import "time"

// Task represents the task model/entity in the database.
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Finished  bool      `json:"finished"`
	CreatedAt time.Time `json:"created_at"`
}
