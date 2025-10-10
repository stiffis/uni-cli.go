package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Colors - Inspired by lazygit
var (
	Primary   = lipgloss.Color("#7C3AED") // Purple
	Secondary = lipgloss.Color("#06B6D4") // Cyan
	Success   = lipgloss.Color("#10B981") // Green
	Warning   = lipgloss.Color("#F59E0B") // Amber
	Danger    = lipgloss.Color("#EF4444") // Red
	Info      = lipgloss.Color("#3B82F6") // Blue
	Muted     = lipgloss.Color("#6B7280") // Gray
	
	Background       = lipgloss.Color("#1E1E2E")
	BackgroundLight  = lipgloss.Color("#313244")
	Foreground       = lipgloss.Color("#CDD6F4")
	Border           = lipgloss.Color("#45475A")
	BorderFocused    = lipgloss.Color("#7C3AED")
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

	// Title is for panel titles
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Primary).
		Padding(0, 1)

	// TitleBar is the top bar style
	TitleBar = lipgloss.NewStyle().
			Background(Primary).
			Foreground(lipgloss.Color("#FFFFFF")).
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
				Foreground(lipgloss.Color("#FFFFFF")).
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
