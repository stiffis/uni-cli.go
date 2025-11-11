package models

import "time"

// Subtask represents a single item in a task's checklist.
type Subtask struct {
	ID          int       `json:"id"`
	TaskID      string    `json:"task_id"`
	Title       string    `json:"title"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
}
