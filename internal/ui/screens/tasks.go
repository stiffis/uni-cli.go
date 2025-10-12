package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stiffis/UniCLI/internal/database"
	"github.com/stiffis/UniCLI/internal/models"
	"github.com/stiffis/UniCLI/internal/ui/styles"
)

// TaskScreen is the tasks view
type TaskScreen struct {
	db       *database.DB
	tasks    []models.Task
	cursor   int
	selected map[int]struct{}
	width    int
	height   int
	loading  bool
	err      error
}

// NewTaskScreen creates a new task screen
func NewTaskScreen(db *database.DB) *TaskScreen {
	return &TaskScreen{
		db:       db,
		tasks:    []models.Task{},
		selected: make(map[int]struct{}),
		loading:  true,
	}
}

// Init initializes the task screen
func (s *TaskScreen) Init() tea.Cmd {
	return s.loadTasks()
}

// Update handles messages
func (s *TaskScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if s.cursor < len(s.tasks)-1 {
				s.cursor++
			}
		case "k", "up":
			if s.cursor > 0 {
				s.cursor--
			}
		case "g":
			s.cursor = 0
		case "G":
			if len(s.tasks) > 0 {
				s.cursor = len(s.tasks) - 1
			}
		case "space":
			if s.cursor < len(s.tasks) {
				return s, s.toggleTask()
			}
		case "r":
			// Refresh tasks
			return s, s.loadTasks()
		case "n":
			// TODO: Open new task form
			return s, nil
		case "d":
			// Delete current task
			if s.cursor < len(s.tasks) {
				return s, s.deleteTask()
			}
			return s, nil
		case "e":
			// TODO: Edit task
			return s, nil
		}

	case tasksLoadedMsg:
		s.tasks = msg.tasks
		s.loading = false
		s.err = msg.err
		if s.cursor >= len(s.tasks) && len(s.tasks) > 0 {
			s.cursor = len(s.tasks) - 1
		}
		return s, nil

	case taskToggledMsg:
		if msg.err != nil {
			s.err = msg.err
		}
		// Reload tasks after toggle
		return s, s.loadTasks()

	case taskDeletedMsg:
		if msg.err != nil {
			s.err = msg.err
		}
		// Reload tasks after delete
		return s, s.loadTasks()
	}

	return s, nil
}

// View renders the task screen
func (s *TaskScreen) View() string {
	// Show loading state
	if s.loading {
		return lipgloss.NewStyle().
			Padding(2).
			Foreground(styles.Info).
			Render("Loading tasks...")
	}

	// Show error state
	if s.err != nil {
		errorMsg := fmt.Sprintf("Error: %v\n\nPress 'r' to retry", s.err)
		return lipgloss.NewStyle().
			Padding(2).
			Foreground(styles.Danger).
			Render(errorMsg)
	}

	if len(s.tasks) == 0 {
		return s.renderEmpty()
	}

	var lines []string
	
	// Section: Today
	lines = append(lines, styles.Title.Render("Today's Tasks"))
	lines = append(lines, "")
	
	todayTasks := s.filterToday()
	if len(todayTasks) == 0 {
		lines = append(lines, styles.Dimmed.Render("  No tasks due today"))
	} else {
		for i, task := range todayTasks {
			lines = append(lines, s.renderTask(task, i == s.cursor))
		}
	}
	
	lines = append(lines, "")
	
	// Section: Upcoming
	lines = append(lines, styles.Title.Render("Upcoming"))
	lines = append(lines, "")
	
	upcomingTasks := s.filterUpcoming()
	if len(upcomingTasks) == 0 {
		lines = append(lines, styles.Dimmed.Render("  No upcoming tasks"))
	} else {
		for i, task := range upcomingTasks {
			idx := len(todayTasks) + i
			lines = append(lines, s.renderTask(task, idx == s.cursor))
		}
	}
	
	lines = append(lines, "")
	
	// Section: All Tasks
	lines = append(lines, styles.Title.Render("All Tasks"))
	lines = append(lines, "")
	
	for i, task := range s.tasks {
		lines = append(lines, s.renderTask(task, i == s.cursor))
	}

	// Shortcuts
	lines = append(lines, "")
	lines = append(lines, s.renderShortcuts())

	return strings.Join(lines, "\n")
}

