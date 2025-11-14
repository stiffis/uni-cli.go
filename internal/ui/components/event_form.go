package components

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stiffis/UniCLI/internal/models"
	"github.com/stiffis/UniCLI/internal/ui/styles"
)

// EventForm is a form for creating/editing events
type EventForm struct {
	originalEvent        *models.Event
	eventID              string // ID of the event being edited (empty if new event)
	titleInput           Input
	descriptionInput     TextArea
	startDateTimeInput   Input
	endDateTimeInput     Input
	recurrenceRuleInput  Input
	recurrenceEndDateInput Input

	// Focus tracking
	focusedField int
	submitted    bool
	cancelled    bool

	width  int
	height int
}

const (
	eventFieldTitle = iota
	eventFieldDescription
	eventFieldStartDateTime
	eventFieldEndDateTime
	eventFieldRecurrenceRule
	eventFieldRecurrenceEndDate
	eventFieldButtons
)

// NewEventForm creates a new event form, optionally pre-filling with existing event data
func NewEventForm(event *models.Event) EventForm {
	titleInput := NewInput("Title:", "Enter event title...")
	descriptionInput := NewTextArea("Description:", "Enter event description...")
	startDateTimeInput := NewInput("Start Time:", "YYYY-MM-DD HH:MM")
	endDateTimeInput := NewInput("End Time (optional):", "YYYY-MM-DD HH:MM")
	recurrenceRuleInput := NewInput("Recurrence:", "none, daily, weekly, monthly")
	recurrenceEndDateInput := NewInput("Recurrence End Date:", "YYYY-MM-DD")

	form := EventForm{
		titleInput:           titleInput,
		descriptionInput:     descriptionInput,
		startDateTimeInput:   startDateTimeInput,
		endDateTimeInput:     endDateTimeInput,
		recurrenceRuleInput:  recurrenceRuleInput,
		recurrenceEndDateInput: recurrenceEndDateInput,
		focusedField:         eventFieldTitle,
		width:                60,
		height:               20,
	}

	// If an event is provided, pre-fill the form fields
	if event != nil {
		form.originalEvent = event
		form.eventID = event.ID
		form.titleInput.SetValue(event.Title)
		form.descriptionInput.SetValue(event.Description)
		form.startDateTimeInput.SetValue(event.StartDatetime.Format("2006-01-02 15:04"))
		if event.EndDatetime != nil {
			form.endDateTimeInput.SetValue(event.EndDatetime.Format("2006-01-02 15:04"))
		}
		form.recurrenceRuleInput.SetValue(event.RecurrenceRule)
		if event.RecurrenceEndDate != nil {
			form.recurrenceEndDateInput.SetValue(event.RecurrenceEndDate.Format("2006-01-02"))
		}
	}

	// Focus first field
	form.titleInput.Focus()

	return form
}

// Init initializes the form
func (f EventForm) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (f EventForm) Update(msg tea.Msg) (EventForm, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Cancel form
			f.cancelled = true
			return f, nil

		case "tab", "down":
			// Move to next field
			f.blurAll()
			f.focusedField = (f.focusedField + 1) % 7
			cmd = f.focusField(f.focusedField)
			return f, cmd

		case "shift+tab", "up":
			// Move to previous field
			f.blurAll()
			f.focusedField = (f.focusedField + 6) % 7
			cmd = f.focusField(f.focusedField)
			return f, cmd

		case "enter":
			if f.focusedField == eventFieldButtons {
				// Submit form
				titleVal := f.titleInput.Value()
				if titleVal != "" {
					f.submitted = true
				}
				return f, nil
			}
		}
	}

	// Update focused field
	switch f.focusedField {
	case eventFieldTitle:
		cmd = f.titleInput.Update(msg)
	case eventFieldDescription:
		cmd = f.descriptionInput.Update(msg)
	case eventFieldStartDateTime:
		cmd = f.startDateTimeInput.Update(msg)
	case eventFieldEndDateTime:
		cmd = f.endDateTimeInput.Update(msg)
	case eventFieldRecurrenceRule:
		cmd = f.recurrenceRuleInput.Update(msg)
	case eventFieldRecurrenceEndDate:
		cmd = f.recurrenceEndDateInput.Update(msg)
	}

	return f, cmd
}

