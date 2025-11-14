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

// WeekView represents the weekly time-blocking view
type WeekView struct {
	db                  *database.DB
	currentWeek         time.Time // Monday of current week
	width               int
	height              int
	events              []models.Event
	categories          []models.Category
	selectedDay         int // 0-6 (Mon-Sun)
	selectedHour        int // 0-23
	selectedMinute      int // 0 or 30
	showEventForm       bool
	eventForm           components.EventForm
	selectedEventID     string
	showDeleteConfirm   bool
	showCategoryManager bool
	categoryManager     *components.CategoryManager
	startHour           int // First hour to display (default 6)
	endHour             int // Last hour to display (default 22)
	scrollOffset        int // For vertical scrolling
	err                 error
	errorMessage        string
}

// NewWeekView creates a new week view
func NewWeekView(db *database.DB, currentDate time.Time) *WeekView {
	monday := getMonday(currentDate)
	
	return &WeekView{
		db:              db,
		currentWeek:     monday,
		selectedDay:     0,
		selectedHour:    9,  // Start at 9 AM
		selectedMinute:  0,  // Start at :00
		categoryManager: components.NewCategoryManager(db),
		startHour:       6,  // 6 AM
		endHour:         22, // 10 PM
		scrollOffset:    0,
	}
}

// getMonday returns the Monday of the week containing the given date
func getMonday(date time.Time) time.Time {
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday is 7, not 0
	}
	daysToSubtract := weekday - 1 // Monday is 1
	monday := date.AddDate(0, 0, -daysToSubtract)
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())
}

// Init initializes the week view
func (w *WeekView) Init() tea.Cmd {
	return tea.Batch(w.fetchWeekEvents(), w.fetchCategoriesCmd())
}

// fetchWeekEvents fetches events for the current week
func (w *WeekView) fetchWeekEvents() tea.Cmd {
	return func() tea.Msg {
		startOfWeek := w.currentWeek

		weekEvents, err := w.db.Events().GetEventsWithCoursesForWeek(startOfWeek, w.db.Courses())
		if err != nil {
			return errMsg{err}
		}

		return weekEventsFetchedMsg(weekEvents)
	}
}

func (w *WeekView) fetchCategoriesCmd() tea.Cmd {
	return func() tea.Msg {
		categories, err := w.db.Categories().FindAll()
		if err != nil {
			return errMsg{err}
		}
		return categoriesFetchedMsg(categories)
	}
}

type weekEventsFetchedMsg []models.Event

func (w *WeekView) Update(msg tea.Msg) (*WeekView, tea.Cmd) {
	var cmd tea.Cmd

	if w.showCategoryManager {
		newModel, newCmd := w.categoryManager.Update(msg)
		w.categoryManager = newModel.(*components.CategoryManager)
		if w.categoryManager.IsQuitting() {
			w.showCategoryManager = false
			return w, w.fetchCategoriesCmd()
		}
		return w, newCmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w.width = msg.Width
		w.height = msg.Height

	case tea.KeyMsg:
		if w.showDeleteConfirm {
			switch msg.String() {
			case "y", "Y":
				w.showDeleteConfirm = false
				if w.selectedEventID != "" {
					return w, w.deleteEvent(w.selectedEventID)
				}
			case "n", "N", "esc":
				w.showDeleteConfirm = false
			}
			return w, nil
		}

		if w.showEventForm {
			w.eventForm, cmd = w.eventForm.Update(msg)
			if w.eventForm.IsSubmitted() {
				event := w.eventForm.GetEvent()
				w.showEventForm = false
				if w.eventForm.IsNewEvent() {
					return w, w.createEvent(event)
				} else {
					return w, w.updateEvent(event)
				}
			} else if w.eventForm.IsCancelled() {
				w.showEventForm = false
			}
			return w, cmd
		}

		// Navigation keys
		switch msg.String() {
		case "h", "left":
			if w.selectedDay > 0 {
				w.selectedDay--
			}
		case "l", "right":
			if w.selectedDay < 6 {
				w.selectedDay++
			}
		case "j", "down":
			// Move down by 30 minutes
			if w.selectedMinute == 0 {
				w.selectedMinute = 30
			} else {
				w.selectedMinute = 0
				if w.selectedHour < w.endHour {
					w.selectedHour++
				}
			}
		case "k", "up":
			// Move up by 30 minutes
			if w.selectedMinute == 30 {
				w.selectedMinute = 0
			} else {
				w.selectedMinute = 30
				if w.selectedHour > w.startHour {
					w.selectedHour--
				}
			}
		case "H":
			// Previous week
			w.currentWeek = w.currentWeek.AddDate(0, 0, -7)
			return w, w.fetchWeekEvents()
		case "L":
			// Next week
			w.currentWeek = w.currentWeek.AddDate(0, 0, 7)
			return w, w.fetchWeekEvents()
		case "n":
			// New event at selected slot
			w.showEventForm = true
			selectedDate := w.currentWeek.AddDate(0, 0, w.selectedDay)
			startTime := time.Date(selectedDate.Year(), selectedDate.Month(), selectedDate.Day(),
				w.selectedHour, w.selectedMinute, 0, 0, selectedDate.Location())
			event := &models.Event{
				ID:            "", // Empty ID means new event
				StartDatetime: startTime,
				Type:          "event",
				CreatedAt:     time.Now(),
			}
			endTime := startTime.Add(1 * time.Hour)
			event.EndDatetime = &endTime
			w.eventForm = components.NewEventForm(event, w.categories)
			return w, nil
		case "e":
			// Edit event at selected slot
			eventID := w.getEventAtSlot(w.selectedDay, w.selectedHour, w.selectedMinute)
			if eventID != "" {
				w.selectedEventID = eventID
				var eventToEdit *models.Event
				for i := range w.events {
					if w.events[i].ID == eventID {
						eventToEdit = &w.events[i]
						break
					}
				}
				if eventToEdit != nil {
					w.showEventForm = true
					w.eventForm = components.NewEventForm(eventToEdit, w.categories)
				}
			}
			return w, nil
		case "d":
			eventID := w.getEventAtSlot(w.selectedDay, w.selectedHour, w.selectedMinute)
			if eventID != "" {
				w.selectedEventID = eventID
				w.showDeleteConfirm = true
			}
			return w, nil
		case "c":
			w.showCategoryManager = true
			w.categoryManager.Reset()
			w.categoryManager.SetSize(w.width, w.height)
			return w, w.categoryManager.Init()
		}

	case weekEventsFetchedMsg:
		w.events = msg
		// Populate category for each event
		for i := range w.events {
			for _, category := range w.categories {
				if category.ID == w.events[i].CategoryID {
					w.events[i].Category = &category
					break
				}
			}
		}
		return w, nil

	case categoriesFetchedMsg:
		w.categories = msg
		return w, nil

	case errMsg:
		// Don't quit on error, just log it
		w.err = msg.err
		w.errorMessage = fmt.Sprintf("Error: %v", msg.err)
		return w, nil
	}

	return w, nil
}

