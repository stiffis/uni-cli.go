package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Colors - Zen & Friendly palette (inspired by Nord + Catppuccin)
var (
	// Main colors - Soft and calming
	Primary   = lipgloss.Color("#88C0D0") // Soft Blue (peaceful)
	Secondary = lipgloss.Color("#A3BE8C") // Soft Green (zen)
	Success   = lipgloss.Color("#A3BE8C") // Soft Green
	Warning   = lipgloss.Color("#EBCB8B") // Soft Yellow
	Danger    = lipgloss.Color("#BF616A") // Soft Red
	Info      = lipgloss.Color("#81A1C1") // Soft Blue-Gray
	Muted     = lipgloss.Color("#4C566A") // Muted Gray
	
	// Backgrounds - Dark but warm
	Background       = lipgloss.Color("#2E3440") // Warm dark gray
	BackgroundLight  = lipgloss.Color("#3B4252") // Slightly lighter
	Foreground       = lipgloss.Color("#ECEFF4") // Soft white
	Border           = lipgloss.Color("#4C566A") // Subtle border
	BorderFocused    = lipgloss.Color("#88C0D0") // Soft blue focus
)

// Base styles
var (
	// Panel represents a bordered container
	Panel = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Border).
		Padding(0, 1)

	// PanelFocused is a focused panel
	PanelFocused = Panel.Copy().
			BorderForeground(BorderFocused)

	// PanelTarget is a panel that is a move target
	PanelTarget = Panel.Copy().
			BorderForeground(Warning)

	// Title is for panel titles
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Primary).
		Padding(0, 1)

	// TitleBar is the top bar style
	TitleBar = lipgloss.NewStyle().
			Background(Primary).
			Foreground(Background).
			Bold(true).
			Padding(0, 1)

	// StatusBar is the bottom status bar
	StatusBar = lipgloss.NewStyle().
			Background(BackgroundLight).
			Foreground(Muted).
			Padding(0, 1)

	// ListItem is for list items
	ListItem = lipgloss.NewStyle().
			Padding(0, 2)

	// ListItemSelected is for selected list items
	ListItemSelected = ListItem.Copy().
				Background(Primary).
				Foreground(Background).
				Bold(true)

	// Shortcut displays keyboard shortcuts
	Shortcut = lipgloss.NewStyle().
			Foreground(Secondary).
			Bold(true)

	// ShortcutText is the description text
	ShortcutText = lipgloss.NewStyle().
			Foreground(Muted)

	// Error message style
	Error = lipgloss.NewStyle().
		Foreground(Danger).
		Bold(true)

	// Success message style
	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	// Dimmed text
	Dimmed = lipgloss.NewStyle().
		Foreground(Muted)

	// Tag style
	Tag = lipgloss.NewStyle().
		Background(Secondary).
		Foreground(Background).
		Padding(0, 1).
		MarginRight(1).
		Bold(true)
)

// Status colors for tasks
func StatusColor(status string) lipgloss.Color {
	switch status {
	case "completed":
		return Success
	case "in_progress":
		return Info
	case "pending":
		return Warning
	case "cancelled":
		return Muted
	default:
		return Foreground
	}
}

// Priority colors for tasks
func PriorityColor(priority string) lipgloss.Color {
	switch priority {
	case "urgent":
		return Danger
	case "high":
		return Warning
	case "medium":
		return Info
	case "low":
		return Muted
	default:
		return Foreground
	}
}
