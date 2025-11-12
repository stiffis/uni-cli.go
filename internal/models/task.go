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

// GetID returns the ID of the task
func (t *Task) GetID() string {
	return t.ID
}

// GetTitle returns the title of the task
func (t *Task) GetTitle() string {
	return t.Title
}

// GetStartTime returns the due date of the task as the start time
func (t *Task) GetStartTime() time.Time {
	if t.DueDate == nil {
		return time.Time{} // Return zero time if no due date
	}
	return *t.DueDate
}

// GetEndTime returns nil for tasks as they typically don't have an end time
func (t *Task) GetEndTime() *time.Time {
	return nil
}

// IsAllDay returns true if the task has a due date but no specific time
func (t *Task) IsAllDay() bool {
	if t.DueDate == nil {
		return false
	}
	// Consider it all-day if the time part is zero
	return t.DueDate.Hour() == 0 && t.DueDate.Minute() == 0 && t.DueDate.Second() == 0
}

// GetType returns the type of the calendar item (task)
func (t *Task) GetType() string {
	return "task"
}
