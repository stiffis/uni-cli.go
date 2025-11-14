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

// DayView represents the daily timeline view
type DayView struct {
	db                  *database.DB
	currentDate         time.Time
	width               int
	height              int
	events              []models.Event
	tasks               []models.Task
	categories          []models.Category
	selectedHour        int // 0-23
	selectedMinute      int // 0 or 30
	startHour           int // First hour to display (default 6)
	endHour             int // Last hour to display (default 22)
	scrollOffset        int
	err                 error
	errorMessage        string
	showEventForm       bool
	eventForm           components.EventForm
	selectedEventID     string
	showDeleteConfirm   bool
	showCategoryManager bool
	categoryManager     *components.CategoryManager
}

// NewDayView creates a new day view for the given date
func NewDayView(db *database.DB, date time.Time) *DayView {
	return &DayView{
		db:              db,
		currentDate:     date,
		selectedHour:    9,
		selectedMinute:  0,
		startHour:       6,
		endHour:         22,
		scrollOffset:    0,
		categoryManager: components.NewCategoryManager(db),
	}
}

// Init initializes the day view
func (d *DayView) Init() tea.Cmd {
	return tea.Batch(d.fetchDayEvents(), d.fetchDayTasks(), d.fetchCategories())
}

// fetchDayEvents fetches events for the current day
func (d *DayView) fetchDayEvents() tea.Cmd {
	return func() tea.Msg {
		// Fetch events for this specific day (including course classes)
		dayEvents, err := d.db.Events().GetEventsWithCoursesForDay(d.currentDate, d.db.Courses())
		if err != nil {
			return errMsg{err}
		}

		return dayEventsFetchedMsg(dayEvents)
	}
}

// fetchDayTasks fetches tasks that are due on the current day
func (d *DayView) fetchDayTasks() tea.Cmd {
	return func() tea.Msg {
		allTasks, err := d.db.Tasks().FindAll()
		if err != nil {
			return errMsg{err}
		}

		// Filter tasks due on this day
		var dayTasks []models.Task
		for _, task := range allTasks {
			if task.DueDate != nil &&
				task.DueDate.Year() == d.currentDate.Year() &&
				task.DueDate.Month() == d.currentDate.Month() &&
				task.DueDate.Day() == d.currentDate.Day() {
				dayTasks = append(dayTasks, task)
			}
		}

		return dayTasksFetchedMsg(dayTasks)
	}
}

// fetchCategories fetches all categories
func (d *DayView) fetchCategories() tea.Cmd {
	return func() tea.Msg {
		categories, err := d.db.Categories().FindAll()
		if err != nil {
			return errMsg{err}
		}
		return categoriesFetchedMsg(categories)
	}
}

type dayEventsFetchedMsg []models.Event
type dayTasksFetchedMsg []models.Task

