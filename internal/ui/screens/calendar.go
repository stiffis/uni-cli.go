package screens

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stiffis/UniCLI/internal/database"
	"github.com/stiffis/UniCLI/internal/models"
	"github.com/stiffis/UniCLI/internal/ui/styles"
)

// CalendarScreen is the model for the calendar view
type CalendarScreen struct {
	db *database.DB
	currentDate time.Time
	width       int
	height      int
	selectedDay   int
	calendarItems []models.CalendarItem
}

// NewCalendarScreen creates a new model for the calendar view
func NewCalendarScreen(db *database.DB) tea.Model {
	now := time.Now()
	return CalendarScreen{
		db:          db,
		currentDate: now,
		selectedDay: now.Day(),
	}
}

// Init initializes the calendar screen
func (m CalendarScreen) Init() tea.Cmd {
	return m.fetchCalendarItemsCmd()
}

// fetchCalendarItemsCmd is a tea.Cmd that fetches tasks and events for the current month
func (m CalendarScreen) fetchCalendarItemsCmd() tea.Cmd {
	return func() tea.Msg {
		year, month := m.currentDate.Year(), m.currentDate.Month()

		// Fetch tasks
		tasks, err := m.db.Tasks().FindAll() // Assuming FindAll can be filtered by month later
		if err != nil {
			return errMsg{err}
		}

		// Filter tasks by due date within the current month
		var monthTasks []models.Task
		for _, task := range tasks {
			if task.DueDate != nil && task.DueDate.Year() == year && task.DueDate.Month() == month {
				monthTasks = append(monthTasks, task)
			}
		}

		// Fetch events
		events, err := m.db.Events().GetEventsByMonth(year, month)
		if err != nil {
			return errMsg{err}
		}

		var items []models.CalendarItem
		for _, task := range monthTasks {
			items = append(items, &task)
		}
		for _, event := range events {
			items = append(items, &event)
		}

		return calendarItemsFetchedMsg(items)
	}
}

// calendarItemsFetchedMsg is a message sent when calendar items are fetched
type calendarItemsFetchedMsg []models.CalendarItem

// errMsg is a message for errors
type errMsg struct {
	err error
}

func (e errMsg) Error() string { return e.err.Error() }



// Update handles messages and updates the calendar screen model
func (m CalendarScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		oldMonth := m.currentDate.Month()
		switch msg.String() {
		case "H":
			m.currentDate = m.currentDate.AddDate(0, -1, 0) // Go back one month
		case "L":
			m.currentDate = m.currentDate.AddDate(0, 1, 0) // Go forward one month
		case "h":
			m.selectedDay = max(1, m.selectedDay-1) // Move left one day
		case "l":
			lastOfMonth := time.Date(m.currentDate.Year(), m.currentDate.Month(), 1, 0, 0, 0, 0, m.currentDate.Location()).AddDate(0, 1, -1).Day()
			m.selectedDay = min(lastOfMonth, m.selectedDay+1) // Move right one day
		case "k":
			m.selectedDay = max(1, m.selectedDay-7) // Move up one week
		case "j":
			lastOfMonth := time.Date(m.currentDate.Year(), m.currentDate.Month(), 1, 0, 0, 0, 0, m.currentDate.Location()).AddDate(0, 1, -1).Day()
			m.selectedDay = min(lastOfMonth, m.selectedDay+7) // Move down one week
		case "esc":
			// Handle escape key if needed for this screen
		}
		if m.currentDate.Month() != oldMonth {
			// If month changed, ensure selectedDay is valid for the new month
			lastOfMonth := time.Date(m.currentDate.Year(), m.currentDate.Month(), 1, 0, 0, 0, 0, m.currentDate.Location()).AddDate(0, 1, -1).Day()
			m.selectedDay = min(m.selectedDay, lastOfMonth)
			return m, m.fetchCalendarItemsCmd() // Re-fetch items if month changed
		}
	case calendarItemsFetchedMsg:
		m.calendarItems = msg
		return m, nil
	case errMsg:
		// Handle error, e.g., display an error message
		return m, tea.Quit // For now, just quit on error
	}
	return m, nil
}

