package screens

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stiffis/UniCLI/internal/database"
	"github.com/stiffis/UniCLI/internal/models"
	"github.com/stiffis/UniCLI/internal/ui/components"
	"github.com/stiffis/UniCLI/internal/ui/styles"
)

// CalendarScreen is the model for the calendar view
type CalendarScreen struct {
	db                  *database.DB
	currentDate         time.Time
	width               int
	height              int
	selectedDay         int
	calendarItems       []models.CalendarItem
	categories          []models.Category
	showDayDetails      bool
	showEventForm       bool
	eventForm           components.EventForm
	selectedEventID     string
	showDeleteConfirm   bool
	selectedItemIndex   int
	showCategoryManager bool
	categoryManager     *components.CategoryManager
	showWeekView        bool
	weekView            *WeekView
	showDayView         bool
	dayView             *DayView
}

// NewCalendarScreen creates a new model for the calendar view
func NewCalendarScreen(db *database.DB) tea.Model {
	now := time.Now()
	return CalendarScreen{
		db:           db,
		currentDate:  now,
		selectedDay:  now.Day(),
		categoryManager: components.NewCategoryManager(db),
	}
}

// Init initializes the calendar screen
func (m CalendarScreen) Init() tea.Cmd {
	return tea.Batch(m.fetchCalendarItemsCmd(), m.fetchCategoriesCmd())
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
			e := event
			items = append(items, &e)
		}

		return calendarItemsFetchedMsg(items)
	}
}

func (m CalendarScreen) fetchCategoriesCmd() tea.Cmd {
	return func() tea.Msg {
		categories, err := m.db.Categories().FindAll()
		if err != nil {
			return errMsg{err}
		}
		return categoriesFetchedMsg(categories)
	}
}

// calendarItemsFetchedMsg is a message sent when calendar items are fetched
type calendarItemsFetchedMsg []models.CalendarItem
type categoriesFetchedMsg []models.Category

// errMsg is a message for errors
type errMsg struct {
	err error
}

func (e errMsg) Error() string { return e.err.Error() }

func (m CalendarScreen) deleteEvent(eventID string) tea.Cmd {
	return func() tea.Msg {
		err := m.db.Events().Delete(eventID)
		if err != nil {
			return errMsg{err}
		}
		return m.fetchCalendarItemsCmd()()
	}
}

func (m CalendarScreen) getItemsForSelectedDay() []models.CalendarItem {
	var itemsForSelectedDay []models.CalendarItem
	for _, item := range m.calendarItems {
		if item.GetStartTime().Day() == m.selectedDay && item.GetStartTime().Month() == m.currentDate.Month() {
			itemsForSelectedDay = append(itemsForSelectedDay, item)
		}
	}
	return itemsForSelectedDay
}

func (m CalendarScreen) IsEventFormActive() bool {
	return m.showEventForm
}

func (m CalendarScreen) IsWeekViewActive() bool {
	return m.showWeekView
}

func (m CalendarScreen) IsWeekViewEventFormActive() bool {
	if m.showWeekView && m.weekView != nil {
		return m.weekView.showEventForm
	}
	return false
}