// getEventAtSlot returns the ID of an event at the given day, hour and minute, or empty string
func (w *WeekView) getEventAtSlot(day, hour, minute int) string {
	selectedDate := w.currentWeek.AddDate(0, 0, day)
	slotMinute := hour*60 + minute
	
	for _, event := range w.events {
		if event.StartDatetime.Day() != selectedDate.Day() ||
			event.StartDatetime.Month() != selectedDate.Month() {
			continue
		}
		
		eventStartMinute := event.StartDatetime.Hour()*60 + event.StartDatetime.Minute()
		var eventEndMinute int
		if event.EndDatetime != nil {
			eventEndMinute = event.EndDatetime.Hour()*60 + event.EndDatetime.Minute()
		} else {
			eventEndMinute = eventStartMinute + 60
		}
		
		if slotMinute >= eventStartMinute && slotMinute < eventEndMinute {
			return event.ID
		}
	}
	return ""
}

func (w *WeekView) createEvent(event *models.Event) tea.Cmd {
	return func() tea.Msg {
		err := w.db.Events().Create(event)
		if err != nil {
			return errMsg{err}
		}
		return w.fetchWeekEvents()()
	}
}

func (w *WeekView) updateEvent(event *models.Event) tea.Cmd {
	return func() tea.Msg {
		err := w.db.Events().Update(event)
		if err != nil {
			return errMsg{err}
		}
		return w.fetchWeekEvents()()
	}
}

func (w *WeekView) deleteEvent(eventID string) tea.Cmd {
	return func() tea.Msg {
		err := w.db.Events().Delete(eventID)
		if err != nil {
			return errMsg{err}
		}
		return w.fetchWeekEvents()()
	}
}

func (w *WeekView) View() string {
	if w.width == 0 || w.height == 0 {
		return "Initializing week view..."
	}

	if w.showCategoryManager {
		return w.categoryManager.View()
	}

	var mainView string
	if w.showEventForm {
		mainView = w.eventForm.View()
	} else {
		mainView = w.renderWeekGrid()
	}

	if w.showDeleteConfirm {
		return w.renderDeleteConfirmDialog(mainView)
	}

	return mainView
}

