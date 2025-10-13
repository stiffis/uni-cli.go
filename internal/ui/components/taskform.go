package components

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stiffis/UniCLI/internal/models"
	"github.com/stiffis/UniCLI/internal/ui/styles"
)

// TaskForm is a form for creating/editing tasks
type TaskForm struct {
	titleInput       Input
	descriptionInput TextArea
	dueDateInput     Input
	
	// Priority selector
	priorities       []models.TaskPriority
	selectedPriority int
	
	// Focus tracking
	focusedField int
	submitted    bool
	cancelled    bool
	
	width  int
	height int
}

const (
	fieldTitle = iota
	fieldDescription
	fieldDueDate
	fieldPriority
	fieldButtons
)

// NewTaskForm creates a new task form
func NewTaskForm() TaskForm {
	titleInput := NewInput("Title:", "Enter task title...")
	descriptionInput := NewTextArea("Description:", "Enter task description...")
	dueDateInput := NewInput("Due Date (optional):", "YYYY-MM-DD or leave empty")

	priorities := []models.TaskPriority{
		models.TaskPriorityLow,
		models.TaskPriorityMedium,
		models.TaskPriorityHigh,
		models.TaskPriorityUrgent,
	}

	form := TaskForm{
		titleInput:       titleInput,
		descriptionInput: descriptionInput,
		dueDateInput:     dueDateInput,
		priorities:       priorities,
		selectedPriority: 1, // Default to Medium
		focusedField:     fieldTitle,
		width:            60,
		height:           20,
	}

	// Focus first field
	form.titleInput.Focus()

	return form
}

// Init initializes the form
func (f TaskForm) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (f TaskForm) Update(msg tea.Msg) (TaskForm, tea.Cmd) {
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
			f.focusedField = (f.focusedField + 1) % 5
			cmd = f.focusField(f.focusedField)
			return f, cmd

		case "shift+tab", "up":
			// Move to previous field
			f.blurAll()
			f.focusedField = (f.focusedField + 4) % 5
			cmd = f.focusField(f.focusedField)
			return f, cmd

		case "left":
			if f.focusedField == fieldPriority {
				if f.selectedPriority > 0 {
					f.selectedPriority--
				}
				return f, nil
			}

		case "right":
			if f.focusedField == fieldPriority {
				if f.selectedPriority < len(f.priorities)-1 {
					f.selectedPriority++
				}
				return f, nil
			}

		case "enter":
			if f.focusedField == fieldButtons {
				// Submit form
				if f.titleInput.Value() != "" {
					f.submitted = true
				}
				return f, nil
			}
		}
	}

	// Update focused field
	switch f.focusedField {
	case fieldTitle:
		cmd = f.titleInput.Update(msg)
	case fieldDescription:
		cmd = f.descriptionInput.Update(msg)
	case fieldDueDate:
		cmd = f.dueDateInput.Update(msg)
	}

	return f, cmd
}

// View renders the form
func (f TaskForm) View() string {
	var sections []string

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		Align(lipgloss.Center).
		Width(f.width).
		Render(" New Task")
	sections = append(sections, title)
	sections = append(sections, "")

	// Title input
	sections = append(sections, f.titleInput.View())
	sections = append(sections, "")

	// Description input
	sections = append(sections, f.descriptionInput.View())
	sections = append(sections, "")

	// Due date input
	sections = append(sections, f.dueDateInput.View())
	sections = append(sections, "")

	// Priority selector
	sections = append(sections, f.renderPrioritySelector())
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

// renderPrioritySelector renders the priority selection
func (f TaskForm) renderPrioritySelector() string {
	labelStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true)

	label := labelStyle.Render("Priority:")

	var options []string
	priorityIcons := map[models.TaskPriority]string{
		models.TaskPriorityUrgent: "",
		models.TaskPriorityHigh:   "🔴",
		models.TaskPriorityMedium: "💛",
		models.TaskPriorityLow:    "🤍",
	}

	for i, priority := range f.priorities {
		icon := priorityIcons[priority]
		text := fmt.Sprintf("%s %s", icon, strings.Title(string(priority)))

		style := lipgloss.NewStyle().Padding(0, 2)

		if f.focusedField == fieldPriority && i == f.selectedPriority {
			// Selected and focused
			style = style.
				Background(styles.Primary).
				Foreground(styles.Background).
				Bold(true)
		} else if i == f.selectedPriority {
			// Just selected
			style = style.
				Foreground(styles.Primary).
				Bold(true)
		}

		options = append(options, style.Render(text))
	}

	optionsLine := lipgloss.JoinHorizontal(lipgloss.Top, options...)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		label,
		optionsLine,
	)
}

// renderButtons renders the action buttons
func (f TaskForm) renderButtons() string {
	createStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Foreground(styles.Success)

	cancelStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Foreground(styles.Muted)

	if f.focusedField == fieldButtons {
		createStyle = createStyle.
			Background(styles.Success).
			Foreground(styles.Background).
			Bold(true)
	}

	create := createStyle.Render("[ Create ]")
	cancel := cancelStyle.Render("[ Cancel (Esc) ]")

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		create,
		"  ",
		cancel,
	)
}

// blurAll removes focus from all fields
func (f *TaskForm) blurAll() {
	f.titleInput.Blur()
	f.descriptionInput.Blur()
	f.dueDateInput.Blur()
}

// focusField focuses a specific field
func (f *TaskForm) focusField(field int) tea.Cmd {
	switch field {
	case fieldTitle:
		return f.titleInput.Focus()
	case fieldDescription:
		return f.descriptionInput.Focus()
	case fieldDueDate:
		return f.dueDateInput.Focus()
	}
	return nil
}

// GetTask returns the task from form data
func (f TaskForm) GetTask() *models.Task {
	task := models.NewTask(f.titleInput.Value())
	task.Description = f.descriptionInput.Value()
	task.Priority = f.priorities[f.selectedPriority]
	task.Status = models.TaskStatusPending // Always create as pending
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	
	// Parse due date if provided
	dueDateStr := strings.TrimSpace(f.dueDateInput.Value())
	if dueDateStr != "" {
		// Try parsing with YYYY-MM-DD format
		if dueDate, err := time.Parse("2006-01-02", dueDateStr); err == nil {
			task.DueDate = &dueDate
		}
		// Could also try other formats like DD/MM/YYYY
		// if dueDate, err := time.Parse("02/01/2006", dueDateStr); err == nil {
		//     task.DueDate = &dueDate
		// }
	}
	
	return task
}

// IsSubmitted returns true if form was submitted
func (f TaskForm) IsSubmitted() bool {
	return f.submitted
}

// IsCancelled returns true if form was cancelled
func (f TaskForm) IsCancelled() bool {
	return f.cancelled
}
