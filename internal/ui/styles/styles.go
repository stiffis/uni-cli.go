package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Colors - Kanagawa Wave palette
var (
	// Main colors
	Primary   = lipgloss.Color("#7E9CD8") // crystalBlue - Functions and Titles
	Secondary = lipgloss.Color("#98BB6C") // springGreen - Strings
	Success   = lipgloss.Color("#76946A") // autumnGreen - Git Add
	Warning   = lipgloss.Color("#FF9E3B") // roninYellow - Diagnostic Warning
	Danger    = lipgloss.Color("#E82424") // samuraiRed - Diagnostic Error
	Info      = lipgloss.Color("#6A9589") // waveAqua1 - Diagnostic Info
	Muted     = lipgloss.Color("#727169") // fujiGray - Comments
	AutumnYellow = lipgloss.Color("#DCA561") // Git Change - as requested by user
	
	// Backgrounds
	Background       = lipgloss.Color("#1F1F28") // sumiInk1 - Default background
	BackgroundLight  = lipgloss.Color("#2A2A37") // sumiInk2 - Lighter background
	Foreground       = lipgloss.Color("#DCD7BA") // fujiWhite - Default foreground
	Border           = lipgloss.Color("#54546D") // sumiInk4 - Darker foreground, float borders
	BorderFocused    = lipgloss.Color("#7E9CD8") // crystalBlue - Same as Primary for focus
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
