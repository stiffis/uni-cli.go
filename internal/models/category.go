package models

import "github.com/google/uuid"

type Category struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

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