// Update handles messages
func (d *DayView) Update(msg tea.Msg) (*DayView, tea.Cmd) {
	var cmd tea.Cmd
	
	// Handle delete confirmation dialog
	if d.showDeleteConfirm {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "y", "Y":
				d.showDeleteConfirm = false
				if d.selectedEventID != "" {
					return d, d.deleteEvent(d.selectedEventID)
				}
			case "n", "N", "esc":
				d.showDeleteConfirm = false
			}
		}
		return d, nil
	}
	
	// Handle event form
	if d.showEventForm {
		d.eventForm, cmd = d.eventForm.Update(msg)
		if d.eventForm.IsSubmitted() {
			event := d.eventForm.GetEvent()
			d.showEventForm = false
			if d.eventForm.IsNewEvent() {
				return d, d.createEvent(event)
			} else {
				return d, d.updateEvent(event)
			}
		} else if d.eventForm.IsCancelled() {
			d.showEventForm = false
		}
		return d, cmd
	}
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height

	case tea.KeyMsg:
		// Navigation keys
		switch msg.String() {
		case "n":
			// New event at selected slot
			d.showEventForm = true
			startTime := time.Date(d.currentDate.Year(), d.currentDate.Month(), d.currentDate.Day(),
				d.selectedHour, d.selectedMinute, 0, 0, d.currentDate.Location())
			event := &models.Event{
				ID:            "",
				StartDatetime: startTime,
				Type:          "event",
				CreatedAt:     time.Now(),
			}
			endTime := startTime.Add(1 * time.Hour)
			event.EndDatetime = &endTime
			d.eventForm = components.NewEventForm(event, d.categories)
			return d, nil
			
		case "e":
			// Edit event at selected slot
			eventID := d.getEventAtSlot(d.selectedHour, d.selectedMinute)
			if eventID != "" {
				d.selectedEventID = eventID
				var eventToEdit *models.Event
				for i := range d.events {
					if d.events[i].ID == eventID {
						eventToEdit = &d.events[i]
						break
					}
				}
				if eventToEdit != nil {
					d.showEventForm = true
					d.eventForm = components.NewEventForm(eventToEdit, d.categories)
				}
			}
			return d, nil
			
		case "d":
			// Delete event at selected slot
			eventID := d.getEventAtSlot(d.selectedHour, d.selectedMinute)
			if eventID != "" {
				d.selectedEventID = eventID
				d.showDeleteConfirm = true
			}
			return d, nil
			
		case "j", "down":
			// Move down by 30 minutes
			if d.selectedMinute == 0 {
				d.selectedMinute = 30
			} else {
				d.selectedMinute = 0
				if d.selectedHour < d.endHour {
					d.selectedHour++
				}
			}
		case "k", "up":
			// Move up by 30 minutes
			if d.selectedMinute == 30 {
				d.selectedMinute = 0
			} else {
				d.selectedMinute = 30
				if d.selectedHour > d.startHour {
					d.selectedHour--
				}
			}
		case "h", "left":
			// Previous day
			d.currentDate = d.currentDate.AddDate(0, 0, -1)
			return d, d.Init()
		case "l", "right":
			// Next day
			d.currentDate = d.currentDate.AddDate(0, 0, 1)
			return d, d.Init()
		}

	case dayEventsFetchedMsg:
		d.events = msg
		// Populate categories for events
		for i := range d.events {
			for _, category := range d.categories {
				if category.ID == d.events[i].CategoryID {
					d.events[i].Category = &category
					break
				}
			}
		}
		return d, nil

	case dayTasksFetchedMsg:
		d.tasks = msg
		return d, nil

	case categoriesFetchedMsg:
		d.categories = msg
		return d, nil

	case errMsg:
		d.err = msg.err
		d.errorMessage = fmt.Sprintf("Error: %v", msg.err)
		return d, nil
	}

	return d, nil
}

// getEventAtSlot returns the ID of an event at the given time slot
func (d *DayView) getEventAtSlot(hour, minute int) string {
	slotMinute := hour*60 + minute
	
	for _, event := range d.events {
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

func (d *DayView) createEvent(event *models.Event) tea.Cmd {
	return func() tea.Msg {
		err := d.db.Events().Create(event)
		if err != nil {
			return errMsg{err}
		}
		return d.fetchDayEvents()()
	}
}

func (d *DayView) updateEvent(event *models.Event) tea.Cmd {
	return func() tea.Msg {
		err := d.db.Events().Update(event)
		if err != nil {
			return errMsg{err}
		}
		return d.fetchDayEvents()()
	}
}

func (d *DayView) deleteEvent(eventID string) tea.Cmd {
	return func() tea.Msg {
		err := d.db.Events().Delete(eventID)
		if err != nil {
			return errMsg{err}
		}
		return d.fetchDayEvents()()
	}
}

// View renders the day view
func (d *DayView) View() string {
	if d.width == 0 || d.height == 0 {
		return "Initializing day view..."
	}
	
	// Show event form if active
	if d.showEventForm {
		return d.eventForm.View()
	}

	// Calculate panel widths - make it more compact
	rightPanelWidth := 22 // Fixed width for summary panel
	leftPanelWidth := d.width - rightPanelWidth - 3 // Timeline panel
	
	// Limit max width to avoid being too wide, but increase from 100 to 120
	maxTotalWidth := 120
	if d.width > maxTotalWidth {
		leftPanelWidth = maxTotalWidth - rightPanelWidth - 3
	}

	// Render panels
	timelinePanel := d.renderTimeline(leftPanelWidth)
	summaryPanel := d.renderSummaryAndTasks(rightPanelWidth)

	// Add space between panels
	spacer := strings.Repeat(" ", 3)

	// Join panels horizontally with spacer
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		timelinePanel,
		spacer,
		summaryPanel,
	)

	// Title
	title := d.renderTitle()

	// Shortcuts
	shortcuts := d.renderShortcuts()

	mainView := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
		shortcuts,
	)
	
	// Show delete confirmation dialog if active
	if d.showDeleteConfirm {
		return d.renderDeleteConfirmDialog(mainView)
	}
	
	return mainView
}