func (w *WeekView) renderDeleteConfirmDialog(baseView string) string {
	var eventTitle string
	for _, event := range w.events {
		if event.ID == w.selectedEventID {
			eventTitle = event.Title
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

// renderWeekGrid renders the full week time-blocking grid
func (w *WeekView) renderWeekGrid() string {
	// Title
	weekStart := w.currentWeek
	weekEnd := weekStart.AddDate(0, 0, 6)
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		Align(lipgloss.Center).
		Width(w.width).
		Render(fmt.Sprintf("Week of %s - %s", 
			weekStart.Format("Jan 02"),
			weekEnd.Format("Jan 02, 2006")))

	// Calculate column widths
	timeColWidth := 6 // "HH:MM"
	availableWidth := w.width - timeColWidth - 4 // Extra padding
	dayColWidth := availableWidth / 7
	if dayColWidth < 8 {
		dayColWidth = 8
	}
	// Limit max width to avoid overflow
	if dayColWidth > 15 {
		dayColWidth = 15
	}

	// Header row with day names (shortened)
	weekdays := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	headerRow := w.renderHeaderRow(weekdays, timeColWidth, dayColWidth)

	// Time rows
	var rows []string
	visibleHours := w.endHour - w.startHour + 1
	maxVisibleRows := (w.height - 6) // Reserve space for title, header, and shortcuts
	
	startRow := w.startHour
	endRow := w.endHour
	
	if visibleHours > maxVisibleRows {
		// Need scrolling
		startRow = w.startHour + w.scrollOffset
		endRow = startRow + maxVisibleRows - 1
		if endRow > w.endHour {
			endRow = w.endHour
		}
	}

	now := time.Now()
	currentHour := now.Hour()
	currentMinute := now.Minute()
	
	for hour := startRow; hour <= endRow; hour++ {
		// Full hour (e.g., 09:00)
		rows = append(rows, w.renderTimeRow(hour, 0, dayColWidth))
		
		if hour == currentHour && currentMinute >= 0 && currentMinute < 30 {
			nowLine := w.renderNowLine(dayColWidth, currentHour, currentMinute)
			rows = append(rows, nowLine)
		}
		
		// Half hour (e.g., 09:30)
		rows = append(rows, w.renderTimeRow(hour, 30, dayColWidth))
		
		if hour == currentHour && currentMinute >= 30 {
			nowLine := w.renderNowLine(dayColWidth, currentHour, currentMinute)
			rows = append(rows, nowLine)
		}
	}

	grid := lipgloss.JoinVertical(lipgloss.Left, rows...)

	// Shortcuts
	shortcuts := w.renderShortcuts()

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		headerRow,
		grid,
		shortcuts,
	)
}

// renderHeaderRow renders the header with day names
func (w *WeekView) renderHeaderRow(weekdays []string, timeColWidth, dayColWidth int) string {
	// Time column header
	timeHeader := lipgloss.NewStyle().
		Width(timeColWidth).
		Align(lipgloss.Right).
		Bold(true).
		Foreground(styles.Accent).
		Render("")

	// Separator after time column
	separator := lipgloss.NewStyle().
		Foreground(styles.Border).
		Render("│")

	var dayHeaders []string
	for i, day := range weekdays {
		date := w.currentWeek.AddDate(0, 0, i)
		dayStr := fmt.Sprintf("%s%d", day, date.Day())
		
		style := lipgloss.NewStyle().
			Width(dayColWidth).
			Align(lipgloss.Center).
			Bold(true).
			Foreground(styles.Accent)

		// Highlight today
		today := time.Now()
		if date.Year() == today.Year() && date.YearDay() == today.YearDay() {
			style = style.Foreground(styles.Primary)
		}

		// Highlight selected day
		if i == w.selectedDay {
			style = style.Background(styles.Primary).Foreground(styles.Background)
		}

		dayHeaders = append(dayHeaders, style.Render(dayStr))
	}

	headerContent := lipgloss.JoinHorizontal(lipgloss.Top, timeHeader, separator, strings.Join(dayHeaders, ""))
	
	return lipgloss.NewStyle().
		Padding(0, 0, 1, 0).
		Render(headerContent)
}

// renderTimeRow renders a single time row across all days (hour and minute)
func (w *WeekView) renderTimeRow(hour, minute, dayColWidth int) string {
	// Time column - only show time label for full hours
	var timeStr string
	if minute == 0 {
		timeStr = fmt.Sprintf("%02d:00", hour)
	} else {
		timeStr = fmt.Sprintf("  :%02d", minute)
	}
	
	timeCell := lipgloss.NewStyle().
		Width(6).
		Align(lipgloss.Right).
		Foreground(styles.Muted).
		Render(timeStr)

	// Separator after time column
	separator := lipgloss.NewStyle().
		Foreground(styles.Border).
		Render("│")

	// Day columns
	var dayCells []string
	for day := 0; day < 7; day++ {
		cell := w.renderDayCell(day, hour, minute, dayColWidth)
		dayCells = append(dayCells, cell)
	}

	rowContent := lipgloss.JoinHorizontal(lipgloss.Top, timeCell, separator, strings.Join(dayCells, ""))

	return rowContent
}

// renderNowLine renders the current time indicator line
func (w *WeekView) renderNowLine(dayColWidth, currentHour, currentMinute int) string {
	nowStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E82424")).
		Bold(true)
	
	// Time indicator
	timeStr := fmt.Sprintf("▶%02d:%02d", currentHour, currentMinute)
	timeCell := lipgloss.NewStyle().
		Width(6).
		Align(lipgloss.Right).
		Foreground(lipgloss.Color("#E82424")).
		Bold(true).
		Render(timeStr)
	
	// Separator
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E82424")).
		Render("│")
	
	// Line across all days
	totalWidth := dayColWidth * 7
	nowLine := nowStyle.Render(strings.Repeat("─", totalWidth))
	
	return lipgloss.JoinHorizontal(lipgloss.Top, timeCell, separator, nowLine)
}

// renderDayCell renders a single cell (day x hour x minute)
func (w *WeekView) renderDayCell(day, hour, minute, width int) string {
	selectedDate := w.currentWeek.AddDate(0, 0, day)
	
	var cellContent string
	var cellStyle lipgloss.Style

	// Find event at this slot
	for _, event := range w.events {
		eventStart := event.StartDatetime
		eventEnd := event.EndDatetime
		
		eventStartMinute := eventStart.Hour()*60 + eventStart.Minute()
		slotMinute := hour*60 + minute
		
		var eventEndMinute int
		if eventEnd != nil {
			eventEndMinute = eventEnd.Hour()*60 + eventEnd.Minute()
		} else {
			eventEndMinute = eventStartMinute + 60 // Default 1 hour
		}
		
		// Event is in this slot if:
		// - Same day
		// - Slot time is between event start and end
		if eventStart.Day() == selectedDate.Day() &&
			eventStart.Month() == selectedDate.Month() &&
			slotMinute >= eventStartMinute && slotMinute < eventEndMinute {
			
			bgColor := styles.Info
			if event.Type == "class" && strings.HasPrefix(event.CategoryID, "course_") {
				// This is a course class, get color from course
				courseID := strings.TrimPrefix(event.CategoryID, "course_")
				course, err := w.db.Courses().GetByID(courseID)
				if err == nil && course.Color != "" {
					bgColor = lipgloss.Color(course.Color)
				}
			} else if event.Category != nil && event.Category.Color != "" {
				bgColor = lipgloss.Color(event.Category.Color)
			}

			if slotMinute == eventStartMinute || (minute == 0 && slotMinute > eventStartMinute && slotMinute < eventStartMinute+30) {
				// Truncate title if too long
				title := event.Title
				maxLen := width - 1
				if maxLen < 3 {
					maxLen = 3
				}
				if len(title) > maxLen {
					if maxLen > 3 {
						title = title[:maxLen-3] + "..."
					} else {
						title = title[:maxLen]
					}
				}
				cellContent = title
			} else {
				// Continuation of event, just show color
				cellContent = " "
			}

			cellStyle = lipgloss.NewStyle().
				Background(bgColor).
				Foreground(styles.Background).
				Width(width).
				Align(lipgloss.Center).
				Padding(0)
			
			break
		}
	}

	// If no event, render empty cell
	if cellContent == "" {
		cellStyle = lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Padding(0)
		cellContent = " "
	}

	// Highlight selected cell - use background instead of extra borders
	if day == w.selectedDay && hour == w.selectedHour && minute == w.selectedMinute {
		cellStyle = cellStyle.
			Background(styles.SelectedBackground).
			Foreground(styles.SelectedForeground).
			BorderForeground(styles.Warning).
			Border(lipgloss.NormalBorder(), false, true, false, false)
	} else {
		cellStyle = cellStyle.
			BorderForeground(styles.Border).
			Border(lipgloss.NormalBorder(), false, true, false, false)
	}

	return cellStyle.Render(cellContent)
}

func (w *WeekView) renderShortcuts() string {
	shortcuts := []string{
		styles.Shortcut.Render("h/j/k/l") + styles.ShortcutText.Render(" navigate"),
		styles.Shortcut.Render("H/L") + styles.ShortcutText.Render(" change week"),
		styles.Shortcut.Render("n") + styles.ShortcutText.Render(" new event"),
		styles.Shortcut.Render("e") + styles.ShortcutText.Render(" edit"),
		styles.Shortcut.Render("d") + styles.ShortcutText.Render(" delete"),
		styles.Shortcut.Render("c") + styles.ShortcutText.Render(" categories"),
		styles.Shortcut.Render("esc") + styles.ShortcutText.Render(" back to month"),
	}

	shortcutLine := strings.Join(shortcuts, "  ")

	return lipgloss.NewStyle().
		Padding(1, 0).
		BorderTop(true).
		BorderForeground(styles.Border).
		Render(shortcutLine)
}
