package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

func (s TaskStatus) String() string {
	return string(s)
}

// TaskPriority represents the priority of a task
type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
	TaskPriorityUrgent TaskPriority = "urgent"
)

func (p TaskPriority) String() string {
	return string(p)
}

// Task represents a task/todo item
type Task struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	Category    string       `json:"category"`
	Tags        []string     `json:"tags"`
	Subtasks    []Subtask    `json:"subtasks"`
	DueDate     *time.Time   `json:"due_date"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	CompletedAt *time.Time   `json:"completed_at"`
}

// NewTask creates a new task with default values
func NewTask(title string) *Task {
	now := time.Now()
	return &Task{
		ID:        uuid.New().String(),
		Title:     title,
		Status:    TaskStatusPending,
		Priority:  TaskPriorityMedium,
		Tags:      []string{},
		Subtasks:  []Subtask{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// IsOverdue returns true if the task is overdue
func (t *Task) IsOverdue() bool {
	if t.DueDate == nil || t.Status == TaskStatusCompleted {
		return false
	}
	// Compare with the start of the current day
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return t.DueDate.Before(startOfToday)
}

// IsDueToday returns true if the task is due today
func (t *Task) IsDueToday() bool {
	if t.DueDate == nil {
		return false
	}
	now := time.Now()
	return t.DueDate.Year() == now.Year() &&
		t.DueDate.Month() == now.Month() &&
		t.DueDate.Day() == now.Day()
}

// CompletionPercentage calculates the percentage of completed subtasks.
func (t *Task) CompletionPercentage() int {
	if len(t.Subtasks) == 0 {
		return 0
	}

	completed := 0
	for _, st := range t.Subtasks {
		if st.IsCompleted {
			completed++
		}
	}

	return int(float64(completed) / float64(len(t.Subtasks)) * 100)
}

// CompletionRatio returns a string like "(3/5)" for completed subtasks.
func (t *Task) CompletionRatio() string {
	if len(t.Subtasks) == 0 {
		return ""
	}

	completed := 0
	for _, st := range t.Subtasks {
		if st.IsCompleted {
			completed++
		}
	}

	return fmt.Sprintf("(%d/%d)", completed, len(t.Subtasks))
}