// Update handles messages and updates the calendar screen model
func (m CalendarScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.showCategoryManager {
		newModel, newCmd := m.categoryManager.Update(msg)
		m.categoryManager = newModel.(*components.CategoryManager)
		if m.categoryManager.IsQuitting() {
			m.showCategoryManager = false
			return m, m.fetchCategoriesCmd()
		}
		return m, newCmd
	}

	if m.showWeekView {
		// Check for escape key to return to month view
		// BUT only if no form is active in week view
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() == "esc" {
				// Only exit to month view if no forms are active
				if !m.weekView.showEventForm && !m.weekView.showCategoryManager && !m.weekView.showDeleteConfirm {
					m.showWeekView = false
					return m, m.fetchCalendarItemsCmd()
				}
			}
		}
		
		// Pass all messages to week view
		var newWeekView *WeekView
		newWeekView, cmd = m.weekView.Update(msg)
		m.weekView = newWeekView
		
		return m, cmd
	}

	if m.showDayView {
		// Check for escape key to return to month view
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() == "esc" {
				m.showDayView = false
				return m, m.fetchCalendarItemsCmd()
			}
		}
		
		// Pass all messages to day view
		var newDayView *DayView
		newDayView, cmd = m.dayView.Update(msg)
		m.dayView = newDayView
		
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if m.showDeleteConfirm {
			switch msg.String() {
			case "y", "Y":
				m.showDeleteConfirm = false
				if m.selectedEventID != "" {
					return m, m.deleteEvent(m.selectedEventID)
				}
			case "n", "N", "esc":
				m.showDeleteConfirm = false
			}
			return m, nil
		}

		if m.showEventForm {
			m.eventForm, cmd = m.eventForm.Update(msg)
			if m.eventForm.IsSubmitted() {
				event := m.eventForm.GetEvent()
				m.showEventForm = false
				if m.eventForm.IsNewEvent() {
					return m, m.createEvent(event)
				} else {
					return m, m.updateEvent(event)
				}
			} else if m.eventForm.IsCancelled() {
				m.showEventForm = false
			}
			return m, cmd
		}

		if m.showDayDetails {
			items := m.getItemsForSelectedDay()
			switch msg.String() {
			case "esc":
				m.showDayDetails = false
				m.selectedItemIndex = 0
				m.selectedEventID = ""
			case "j", "down":
				if m.selectedItemIndex < len(items)-1 {
					m.selectedItemIndex++
				}
			case "k", "up":
				if m.selectedItemIndex > 0 {
					m.selectedItemIndex--
				}
			case "n":
				m.showEventForm = true
				m.eventForm = components.NewEventForm(nil, m.categories)
				return m, nil
			case "e":
				if m.selectedItemIndex >= 0 && m.selectedItemIndex < len(items) {
					m.selectedEventID = items[m.selectedItemIndex].GetID()
					var eventToEdit *models.Event
					for _, item := range m.calendarItems {
						if item.GetID() == m.selectedEventID {
							if event, ok := item.(*models.Event); ok {
								eventToEdit = event
								break
							}
						}
					}
					if eventToEdit != nil {
						m.showEventForm = true
						m.eventForm = components.NewEventForm(eventToEdit, m.categories)
					}
				}
				return m, nil
			case "d":
				if m.selectedItemIndex >= 0 && m.selectedItemIndex < len(items) {
					m.selectedEventID = items[m.selectedItemIndex].GetID()
					m.showDeleteConfirm = true
				}
				return m, nil
			}
			return m, nil
		}

		oldMonth := m.currentDate.Month()
		switch msg.String() {
		case "s":
			// Switch to week view
			m.showWeekView = true
			m.weekView = NewWeekView(m.db, m.currentDate)
			m.weekView.width = m.width
			m.weekView.height = m.height
			return m, m.weekView.Init()
		case "c":
			m.showCategoryManager = true
			m.categoryManager.Reset()
			m.categoryManager.SetSize(m.width, m.height)
			return m, m.categoryManager.Init()
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
		case "enter":
			// Open day view instead of simple details
			selectedDate := time.Date(m.currentDate.Year(), m.currentDate.Month(), m.selectedDay, 0, 0, 0, 0, m.currentDate.Location())
			m.showDayView = true
			m.dayView = NewDayView(m.db, selectedDate)
			m.dayView.width = m.width
			m.dayView.height = m.height
			return m, m.dayView.Init()
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
		// Populate category for each event
		for _, item := range m.calendarItems {
			if event, ok := item.(*models.Event); ok {
				for _, category := range m.categories {
					if category.ID == event.CategoryID {
						event.Category = &category
						break
					}
				}
			}
		}
		return m, nil
	case categoriesFetchedMsg:
		m.categories = msg
		return m, nil
	case errMsg:
		// Handle error, e.g., display an error message
		return m, tea.Quit // For now, just quit on error
	}
	return m, nil
}

func (m CalendarScreen) createEvent(event *models.Event) tea.Cmd {
	return func() tea.Msg {
		err := m.db.Events().Create(event)
		if err != nil {
			return errMsg{err}
		}
		return m.fetchCalendarItemsCmd()()
	}
}

func (m CalendarScreen) updateEvent(event *models.Event) tea.Cmd {
	return func() tea.Msg {
		err := m.db.Events().Update(event)
		if err != nil {
			return errMsg{err}
		}
		return m.fetchCalendarItemsCmd()()
	}
}

// View renders the calendar screen
func (m CalendarScreen) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing calendar..."
	}

	if m.showCategoryManager {
		return m.categoryManager.View()
	}

	if m.showWeekView {
		return m.weekView.View()
	}

	if m.showDayView {
		return m.dayView.View()
	}

	var mainView string
	if m.showEventForm {
		mainView = m.eventForm.View()
	} else if m.showDayDetails {
		mainView = m.renderDayDetails()
	} else {
		mainView = m.renderCalendar()
	}

	if m.showDeleteConfirm {
		return m.renderDeleteConfirmDialog(mainView)
	}

	return mainView
}

func (m CalendarScreen) renderDeleteConfirmDialog(baseView string) string {
	var eventTitle string
	for _, item := range m.calendarItems {
		if item.GetID() == m.selectedEventID {
			eventTitle = item.GetTitle()
			break
		}
	}

	question := fmt.Sprintf("Delete event \"%s\"?", eventTitle)

	dialog := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Danger).
		Padding(1, 2).
		Render(lipgloss.JoinVertical(
			lipgloss.Center,
			styles.Title.Render(question),
			"",
			styles.Dimmed.Render("This action cannot be undone."),
			"",
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				styles.Shortcut.Render("y")+styles.ShortcutText.Render(" delete"),
				"  ",
				styles.Shortcut.Render("n")+styles.ShortcutText.Render(" cancel"),
			),
		))

	return dialog
}