// View renders the form
func (f EventForm) View() string {
	var sections []string

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		Align(lipgloss.Center).
		Width(f.width).
		Render(" New Event")
	sections = append(sections, title)
	sections = append(sections, "")

	// Title input
	sections = append(sections, f.titleInput.View())
	sections = append(sections, "")

	// Description input
	sections = append(sections, f.descriptionInput.View())
	sections = append(sections, "")

	// Start date/time input
	sections = append(sections, f.startDateTimeInput.View())
	sections = append(sections, "")

	// End date/time input
	sections = append(sections, f.endDateTimeInput.View())
	sections = append(sections, "")

	// Recurrence rule input
	sections = append(sections, f.recurrenceRuleInput.View())
	sections = append(sections, "")

	// Recurrence end date input
	sections = append(sections, f.recurrenceEndDateInput.View())
	sections = append(sections, "")

	// Buttons
	sections = append(sections, f.renderButtons())
	sections = append(sections, "")

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(styles.Muted).
		Italic(true)

	help := helpStyle.Render("Tab: next field  |  Esc: cancel  |  Enter: submit")
	sections = append(sections, help)

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Create modal box
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Primary).
		Padding(1, 2).
		Width(f.width)

	return modalStyle.Render(content)
}

// renderButtons renders the action buttons
func (f EventForm) renderButtons() string {
	var submitText string
	if f.eventID != "" {
		submitText = "[ Save ]"
	} else {
		submitText = "[ Create ]"
	}

	submitStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Foreground(styles.Success)

	cancelStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Foreground(styles.Muted)

	if f.focusedField == eventFieldButtons {
		submitStyle = submitStyle.
			Background(styles.Success).
			Foreground(styles.Background).
			Bold(true)
	}

	submit := submitStyle.Render(submitText)
	cancel := cancelStyle.Render("[ Cancel (Esc) ]")

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		submit,
		"  ",
		cancel,
	)
}

// blurAll removes focus from all fields
func (f *EventForm) blurAll() {
	f.titleInput.Blur()
	f.descriptionInput.Blur()
	f.startDateTimeInput.Blur()
	f.endDateTimeInput.Blur()
	f.recurrenceRuleInput.Blur()
	f.recurrenceEndDateInput.Blur()
}

// focusField focuses a specific field
func (f *EventForm) focusField(field int) tea.Cmd {
	switch field {
	case eventFieldTitle:
		return f.titleInput.Focus()
	case eventFieldDescription:
		return f.descriptionInput.Focus()
	case eventFieldStartDateTime:
		return f.startDateTimeInput.Focus()
	case eventFieldEndDateTime:
		return f.endDateTimeInput.Focus()
	case eventFieldRecurrenceRule:
		return f.recurrenceRuleInput.Focus()
	case eventFieldRecurrenceEndDate:
		return f.recurrenceEndDateInput.Focus()
	}
	return nil
}

// GetEvent returns the event from form data
func (f EventForm) GetEvent() *models.Event {
	var event *models.Event
	if f.originalEvent != nil {
		// Editing existing event
		event = f.originalEvent
	} else {
		// Creating new event
		event = models.NewEvent("", time.Now()) // Provide dummy title and start time
	}

	event.Title = f.titleInput.Value()
	event.Description = f.descriptionInput.Value()

	// Parse start date/time
	startDateTimeStr := strings.TrimSpace(f.startDateTimeInput.Value())
	if startDateTime, err := time.ParseInLocation("2006-01-02 15:04", startDateTimeStr, time.Local); err == nil {
		event.StartDatetime = startDateTime
	} else {
		// If parsing fails, set a default valid time and prevent submission
		event.StartDatetime = time.Now()
		f.submitted = false // Prevent submission if start date is invalid
	}

	// Parse end date/time if provided
	endDateTimeStr := strings.TrimSpace(f.endDateTimeInput.Value())
	if endDateTimeStr != "" {
		if endDateTime, err := time.ParseInLocation("2006-01-02 15:04", endDateTimeStr, time.Local); err == nil {
			event.EndDatetime = &endDateTime
		}
	} else {
		event.EndDatetime = nil
	}

	// Recurrence
	event.RecurrenceRule = f.recurrenceRuleInput.Value()
	recurrenceEndDateStr := strings.TrimSpace(f.recurrenceEndDateInput.Value())
	if recurrenceEndDateStr != "" {
		if recurrenceEndDate, err := time.ParseInLocation("2006-01-02", recurrenceEndDateStr, time.Local); err == nil {
			event.RecurrenceEndDate = &recurrenceEndDate
		}
	} else {
		event.RecurrenceEndDate = nil
	}

	return event
}

// IsSubmitted returns true if form was submitted
func (f EventForm) IsSubmitted() bool {
	return f.submitted
}

// IsCancelled returns true if form was cancelled
func (f EventForm) IsCancelled() bool {
	return f.cancelled
}

// IsNewEvent returns true if this is a new event (not editing existing)
func (f EventForm) IsNewEvent() bool {
	return f.eventID == ""
}
