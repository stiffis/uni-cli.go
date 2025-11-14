package models

import "github.com/google/uuid"

// Category represents a category for events or tasks
type Category struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// NewCategory creates a new category with a unique ID
func NewCategory(name, color string) *Category {
	return &Category{
		ID:    uuid.New().String(),
		Name:  name,
		Color: color,
	}
}

func (c Category) FilterValue() string { return c.Name }
func (c Category) Title() string       { return c.Name }
func (c Category) Description() string { return "Color: " + c.Color }