func (m CalendarScreen) renderCalendar() string {
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
	
	// Calculate cell width (accounting for borders)
	// Each cell has 2 border characters (left and right), so we need to account for that
	baseCellWidth := (containerWidth - 14) / 7 // 14 = 2 borders * 7 cells
	if baseCellWidth < 4 {
		baseCellWidth = 4
	}
	
	// Calendar grid
	var rows []string
	var row []string

	firstOfMonth := time.Date(m.currentDate.Year(), m.currentDate.Month(), 1, 0, 0, 0, 0, m.currentDate.Location())
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	firstWeekday := (int(firstOfMonth.Weekday()) + 6) % 7

	// Create a cell style with borders
	cellStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true).
		BorderForeground(styles.Muted).
		Width(baseCellWidth).
		Height(3). // 1 line for content + 2 for borders
		Align(lipgloss.Center)
	
	// Create weekday header style to match cell width INCLUDING borders
	// The header cells should have the same total width as calendar cells (content + borders)
	weekdayStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Accent).
		Align(lipgloss.Center).
		Width(baseCellWidth + 2) // +2 to match the total width of cells with borders

	var weekdayHeaders []string
	for _, day := range weekdays {
		weekdayHeaders = append(weekdayHeaders, weekdayStyle.Render(day))
	}
	weekdayHeader := lipgloss.JoinHorizontal(lipgloss.Top, weekdayHeaders...)

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
				var icon string
				var color lipgloss.Color
				if item.GetType() == "task" {
					icon = ""
					color = styles.Info
				} else if event, ok := item.(*models.Event); ok {
					icon = ""
					if event.Category != nil && event.Category.Color != "" {
						color = lipgloss.Color(event.Category.Color)
					} else {
						color = styles.SakuraPink
					}
				}
				icons = append(icons, lipgloss.NewStyle().Foreground(color).Render(icon))
			}
		}
		if len(icons) > 0 {
			iconContent = strings.Join(icons, " ")
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

	shortcuts := m.renderShortcuts()

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		weekdayHeader,
		calendarGrid,
		shortcuts,
	)
}

func (m CalendarScreen) renderDayDetails() string {
	itemsForSelectedDay := m.getItemsForSelectedDay()

	// Build details content
	var detailsContent string
	if len(itemsForSelectedDay) > 0 {
		var itemStrings []string
		for i, item := range itemsForSelectedDay {
			var icon string
			var itemString string
			if item.GetType() == "task" {
				icon = lipgloss.NewStyle().Foreground(styles.Info).Render("")
				itemString = fmt.Sprintf("%s %s", icon, item.GetTitle())
			} else if event, ok := item.(*models.Event); ok {
				var color lipgloss.Color
				if event.Category != nil && event.Category.Color != "" {
					color = lipgloss.Color(event.Category.Color)
				} else {
					color = styles.SakuraPink
				}
				icon = lipgloss.NewStyle().Foreground(color).Render("")
				itemString = fmt.Sprintf("%s %s (%s)", icon, item.GetTitle(), item.GetStartTime().Format("15:04"))
			}

			if i == m.selectedItemIndex {
				itemString = lipgloss.NewStyle().Background(styles.SelectedBackground).Foreground(styles.SelectedForeground).Render(itemString)
			}
			itemStrings = append(itemStrings, itemString)
		}
		detailsContent = lipgloss.JoinVertical(lipgloss.Left, itemStrings...)
	} else {
		detailsContent = "No items for this day."
	}

	// Style the details panel
	detailsPanelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Primary).
		Padding(1, 2).
		Width(m.width / 2).
		Height(m.height / 2)

	renderedDetails := detailsPanelStyle.Render(detailsContent)
	shortcuts := m.renderShortcuts()

	return lipgloss.JoinVertical(lipgloss.Left,
		renderedDetails,
		shortcuts,
	)
}

func (m CalendarScreen) renderShortcuts() string {
	var shortcuts []string

	if m.showDayDetails {
		shortcuts = []string{
			styles.Shortcut.Render("esc") + styles.ShortcutText.Render(" close details"),
			styles.Shortcut.Render("n") + styles.ShortcutText.Render(" new event"),
			styles.Shortcut.Render("e") + styles.ShortcutText.Render(" edit event"),
			styles.Shortcut.Render("d") + styles.ShortcutText.Render(" delete event"),
		}
	} else {
		shortcuts = []string{
			styles.Shortcut.Render("h/j/k/l") + styles.ShortcutText.Render(" navigate"),
			styles.Shortcut.Render("H/L") + styles.ShortcutText.Render(" change month"),
			styles.Shortcut.Render("enter") + styles.ShortcutText.Render(" view day details"),
			styles.Shortcut.Render("s") + styles.ShortcutText.Render(" week view"),
			styles.Shortcut.Render("c") + styles.ShortcutText.Render(" categories"),
			styles.Shortcut.Render("n") + styles.ShortcutText.Render(" new event"),
		}
	}

	shortcutLine := strings.Join(shortcuts, "  ")

	return lipgloss.NewStyle().Padding(1, 0).Render(shortcutLine)
}
