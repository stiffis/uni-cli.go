package models

import "time"

// CalendarItem is an interface for items that can be displayed on a calendar.
// Both Task and Event models will implement this interface.
type CalendarItem interface {
	GetID() string
	GetTitle() string
	GetStartTime() time.Time
	GetEndTime() *time.Time
	IsAllDay() bool
	GetType() string
}
