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
	ViewWelcome View = iota
	ViewTasks
	ViewCalendar
	ViewClasses
	ViewGrades
	ViewNotes
	ViewStats
	ViewSettings
)

// Model is the main application model
type Model struct {
	db          *database.DB
	cfg         *config.Config
	currentView View
	width       int
	height      int
	taskScreen  tea.Model
	ready       bool
	err         error

	// Command mode (like vim)
	commandMode  bool
	commandInput string

	// Sidebar navigation
	sidebarMode   bool
	sidebarCursor int
}

// NewModel creates a new application model
func NewModel(db *database.DB, cfg *config.Config) Model {
	return Model{
		db:          db,
		cfg:         cfg,
		currentView: ViewWelcome, // Start with Welcome screen
		taskScreen:  screens.NewTaskScreen(db),
	}
}

// Init initializes the application
func (m Model) Init() tea.Cmd {
	// Initialize the task screen
	return m.taskScreen.Init()
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		// Pass window size to current screen
		if m.currentView == ViewTasks {
			var cmd tea.Cmd
			m.taskScreen, cmd = m.taskScreen.Update(msg)
			return m, cmd
		}
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

		// If in sidebar mode, handle sidebar navigation
		if m.sidebarMode {
			switch msg.String() {
			case "j", "down":
				if m.sidebarCursor < 6 {
					m.sidebarCursor++
				}
				return m, nil
			case "k", "up":
				if m.sidebarCursor > 0 {
					m.sidebarCursor--
				}
				return m, nil
			case "enter":
				// Select the view (add 1 to skip ViewWelcome)
				newView := View(m.sidebarCursor + 1)
				m.currentView = newView
				m.sidebarMode = false

				// Initialize view if needed
				var cmd tea.Cmd
				if newView == ViewTasks {
					// Reload tasks when entering the view
					cmd = m.taskScreen.Init()
				}
				return m, cmd
			case "esc":
				m.sidebarMode = false
				return m, nil
			}
			return m, nil
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
		}
	}

	// Update current screen - ALWAYS pass messages to active view
	// This is important for async messages like tasksLoadedMsg
	var cmd tea.Cmd
	switch m.currentView {
	case ViewTasks:
		m.taskScreen, cmd = m.taskScreen.Update(msg)
	}
	return m, cmd
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
	case "s", "sidebar":
		// Enter sidebar navigation mode
		m.sidebarMode = true
		// Start at current view, but adjust for ViewWelcome offset
		if m.currentView == ViewWelcome {
			m.sidebarCursor = 0 // Start at Tasks
		} else {
			m.sidebarCursor = int(m.currentView) - 1 // -1 to adjust for ViewWelcome
		}
		return m, nil
	}

	return m, nil
}

// View renders the application
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Title bar with Nerd Font icons
	titleLeft := " UniCLI - Student Organization Manager"
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

	// Check if terminal is too small
	if m.width < 50 || m.height < 15 {
		minSizeMsg := lipgloss.NewStyle().
			Foreground(styles.Warning).
			Bold(true).
			Padding(2).
			Render("Terminal too small!\n\nMinimum size: 50x15\nCurrent: " +
				fmt.Sprintf("%dx%d", m.width, m.height))

		return lipgloss.JoinVertical(
			lipgloss.Left,
			titleBar,
			minSizeMsg,
			m.renderStatusBar(),
		)
	}

	// Main content area
	contentHeight := m.height - 4 // Reserve space for title and status bars

	// Responsive sidebar width
	sidebarWidth := 20
	if m.width < 80 {
		sidebarWidth = 15 // Smaller sidebar for narrow terminals
	}

	sidebar := m.renderSidebar(contentHeight, sidebarWidth)

	// Main content - account for borders (4 chars: 2 for sidebar border, 2 for content border)
	contentWidth := m.width - sidebarWidth - 4
	if contentWidth < 10 {
		contentWidth = 10
	}

	// Main content
	var content string
	switch m.currentView {
	case ViewWelcome:
		content = m.renderWelcome()
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

	// Create content panel with conditional border color
	contentPanelStyle := styles.Panel.
		Width(contentWidth).
		Height(contentHeight)

	// Highlight content panel border when NOT in sidebar mode
	if !m.sidebarMode {
		contentPanelStyle = contentPanelStyle.BorderForeground(styles.Primary)
	}

	contentPanel := contentPanelStyle.Render(content)

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
		icon  string
		label string
	}{
		{ViewTasks, "1", "", "Tasks"},
		{ViewCalendar, "2", "󰃭", "Calendar"},
		{ViewClasses, "3", "󱉟", "Classes"},
		{ViewGrades, "4", "", "Grades"},
		{ViewNotes, "5", "󰷈", "Notes"},
		{ViewStats, "6", "󰄨", "Stats"},
		{ViewSettings, "7", "", "Settings"},
	}

	var items []string
	for i, v := range views {
		style := lipgloss.NewStyle().Padding(0, 1)

		// If in sidebar mode, highlight cursor position
		if m.sidebarMode && i == m.sidebarCursor {
			style = style.
				Background(styles.Secondary).
				Foreground(styles.Background).
				Bold(true)
		} else if v.view == m.currentView {
			// Otherwise highlight current view
			style = style.
				Background(styles.Primary).
				Foreground(styles.Background).
				Bold(true)
		}

		item := style.Render(fmt.Sprintf("%s [%s] %s", v.icon, v.key, v.label))
		items = append(items, item)
	}

	sidebarContent := lipgloss.JoinVertical(lipgloss.Left, items...)

	// Choose border color based on sidebar mode
	sidebarPanel := styles.Panel.
		Width(width).
		Height(height)

	if m.sidebarMode {
		// Highlight border when in sidebar mode
		sidebarPanel = sidebarPanel.
			BorderForeground(styles.Secondary)
	}

	return sidebarPanel.Render(sidebarContent)
}