// renderDeleteConfirmDialog renders the delete confirmation dialog
func (d *DayView) renderDeleteConfirmDialog(mainView string) string {
	dialog := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Warning).
		Padding(1, 2).
		Render("Are you sure you want to delete this event?\n\n[y] Yes  [n] No")

	return lipgloss.Place(
		d.width,
		d.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, mainView, dialog),
	)
}

// renderTitle renders the title bar
func (d *DayView) renderTitle() string {
	titleStr := d.currentDate.Format("Monday, January 02, 2006")
	
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		Render(titleStr)
}

// renderTimeline renders the left timeline panel
func (d *DayView) renderTimeline(width int) string {
	var rows []string

	// Render hours with events
	now := time.Now()
	isToday := now.Year() == d.currentDate.Year() && 
	           now.Month() == d.currentDate.Month() && 
	           now.Day() == d.currentDate.Day()
	currentHour := now.Hour()
	currentMinute := now.Minute()
	
	for hour := d.startHour; hour <= d.endHour; hour++ {
		for _, minute := range []int{0, 30} {
			timeStr := fmt.Sprintf("%02d:%02d", hour, minute)
			
			// Highlight selected time slot
			timeStyle := lipgloss.NewStyle().
				Foreground(styles.Muted).
				Width(5).
				Align(lipgloss.Right)
			
			if hour == d.selectedHour && minute == d.selectedMinute {
				timeStyle = timeStyle.
					Foreground(styles.Primary).
					Bold(true)
			}
			
			separator := lipgloss.NewStyle().
				Foreground(styles.Border).
				Render("│")
			
			// Check if there's an event at this time slot
			eventContent := d.renderEventAtSlot(hour, minute, width-7)
			
			row := timeStyle.Render(timeStr) + separator + eventContent
			rows = append(rows, row)
			
			// Add NOW line if this is today and we're between time slots
			if isToday && hour == currentHour && minute == 0 && currentMinute >= 0 && currentMinute < 30 {
				nowLineStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#E82424")).
					Bold(true)
				nowLine := nowLineStyle.Render(fmt.Sprintf("      ▶ ───────── NOW: %02d:%02d ─────────", currentHour, currentMinute))
				rows = append(rows, nowLine)
			} else if isToday && hour == currentHour && minute == 30 && currentMinute >= 30 {
				nowLineStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#E82424")).
					Bold(true)
				nowLine := nowLineStyle.Render(fmt.Sprintf("      ▶ ───────── NOW: %02d:%02d ─────────", currentHour, currentMinute))
				rows = append(rows, nowLine)
			}
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)

	return content
}

// renderEventAtSlot renders the event at a specific time slot
func (d *DayView) renderEventAtSlot(hour, minute, width int) string {
	slotMinute := hour*60 + minute
	
	// Find event at this slot
	for _, event := range d.events {
		eventStartMinute := event.StartDatetime.Hour()*60 + event.StartDatetime.Minute()
		
		var eventEndMinute int
		if event.EndDatetime != nil {
			eventEndMinute = event.EndDatetime.Hour()*60 + event.EndDatetime.Minute()
		} else {
			eventEndMinute = eventStartMinute + 60 // Default 1 hour
		}
		
		// Check if this slot is within the event time
		if slotMinute >= eventStartMinute && slotMinute < eventEndMinute {
			// Get color - either from category or from course
			bgColor := styles.Info
			if event.Type == "class" && strings.HasPrefix(event.CategoryID, "course_") {
				// This is a course class, get color from course
				courseID := strings.TrimPrefix(event.CategoryID, "course_")
				course, err := d.db.Courses().GetByID(courseID)
				if err == nil && course.Color != "" {
					bgColor = lipgloss.Color(course.Color)
				}
			} else if event.Category != nil && event.Category.Color != "" {
				bgColor = lipgloss.Color(event.Category.Color)
			}
			
			// Only show title on the first slot
			var content string
			if slotMinute == eventStartMinute {
				title := event.Title
				// Truncate if too long
				if len(title) > width-2 {
					title = title[:width-5] + "..."
				}
				content = " " + title
			} else {
				// Continuation of event
				content = " "
			}
			
			// Render with background color
			return lipgloss.NewStyle().
				Background(bgColor).
				Foreground(styles.Background).
				Width(width).
				Render(content)
		}
	}
	
	// No event at this slot, render empty space
	return strings.Repeat(" ", width)
}

// renderSummaryAndTasks renders the right summary panel
func (d *DayView) renderSummaryAndTasks(width int) string {
	// Summary section
	summary := d.renderSummary(width)
	
	// Tasks section
	tasksSection := d.renderTasks(width)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		summary,
		"",
		tasksSection,
	)

	return content
}

