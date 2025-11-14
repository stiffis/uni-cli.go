package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID                string     `json:"id"`
	Title             string     `json:"title"`
	Description       string     `json:"description"`
	StartDatetime     time.Time  `json:"start_datetime"`
	EndDatetime       *time.Time `json:"end_datetime"`
	Type              string     `json:"type"`
	CategoryID        string     `json:"category_id"`
	Category          *Category  `json:"category"`
	RecurrenceRule    string     `json:"recurrence_rule"`
	RecurrenceEndDate *time.Time `json:"recurrence_end_date"`
	CreatedAt         time.Time  `json:"created_at"`
}

func NewEvent(title string, start time.Time) *Event {
	now := time.Now()
	return &Event{
		ID:            uuid.New().String(),
		Title:         title,
		StartDatetime: start,
		Type:          "event",
		CreatedAt:     now,
	}
}

func (e *Event) GetID() string {
	return e.ID
}

func (e *Event) GetTitle() string {
	return e.Title
}

func (e *Event) GetStartTime() time.Time {
	return e.StartDatetime
}

func (e *Event) GetEndTime() *time.Time {
	return e.EndDatetime
}

func (e *Event) GetType() string {
	return e.Type
}

func (e *Event) IsAllDay() bool {
	if e.EndDatetime == nil {
		return true
	}
	duration := e.EndDatetime.Sub(e.StartDatetime)
	return e.StartDatetime.Hour() == 0 && e.StartDatetime.Minute() == 0 && duration == 24*time.Hour
}