// View renders the calendar screen
func (m CalendarScreen) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing calendar..."
	}

	// Use a portion of the full width for the calendar
	containerWidth := m.width - 18 
	if containerWidth < 40 { // Minimum width for the calendar grid
		containerWidth = 40
	}

	// Header for month and year
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		Align(lipgloss.Center).
		Width(containerWidth).
		Render(m.currentDate.Format("January 2006"))

	// Days of the week header
	weekdays := []string{"Lun", "Mar", "Mié", "Jue", "Vie", "Sáb", "Dom"}
	dayWidth := containerWidth / 7
	if dayWidth < 4 { // Minimum width for each day cell
		dayWidth = 4
	}
	weekdayStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Accent).
		Align(lipgloss.Center).
		Width(dayWidth)

	var weekdayHeaders []string
	for _, day := range weekdays {
		weekdayHeaders = append(weekdayHeaders, weekdayStyle.Render(day))
	}
	weekdayHeader := lipgloss.JoinHorizontal(lipgloss.Top, weekdayHeaders...)

	// Calendar grid
	var rows []string
	var row []string

	firstOfMonth := time.Date(m.currentDate.Year(), m.currentDate.Month(), 1, 0, 0, 0, 0, m.currentDate.Location())
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	firstWeekday := (int(firstOfMonth.Weekday()) + 6) % 7

	// Create a cell style
	cellStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true).
		BorderForeground(styles.Muted).
		Width(dayWidth).
		Height(3). // 1 line for content + 2 for borders
		Align(lipgloss.Center)

	// Print leading empty cells
	for i := 0; i < firstWeekday; i++ {
		row = append(row, cellStyle.Render(""))
	}

	// Print days
	for day := 1; day <= lastOfMonth.Day(); day++ {
		dayStyle := cellStyle
		dayContent := fmt.Sprintf("%d", day)
		var iconContent string

		// Add icons for tasks and events
		var icons []string
		for _, item := range m.calendarItems {
			if item.GetStartTime().Day() == day && item.GetStartTime().Month() == m.currentDate.Month() {
				if item.GetType() == "task" {
					icons = append(icons, "")
				} else if item.GetType() == "event" {
					icons = append(icons, "")
				}
			}
		}
		if len(icons) > 0 {
			iconContent = lipgloss.NewStyle().Foreground(styles.AutumnRed).Render(strings.Join(icons, ""))
		}

		// Combine day number and icons
		cellRenderContent := lipgloss.JoinVertical(lipgloss.Center, dayContent, iconContent)

		// Highlight selected day
		isSelected := day == m.selectedDay
		isToday := time.Now().Year() == m.currentDate.Year() && time.Now().Month() == m.currentDate.Month() && day == time.Now().Day()

		        		if isSelected {
		        			dayStyle = dayStyle.Copy().BorderForeground(styles.Warning)
		        		} else if isToday {
		        			dayStyle = dayStyle.Copy().BorderForeground(styles.Secondary)
		        		}	
		row = append(row, dayStyle.Render(cellRenderContent))

		if (firstWeekday+day)%7 == 0 {
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
			row = []string{}
		}
	}

	// Print trailing empty cells
	if len(row) > 0 {
		for i := len(row); i < 7; i++ {
			row = append(row, cellStyle.Render(""))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
	}

	calendarGrid := lipgloss.JoinVertical(lipgloss.Left, rows...)

	shortcuts := renderShortcuts()

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		weekdayHeader,
		calendarGrid,
		shortcuts,
	)
}

func renderShortcuts() string {
	shortcuts := []string{
		styles.Shortcut.Render("h/j/k/l") + styles.ShortcutText.Render(" navigate"),
		styles.Shortcut.Render("H/L") + styles.ShortcutText.Render(" change month"),
	}

	shortcutLine := strings.Join(shortcuts, "  ")

	return lipgloss.NewStyle().Padding(1, 0).Render(shortcutLine)
}