// renderSummary renders the day summary stats
func (d *DayView) renderSummary(width int) string {
	// Calculate stats
	totalEvents := len(d.events)
	totalTasks := len(d.tasks)
	
	// Calculate busy time
	var busyMinutes int
	for _, event := range d.events {
		if event.EndDatetime != nil {
			duration := event.EndDatetime.Sub(event.StartDatetime)
			busyMinutes += int(duration.Minutes())
		} else {
			busyMinutes += 60 // Default 1 hour
		}
	}
	
	busyHours := busyMinutes / 60
	busyMins := busyMinutes % 60
	
	// Free time (assuming 16 hour workday: 6am - 10pm)
	workdayMinutes := 16 * 60
	freeMinutes := workdayMinutes - busyMinutes
	freeHours := freeMinutes / 60
	freeMins := freeMinutes % 60

	summaryTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Accent).
		Render("📊 DAY SUMMARY")

	stats := []string{
		fmt.Sprintf("🕐 %d events", totalEvents),
		fmt.Sprintf("⏱️  %dh %dm busy", busyHours, busyMins),
		fmt.Sprintf("⏳ %dh %dm free", freeHours, freeMins),
		fmt.Sprintf("🔴 %d tasks due", totalTasks),
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		summaryTitle,
		strings.Repeat("─", width-4),
		strings.Join(stats, "\n"),
	)
}

// renderTasks renders the tasks section
func (d *DayView) renderTasks(width int) string {
	tasksTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Accent).
		Render("📋 TASKS DUE")

	if len(d.tasks) == 0 {
		noTasks := lipgloss.NewStyle().
			Foreground(styles.Muted).
			Render("No tasks due today")
		
		return lipgloss.JoinVertical(
			lipgloss.Left,
			tasksTitle,
			strings.Repeat("─", width-4),
			noTasks,
		)
	}

	var taskLines []string
	for _, task := range d.tasks {
		// Checkbox
		checkbox := "☐"
		if task.Status == models.TaskStatusCompleted {
			checkbox = "☑"
		}
		
		// Priority indicator with emoji
		var priorityEmoji string
		switch task.Priority {
		case models.TaskPriorityUrgent:
			priorityEmoji = "☠️" // Skull for urgent
		case models.TaskPriorityHigh:
			priorityEmoji = "🔴" // Red circle for high
		case models.TaskPriorityMedium:
			priorityEmoji = "💛" // Yellow heart for medium
		case models.TaskPriorityLow:
			priorityEmoji = "🤍" // White heart for low
		default:
			priorityEmoji = "⚪" // White circle for none/default
		}
		
		// Truncate title if too long
		title := task.Title
		maxLen := width - 10
		if len(title) > maxLen {
			title = title[:maxLen-3] + "..."
		}
		
		taskLine := fmt.Sprintf("%s %s %s", checkbox, title, priorityEmoji)
		taskLines = append(taskLines, taskLine)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		tasksTitle,
		strings.Repeat("─", width-4),
		strings.Join(taskLines, "\n"),
	)
}

// renderShortcuts renders the shortcuts bar
func (d *DayView) renderShortcuts() string {
	shortcuts := []string{
		styles.Shortcut.Render("h/l") + styles.ShortcutText.Render(" prev/next day"),
		styles.Shortcut.Render("j/k") + styles.ShortcutText.Render(" navigate time"),
		styles.Shortcut.Render("n") + styles.ShortcutText.Render(" new event"),
		styles.Shortcut.Render("e") + styles.ShortcutText.Render(" edit"),
		styles.Shortcut.Render("d") + styles.ShortcutText.Render(" delete"),
		styles.Shortcut.Render("esc") + styles.ShortcutText.Render(" back"),
	}

	return strings.Join(shortcuts, "  ")
}
