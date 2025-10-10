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
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "1", "t":
			m.currentView = ViewTasks
			return m, nil
		case "2", "c":
			m.currentView = ViewCalendar
			return m, nil
		case "3", "s":
			m.currentView = ViewClasses
			return m, nil
		case "4", "g":
			m.currentView = ViewGrades
			return m, nil
		case "5", "n":
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

	// Update current screen
	var cmd tea.Cmd
	switch m.currentView {
	case ViewTasks:
		m.taskScreen, cmd = m.taskScreen.Update(msg)
	}

	return m, cmd
}

// View renders the application
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Title bar
	titleBar := styles.TitleBar.
		Width(m.width).
		Render("🎓 UniCLI - Student Organization Manager")

	// Main content area
	contentHeight := m.height - 4 // Reserve space for title and status bars

	// Sidebar
	sidebar := m.renderSidebar(contentHeight)

	// Main content
	var content string
	switch m.currentView {
	case ViewTasks:
		content = m.taskScreen.View()
	case ViewCalendar:
		content = "📅 Calendar View (Coming Soon)"
	case ViewClasses:
		content = "🎒 Classes View (Coming Soon)"
	case ViewGrades:
		content = "📊 Grades View (Coming Soon)"
	case ViewNotes:
		content = "📝 Notes View (Coming Soon)"
	case ViewStats:
		content = "📈 Statistics View (Coming Soon)"
	case ViewSettings:
		content = "⚙️  Settings View (Coming Soon)"
	}

	contentPanel := styles.Panel.
		Width(m.width - 22).
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
func (m Model) renderSidebar(height int) string {
	views := []struct {
		view  View
		key   string
		label string
		icon  string
	}{
		{ViewTasks, "1", "Tasks", "📋"},
		{ViewCalendar, "2", "Calendar", "📅"},
		{ViewClasses, "3", "Classes", "🎒"},
		{ViewGrades, "4", "Grades", "📊"},
		{ViewNotes, "5", "Notes", "📝"},
		{ViewStats, "6", "Stats", "📈"},
		{ViewSettings, "7", "Settings", "⚙️"},
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
		
		item := style.Render(fmt.Sprintf("%s %s %s", v.icon, v.key, v.label))
		items = append(items, item)
	}

	sidebarContent := lipgloss.JoinVertical(lipgloss.Left, items...)

	return styles.Panel.
		Width(20).
		Height(height).
		Render(sidebarContent)
}

// renderStatusBar renders the bottom status bar
func (m Model) renderStatusBar() string {
	leftHelp := styles.Shortcut.Render("?") + 
		styles.ShortcutText.Render(" Help  ") +
		styles.Shortcut.Render("q") + 
		styles.ShortcutText.Render(" Quit")

	rightHelp := styles.Dimmed.Render("Navigate: 1-7 or arrows")

	statusContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftHelp,
		strings.Repeat(" ", max(0, m.width-lipgloss.Width(leftHelp)-lipgloss.Width(rightHelp)-2)),
		rightHelp,
	)

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
