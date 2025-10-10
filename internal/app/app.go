package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stiffis/UniCLI/internal/config"
	"github.com/stiffis/UniCLI/internal/database"
	"github.com/stiffis/UniCLI/internal/ui/screens"
	"github.com/stiffis/UniCLI/internal/ui/styles"
)

// View represents different screens in the app
type View int

const (
	ViewTasks View = iota
	ViewCalendar
	ViewClasses
	ViewGrades
	ViewNotes
	ViewStats
	ViewSettings
)

// Model is the main application model
type Model struct {
	db            *database.DB
	cfg           *config.Config
	currentView   View
	width         int
	height        int
	taskScreen    tea.Model
	ready         bool
	err           error
	
	// Command mode (like vim)
	commandMode   bool
	commandInput  string
}

// NewModel creates a new application model
func NewModel(db *database.DB, cfg *config.Config) Model {
	return Model{
		db:          db,
		cfg:         cfg,
		currentView: ViewTasks,
		taskScreen:  screens.NewTaskScreen(db),
	}
}

// Init initializes the application
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		// If in command mode, handle command input
		if m.commandMode {
			switch msg.String() {
			case "enter":
				return m.executeCommand()
			case "esc":
				m.commandMode = false
				m.commandInput = ""
				return m, nil
			case "backspace":
				if len(m.commandInput) > 0 {
					m.commandInput = m.commandInput[:len(m.commandInput)-1]
				}
				return m, nil
			default:
				// Add character to command input
				if len(msg.String()) == 1 {
					m.commandInput += msg.String()
				}
				return m, nil
			}
		}

		// Normal mode key handling
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case ":":
			// Enter command mode
			m.commandMode = true
			m.commandInput = ""
			return m, nil

		case "1":
			m.currentView = ViewTasks
			return m, nil
		case "2":
			m.currentView = ViewCalendar
			return m, nil
		case "3":
			m.currentView = ViewClasses
			return m, nil
		case "4":
			m.currentView = ViewGrades
			return m, nil
		case "5":
			m.currentView = ViewNotes
			return m, nil
		case "6":
			m.currentView = ViewStats
			return m, nil
		case "7":
			m.currentView = ViewSettings
			return m, nil
		}
	}

	// Update current screen (only if not in command mode)
	if !m.commandMode {
		var cmd tea.Cmd
		switch m.currentView {
		case ViewTasks:
			m.taskScreen, cmd = m.taskScreen.Update(msg)
		}
		return m, cmd
	}

	return m, nil
}

// executeCommand processes the command entered by the user
func (m Model) executeCommand() (tea.Model, tea.Cmd) {
	cmd := strings.TrimSpace(m.commandInput)
	
	// Reset command mode
	m.commandMode = false
	m.commandInput = ""
	
	switch cmd {
	case "q", "quit":
		return m, tea.Quit
	case "h", "help":
		// TODO: Show help screen
		m.currentView = ViewSettings // Placeholder for now
		return m, nil
	case "1":
		m.currentView = ViewTasks
	case "2":
		m.currentView = ViewCalendar
	case "3":
		m.currentView = ViewClasses
	case "4":
		m.currentView = ViewGrades
	case "5":
		m.currentView = ViewNotes
	case "6":
		m.currentView = ViewStats
	case "7":
		m.currentView = ViewSettings
	}
	
	return m, nil
}

// View renders the application
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Title bar
	titleLeft := "🎓 UniCLI - Student Organization Manager"
	titleRight := "[:h] Help  [:q] Quit"
	
	// Calculate spacing between left and right text
	leftWidth := lipgloss.Width(titleLeft)
	rightWidth := lipgloss.Width(titleRight)
	titleSpacing := m.width - leftWidth - rightWidth - 2 // -2 for padding
	if titleSpacing < 1 {
		titleSpacing = 1
	}
	
	titleContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		titleLeft,
		strings.Repeat(" ", titleSpacing),
		titleRight,
	)
	
	titleBar := styles.TitleBar.
		Width(m.width).
		Render(titleContent)

	// Main content area
	contentHeight := m.height - 4 // Reserve space for title and status bars

	// Sidebar
	sidebarWidth := 20
	sidebar := m.renderSidebar(contentHeight, sidebarWidth)

	// Main content - account for borders (4 chars: 2 for sidebar border, 2 for content border)
	contentWidth := m.width - sidebarWidth - 4
	if contentWidth < 10 {
		contentWidth = 10
	}

	// Main content
	var content string
	switch m.currentView {
	case ViewTasks:
		content = m.taskScreen.View()
	case ViewCalendar:
		content = "Calendar View (Coming Soon)"
	case ViewClasses:
		content = "Classes View (Coming Soon)"
	case ViewGrades:
		content = "Grades View (Coming Soon)"
	case ViewNotes:
		content = "Notes View (Coming Soon)"
	case ViewStats:
		content = "Statistics View (Coming Soon)"
	case ViewSettings:
		content = "Settings View (Coming Soon)"
	}

	contentPanel := styles.Panel.
		Width(contentWidth).
		Height(contentHeight).
		Render(content)

	// Combine sidebar and content
	mainArea := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebar,
		contentPanel,
	)

	// Status bar
	statusBar := m.renderStatusBar()

	// Combine all elements
	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleBar,
		mainArea,
		statusBar,
	)
}

// renderSidebar renders the navigation sidebar
func (m Model) renderSidebar(height int, width int) string {
	views := []struct {
		view  View
		key   string
		label string
	}{
		{ViewTasks, "1", "Tasks"},
		{ViewCalendar, "2", "Calendar"},
		{ViewClasses, "3", "Classes"},
		{ViewGrades, "4", "Grades"},
		{ViewNotes, "5", "Notes"},
		{ViewStats, "6", "Stats"},
		{ViewSettings, "7", "Settings"},
	}

	var items []string
	for _, v := range views {
		style := lipgloss.NewStyle().Padding(0, 1)
		if v.view == m.currentView {
			style = style.
				Background(styles.Primary).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)
		}
		
		item := style.Render(fmt.Sprintf("[%s] %s", v.key, v.label))
		items = append(items, item)
	}

	sidebarContent := lipgloss.JoinVertical(lipgloss.Left, items...)

	return styles.Panel.
		Width(width).
		Height(height).
		Render(sidebarContent)
}

// renderStatusBar renders the bottom status bar
func (m Model) renderStatusBar() string {
	// If in command mode, show command input
	if m.commandMode {
		commandPrompt := lipgloss.NewStyle().
			Foreground(styles.Primary).
			Bold(true).
			Render(":") + m.commandInput
		
		// Add cursor
		cursor := lipgloss.NewStyle().
			Background(styles.Primary).
			Foreground(lipgloss.Color("#FFFFFF")).
			Render(" ")
		
		return styles.StatusBar.
			Width(m.width).
			Render(commandPrompt + cursor)
	}
	
	// Normal mode status bar
	statusContent := styles.Dimmed.Render("Navigate: 1-7  |  [:] Command mode")

	return styles.StatusBar.
		Width(m.width).
		Render(statusContent)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