// renderStatusBar renders the bottom status bar
func (m Model) renderStatusBar() string {
	// Terminal size indicator (right side)
	terminalSize := styles.Dimmed.Render(fmt.Sprintf("%dx%d", m.width, m.height))

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

		leftContent := commandPrompt + cursor
		spacing := m.width - lipgloss.Width(leftContent) - lipgloss.Width(terminalSize) - 2
		if spacing < 0 {
			spacing = 0
		}

		statusLine := lipgloss.JoinHorizontal(
			lipgloss.Top,
			leftContent,
			strings.Repeat(" ", spacing),
			terminalSize,
		)

		return styles.StatusBar.
			Width(m.width).
			Render(statusLine)
	}

	// If in sidebar mode, show navigation help
	if m.sidebarMode {
		leftContent := styles.Dimmed.Render("SIDEBAR: j/k to navigate  |  Enter to select  |  Esc to exit")
		spacing := m.width - lipgloss.Width(leftContent) - lipgloss.Width(terminalSize) - 2
		if spacing < 0 {
			spacing = 0
		}

		statusLine := lipgloss.JoinHorizontal(
			lipgloss.Top,
			leftContent,
			strings.Repeat(" ", spacing),
			terminalSize,
		)

		return styles.StatusBar.
			Width(m.width).
			Render(statusLine)
	}

	// Normal mode status bar
	leftContent := styles.Dimmed.Render("[:s] Sidebar  |  [:h] Help  |  [:q] Quit")
	spacing := m.width - lipgloss.Width(leftContent) - lipgloss.Width(terminalSize) - 2
	if spacing < 0 {
		spacing = 0
	}

	statusLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftContent,
		strings.Repeat(" ", spacing),
		terminalSize,
	)

	return styles.StatusBar.
		Width(m.width).
		Render(statusLine)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// renderWelcome renders the welcome screen