// renderTask renders a single task
func (s *TaskScreen) renderTask(task models.Task, selected bool) string {
	// Status icon
	var icon string
	switch task.Status {
	case models.TaskStatusCompleted:
		icon = "[✓]"
	case models.TaskStatusInProgress:
		icon = "[~]"
	case models.TaskStatusCancelled:
		icon = "[x]"
	default:
		icon = "[ ]"
	}

	// Priority indicator
	var priorityIndicator string
	switch task.Priority {
	case models.TaskPriorityUrgent:
		priorityIndicator = "!!"
	case models.TaskPriorityHigh:
		priorityIndicator = "! "
	case models.TaskPriorityMedium:
		priorityIndicator = "- "
	case models.TaskPriorityLow:
		priorityIndicator = "  "
	}

	// Due date
	var dueStr string
	if task.DueDate != nil {
		if task.IsOverdue() {
			dueStr = lipgloss.NewStyle().
				Foreground(styles.Danger).
				Render(fmt.Sprintf(" (Overdue: %s)", task.DueDate.Format("Jan 02")))
		} else if task.IsDueToday() {
			dueStr = lipgloss.NewStyle().
				Foreground(styles.Warning).
				Render(" (Due today)")
		} else {
			dueStr = lipgloss.NewStyle().
				Foreground(styles.Muted).
				Render(fmt.Sprintf(" (%s)", task.DueDate.Format("Jan 02")))
		}
	}

	// Task title with status color
	titleStyle := lipgloss.NewStyle().
		Foreground(styles.StatusColor(string(task.Status)))
	
	if task.Status == models.TaskStatusCompleted {
		titleStyle = titleStyle.Strikethrough(true)
	}

	taskLine := fmt.Sprintf("  %s %s %s%s",
		icon,
		priorityIndicator,
		titleStyle.Render(task.Title),
		dueStr,
	)

	// Apply selection style
	if selected {
		return styles.ListItemSelected.Render(taskLine)
	}
	return taskLine
}

// renderEmpty renders empty state
func (s *TaskScreen) renderEmpty() string {
	return lipgloss.NewStyle().
		Padding(2).
		Render(strings.Join([]string{
			styles.Title.Render("No tasks yet!"),
			"",
			styles.Dimmed.Render("Press 'n' to create your first task"),
		}, "\n"))
}

// renderShortcuts renders keyboard shortcuts
func (s *TaskScreen) renderShortcuts() string {
	shortcuts := []string{
		styles.Shortcut.Render("n") + styles.ShortcutText.Render(" new"),
		styles.Shortcut.Render("e") + styles.ShortcutText.Render(" edit"),
		styles.Shortcut.Render("d") + styles.ShortcutText.Render(" delete"),
		styles.Shortcut.Render("space") + styles.ShortcutText.Render(" toggle"),
		styles.Shortcut.Render("r") + styles.ShortcutText.Render(" refresh"),
		styles.Shortcut.Render("j/k") + styles.ShortcutText.Render(" navigate"),
	}
	return strings.Join(shortcuts, "  ")
}

// filterToday returns tasks due today
func (s *TaskScreen) filterToday() []models.Task {
	var result []models.Task
	for _, task := range s.tasks {
		if task.IsDueToday() && task.Status != models.TaskStatusCompleted {
			result = append(result, task)
		}
	}
	return result
}

// filterUpcoming returns upcoming tasks (next 7 days)
func (s *TaskScreen) filterUpcoming() []models.Task {
	var result []models.Task
	for _, task := range s.tasks {
		if task.DueDate != nil && !task.IsDueToday() && !task.IsOverdue() && task.Status != models.TaskStatusCompleted {
			result = append(result, task)
		}
	}
	return result
}

// toggleTask toggles task completion
func (s *TaskScreen) toggleTask() tea.Cmd {
	if s.cursor >= len(s.tasks) {
		return nil
	}

	taskID := s.tasks[s.cursor].ID

	return func() tea.Msg {
		err := s.db.Tasks().ToggleComplete(taskID)
		return taskToggledMsg{err: err}
	}
}

// deleteTask deletes the current task
func (s *TaskScreen) deleteTask() tea.Cmd {
	if s.cursor >= len(s.tasks) {
		return nil
	}

	taskID := s.tasks[s.cursor].ID

	return func() tea.Msg {
		err := s.db.Tasks().Delete(taskID)
		return taskDeletedMsg{err: err}
	}
}

// loadTasks loads tasks from database
func (s *TaskScreen) loadTasks() tea.Cmd {
	return func() tea.Msg {
		tasks, err := s.db.Tasks().FindAll()
		if err != nil {
			return tasksLoadedMsg{tasks: []models.Task{}, err: err}
		}
		return tasksLoadedMsg{tasks: tasks, err: nil}
	}
}

// Messages
type tasksLoadedMsg struct {
	tasks []models.Task
	err   error
}

type taskToggledMsg struct {
	err error
}

type taskDeletedMsg struct {
	err error
}
