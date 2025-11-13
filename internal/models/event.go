package models

import (
	"time"

	"github.com/google/uuid"
)

// Event represents an event in the calendar
type Event struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	StartDatetime time.Time  `json:"start_datetime"`
	EndDatetime   *time.Time `json:"end_datetime"`
	Type          string     `json:"type"` // e.g., "class", "meeting", "appointment"
	CreatedAt     time.Time  `json:"created_at"`
}

// NewEvent creates a new event with default values
func NewEvent(title string, start time.Time) *Event {
	now := time.Now()
	return &Event{
		ID:            uuid.New().String(),
		Title:         title,
		StartDatetime: start,
		Type:          "event", // Default type
		CreatedAt:     now,
	}
}

// GetID returns the ID of the event
func (e *Event) GetID() string {
	return e.ID
}

// GetTitle returns the title of the event
func (e *Event) GetTitle() string {
	return e.Title
}

// GetStartTime returns the start time of the event
func (e *Event) GetStartTime() time.Time {
	return e.StartDatetime
}

// GetEndTime returns the end time of the event
func (e *Event) GetEndTime() *time.Time {
	return e.EndDatetime
}

// GetType returns the type of the event
func (e *Event) GetType() string {
	return e.Type
}

// IsAllDay checks if the event is an all-day event
func (e *Event) IsAllDay() bool {
	// An event is considered all-day if it has no specific end time
	// or if its duration is exactly 24 hours and starts at midnight.
	if e.EndDatetime == nil {
		return true
	}
	duration := e.EndDatetime.Sub(e.StartDatetime)
	return e.StartDatetime.Hour() == 0 && e.StartDatetime.Minute() == 0 && duration == 24*time.Hour
}