func (m Model) renderWelcome() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		Render(" Welcome to UniCLI")

	subtitle := lipgloss.NewStyle().
		Foreground(styles.Foreground).
		Render("Your terminal-based student organization manager")

	gettingStarted := lipgloss.NewStyle().
		Foreground(styles.Secondary).
		Render("  Getting Started:")

	availableViews := lipgloss.NewStyle().
		Foreground(styles.Info).
		Render("  Available Views:")

	footer := lipgloss.NewStyle().
		Foreground(styles.Muted).
		Italic(true).
		Render("Press :s to get started!")

	// Left side content with Nerd Font icons
	leftLines := []string{
		"",
		title,
		"",
		subtitle,
		"",
		"",
		gettingStarted,
		"",
		styles.Dimmed.Render("  • Press :s to open sidebar"),
		styles.Dimmed.Render("  • Press :h for help"),
		styles.Dimmed.Render("  • Press :q to quit"),
		"",
		"",
		availableViews,
		"",
		styles.Dimmed.Render("   Tasks      - Manage tasks (Kanban board)"),
		styles.Dimmed.Render("  󰃭 Calendar   - View schedule"),
		styles.Dimmed.Render("  󱉟 Classes    - Organize classes"),
		styles.Dimmed.Render("   Grades     - Track grades"),
		styles.Dimmed.Render("  󰷈 Notes      - Quick notes"),
		styles.Dimmed.Render("  󰄨 Stats      - Productivity stats"),
		"",
		"",
		footer,
	}

	leftContent := lipgloss.NewStyle().
		Width(45).
		Render(strings.Join(leftLines, "\n"))

	var combined string

	// Only show ASCII art if terminal is wide enough (>= 100 columns)
	if m.width >= 100 {
		// ASCII art from file
		asciiArt := []string{
			"       ⢀⣴⡾⠃⠄⠄⠄⠄⠄⠈⠺⠟⠛⠛⠛⠛⠻⢿⣿⣿⣿⣿⣶⣤⡀ ",
			"     ⢀⣴⣿⡿⠁⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⣸⣿⣿⣿⣿⣿⣿⣿⣷ ",
			"    ⣴⣿⡿⡟⡼⢹⣷⢲⡶⣖⣾⣶⢄⠄⠄⠄⠄⠄⢀⣼⣿⢿⣿⣿⣿⣿⣿⣿⣿",
			"   ⣾⣿⡟⣾⡸⢠⡿⢳⡿⠍⣼⣿⢏⣿⣷⢄⡀⠄⢠⣾⢻⣿⣸⣿⣿⣿⣿⣿⣿⣿",
			" ⣡⣿⣿⡟⡼⡁⠁⣰⠂⡾⠉⢨⣿⠃⣿⡿⠍⣾⣟⢤⣿⢇⣿⢇⣿⣿⢿⣿⣿⣿⣿⣿",
			"⣱⣿⣿⡟⡐⣰⣧⡷⣿⣴⣧⣤⣼⣯⢸⡿⠁⣰⠟⢀⣼⠏⣲⠏⢸⣿⡟⣿⣿⣿⣿⣿⣿",
			"⣿⣿⡟⠁⠄⠟⣁⠄⢡⣿⣿⣿⣿⣿⣿⣦⣼⢟⢀⡼⠃⡹⠃⡀⢸⡿⢸⣿⣿⣿⣿⣿⡟",
			"⣿⣿⠃⠄⢀⣾⠋⠓⢰⣿⣿⣿⣿⣿⣿⠿⣿⣿⣾⣅⢔⣕⡇⡇⡼⢁⣿⣿⣿⣿⣿⣿⢣",
			"⣿⡟⠄⠄⣾⣇⠷⣢⣿⣿⣿⣿⣿⣿⣿⣭⣀⡈⠙⢿⣿⣿⡇⡧⢁⣾⣿⣿⣿⣿⣿⢏⣾",
			"⣿⡇⠄⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠟⢻⠇⠄⠄⢿⣿⡇⢡⣾⣿⣿⣿⣿⣿⣏⣼⣿",
			"⣿⣷⢰⣿⣿⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⢰⣧⣀⡄⢀⠘⡿⣰⣿⣿⣿⣿⣿⣿⠟⣼⣿⣿",
			"⢹⣿⢸⣿⣿⠟⠻⢿⣿⣿⣿⣿⣿⣿⣿⣶⣭⣉⣤⣿⢈⣼⣿⣿⣿⣿⣿⣿⠏⣾⣹⣿⣿",
			"⢸⠇⡜⣿⡟⠄⠄⠄⠈⠙⣿⣿⣿⣿⣿⣿⣿⣿⠟⣱⣻⣿⣿⣿⣿⣿⠟⠁⢳⠃⣿⣿⣿",
			" ⣰⡗⠹⣿⣄⠄⠄⠄⢀⣿⣿⣿⣿⣿⣿⠟⣅⣥⣿⣿⣿⣿⠿⠋  ⣾⡌⢠⣿⡿⠃",
			"⠜⠋⢠⣷⢻⣿⣿⣶⣾⣿⣿⣿⣿⠿⣛⣥⣾⣿⠿⠟⠛⠉           ",
		}

		rightContent := lipgloss.NewStyle().
			Foreground(styles.Primary).
			Padding(2, 0).
			Render(strings.Join(asciiArt, "\n"))

		// Combine left and right side by side
		combined = lipgloss.JoinHorizontal(
			lipgloss.Top,
			leftContent,
			rightContent,
		)
	} else {
		// For narrow terminals, just show the text content
		combined = leftContent
	}

	welcome := lipgloss.NewStyle().
		Padding(2, 4).
		Render(combined)

	return welcome
}
