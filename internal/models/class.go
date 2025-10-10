package models

import (
	"time"

	"github.com/google/uuid"
)

// DayOfWeek represents a day of the week
type DayOfWeek int

const (
	Monday DayOfWeek = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

// Schedule represents a class schedule entry
type Schedule struct {
	DayOfWeek DayOfWeek `json:"day_of_week"`
	StartTime string    `json:"start_time"` // Format: "HH:MM"
	EndTime   string    `json:"end_time"`   // Format: "HH:MM"
}

// Class represents a university/school class
type Class struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Professor string     `json:"professor"`
	Room      string     `json:"room"`
	Color     string     `json:"color"`
	Semester  string     `json:"semester"`
	Credits   int        `json:"credits"`
	Schedules []Schedule `json:"schedules"`
	CreatedAt time.Time  `json:"created_at"`
}

// NewClass creates a new class
func NewClass(name string) *Class {
	return &Class{
		ID:        uuid.New().String(),
		Name:      name,
		Schedules: []Schedule{},
		CreatedAt: time.Now(),
		Color:     "#7C3AED", // Default purple
	}
}
