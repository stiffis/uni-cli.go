package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stiffis/UniCLI/internal/database"
	"github.com/stiffis/UniCLI/internal/models"
	"github.com/stiffis/UniCLI/internal/ui/components"
	"github.com/stiffis/UniCLI/internal/ui/styles"
)

// Column represents a kanban column
type Column int

const (
	ColumnTodo Column = iota
	ColumnInProgress
	ColumnDone
)

// TaskScreen is the tasks view
type TaskScreen struct {
	db             *database.DB
	tasks          []models.Task
	activeColumn   Column         // Which column has focus
	cursors        map[Column]int // Cursor position for each column
	selectedTaskID string         // ID of selected task (empty if none)
	width          int
	height         int
	loading        bool
	err            error

	// Form state
	showForm bool
	taskForm components.TaskForm
}

// NewTaskScreen creates a new task screen
func NewTaskScreen(db *database.DB) *TaskScreen {
	return &TaskScreen{
		db:           db,
		tasks:        []models.Task{},
		activeColumn: ColumnTodo,
		cursors: map[Column]int{
			ColumnTodo:       0,
			ColumnInProgress: 0,
			ColumnDone:       0,
		},
		selectedTaskID: "",
		loading:        true,
		showForm:       false,
	}
}

// Init initializes the task screen
func (s *TaskScreen) Init() tea.Cmd {
	return s.loadTasks()
}

// Update handles messages
func (s *TaskScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If form is shown, handle form updates
	if s.showForm {
		var cmd tea.Cmd
		s.taskForm, cmd = s.taskForm.Update(msg)

		// Check if form was submitted or cancelled
		if s.taskForm.IsSubmitted() {
			// Get the task from the form
			task := s.taskForm.GetTask()
			s.showForm = false

			// Check if it's a new task or editing existing
			if s.taskForm.IsNewTask() {
				return s, s.createTask(task)
			} else {
				return s, s.updateTask(task)
			}
		}

		if s.taskForm.IsCancelled() {
			s.showForm = false
			return s, nil
		}

		return s, cmd
	}

	// Normal kanban view handling
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// Move to next column
			s.activeColumn = (s.activeColumn + 1) % 3
			return s, nil

		case "shift+tab":
			// Move to previous column
			s.activeColumn = (s.activeColumn + 2) % 3
			return s, nil

		case "j", "down":
			// Move cursor down in current column
			tasks := s.getTasksForColumn(s.activeColumn)
			if s.cursors[s.activeColumn] < len(tasks)-1 {
				s.cursors[s.activeColumn]++
			}
			return s, nil

		case "k", "up":
			// Move cursor up in current column
			if s.cursors[s.activeColumn] > 0 {
				s.cursors[s.activeColumn]--
			}
			return s, nil

		case "g":
			// Go to top of column
			s.cursors[s.activeColumn] = 0
			return s, nil

		case "G":
			// Go to bottom of column
			tasks := s.getTasksForColumn(s.activeColumn)
			if len(tasks) > 0 {
				s.cursors[s.activeColumn] = len(tasks) - 1
			}
			return s, nil

		case "enter":
			// Select/deselect task
			tasks := s.getTasksForColumn(s.activeColumn)
			if s.cursors[s.activeColumn] < len(tasks) {
				task := tasks[s.cursors[s.activeColumn]]
				if s.selectedTaskID == task.ID {
					// Deselect
					s.selectedTaskID = ""
				} else {
					// Select
					s.selectedTaskID = task.ID
				}
			}
			return s, nil

		case "left", "h":
			// Move selected task to previous column
			if s.selectedTaskID != "" {
				return s, s.moveTaskToColumn(s.selectedTaskID, s.getPreviousColumn())
			}
			return s, nil

		case "right", "l":
			// Move selected task to next column
			if s.selectedTaskID != "" {
				return s, s.moveTaskToColumn(s.selectedTaskID, s.getNextColumn())
			}
			return s, nil

		case "delete", "backspace":
			// Delete selected task
			if s.selectedTaskID != "" {
				return s, s.deleteTask(s.selectedTaskID)
			}
			return s, nil

		case "r":
			// Refresh tasks
			return s, s.loadTasks()

		case "n":
			// Open new task form
			s.showForm = true
			s.taskForm = components.NewTaskForm(nil) // Pass nil for new task
			return s, nil

		case "e":
			// Edit selected task
			if s.selectedTaskID != "" {
				// Find the task to edit
				var taskToEdit *models.Task
				for _, t := range s.tasks {
					if t.ID == s.selectedTaskID {
						taskToEdit = &t
						break
					}
				}
				if taskToEdit != nil {
					s.showForm = true
					s.taskForm = components.NewTaskForm(taskToEdit)
				}
				return s, nil
			}
			return s, nil
		}

	case tasksLoadedMsg:
		s.tasks = msg.tasks
		s.loading = false
		s.err = msg.err
		// Adjust cursors if needed
		for col := ColumnTodo; col <= ColumnDone; col++ {
			tasks := s.getTasksForColumn(col)
			if s.cursors[col] >= len(tasks) && len(tasks) > 0 {
				s.cursors[col] = len(tasks) - 1
			}
		}
		return s, nil

	case taskCreatedMsg:
		if msg.err != nil {
			s.err = msg.err
		} else {
		}
		// Reload tasks after creation
		return s, s.loadTasks()

	case taskUpdatedMsg:
		if msg.err != nil {
			s.err = msg.err
		} else {
			// Deselect after successful update
			s.selectedTaskID = ""
		}
		// Reload tasks after update
		return s, s.loadTasks()

	case taskMovedMsg:
		if msg.err != nil {
			s.err = msg.err
		} else {
			// Deselect after successful move
			s.selectedTaskID = ""
		}
		// Reload tasks after move
		return s, s.loadTasks()

	case taskDeletedMsg:
		if msg.err != nil {
			s.err = msg.err
		} else {
			s.selectedTaskID = ""
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

	// Render base kanban view
	baseView := s.renderKanban()

	// If form is shown, overlay it on top
	if s.showForm {
		return s.overlayForm(baseView)
	}

	return baseView
}

// renderKanban renders the kanban board
func (s *TaskScreen) renderKanban() string {
	// Get tasks for each column
	todoTasks := s.getTasksForColumn(ColumnTodo)
	inProgressTasks := s.getTasksForColumn(ColumnInProgress)
	doneTasks := s.getTasksForColumn(ColumnDone)

	// Calculate column width (divide available width by 3, minus borders and spacing)
	// Add 3 extra characters per column for more width
	columnWidth := ((s.width - 10) / 3) + 3
	if columnWidth < 23 {
		columnWidth = 23
	}

	// Render each column
	todoColumn := s.renderColumn("To Do", todoTasks, ColumnTodo, columnWidth)
	inProgressColumn := s.renderColumn("In Progress", inProgressTasks, ColumnInProgress, columnWidth)
	doneColumn := s.renderColumn("Done", doneTasks, ColumnDone, columnWidth)

	// Combine columns horizontally
	kanbanBoard := lipgloss.JoinHorizontal(
		lipgloss.Top,
		todoColumn,
		inProgressColumn,
		doneColumn,
	)

	// Add shortcuts at the bottom
	shortcuts := s.renderShortcuts()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		kanbanBoard,
		"",
		shortcuts,
	)
}

// overlayForm overlays the form on top of the kanban view
func (s *TaskScreen) overlayForm(baseView string) string {
	// Render the form
	formView := s.taskForm.View()

	// Calculate position to center the form
	baseWidth := s.width
	baseHeight := s.height

	// Place form using lipgloss Place
	overlay := lipgloss.Place(
		baseWidth,
		baseHeight,
		lipgloss.Center,
		lipgloss.Center,
		formView,
		lipgloss.WithWhitespaceChars(""),
	)

	return overlay
}

// renderColumn renders a single kanban column
func (s *TaskScreen) renderColumn(title string, tasks []models.Task, column Column, width int) string {
	// Column header with count
	headerText := fmt.Sprintf("%s (%d)", title, len(tasks))

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		Width(width).
		Align(lipgloss.Center)

	// Highlight active column
	columnStyle := styles.Panel.
		Width(width).
		Height(s.height - 8)

	if s.activeColumn == column {
		columnStyle = columnStyle.BorderForeground(styles.Primary)
		headerStyle = headerStyle.Foreground(styles.Secondary)
	}

	header := headerStyle.Render(headerText)

	// Render tasks
	var taskLines []string
	for i, task := range tasks {
		isSelected := task.ID == s.selectedTaskID
		isCursor := i == s.cursors[column] && s.activeColumn == column
		taskLine := s.renderKanbanTask(task, isSelected, isCursor)
		taskLines = append(taskLines, taskLine)

		// Add divider after each task (except the last one)
		if i < len(tasks)-1 {
			divider := lipgloss.NewStyle().
				Foreground(styles.Border).
				Render("─────────────────────")
			taskLines = append(taskLines, divider)
		}
	}

	// Empty state
	if len(taskLines) == 0 {
		emptyMsg := lipgloss.NewStyle().
			Foreground(styles.Muted).
			Italic(true).
			Render("No tasks")
		taskLines = append(taskLines, emptyMsg)
	}

	content := strings.Join(taskLines, "\n")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		columnStyle.Render(content),
	)
}

// renderKanbanTask renders a single task in kanban view
func (s *TaskScreen) renderKanbanTask(task models.Task, isSelected bool, isCursor bool) string {
	// Priority indicator with Nerd Font icons
	var priorityIndicator string

	switch task.Priority {
	case models.TaskPriorityUrgent:
		priorityIndicator = "☠️"
	case models.TaskPriorityHigh:
		priorityIndicator = "🔴"
	case models.TaskPriorityMedium:
		priorityIndicator = "💛"
	case models.TaskPriorityLow:
		priorityIndicator = "🤍"
	}

	// Task title
	titleStyle := lipgloss.NewStyle()
	if task.Status == models.TaskStatusCompleted {
		titleStyle = titleStyle.Strikethrough(true).Foreground(styles.Muted)
	}

	title := titleStyle.Render(task.Title)

	// Due date indicator with Nerd Font icons
	var dueInfo string
	if task.DueDate != nil {
		if task.IsOverdue() {
			dueInfo = lipgloss.NewStyle().
				Foreground(styles.Danger).
				Render("  " + task.DueDate.Format("Jan 02"))
		} else if task.IsDueToday() {
			dueInfo = lipgloss.NewStyle().
				Foreground(styles.Warning).
				Render(" 󰃭 Today")
		} else { // Not overdue and not due today, so it's in the future
			dueInfo = lipgloss.NewStyle().
				Foreground(styles.Success).
				Render(" 󰃭 " + task.DueDate.Format("Jan 02"))
		}
	}

	taskText := fmt.Sprintf("%s %s%s", priorityIndicator, title, dueInfo)

	// Apply selection/cursor styles
	taskStyle := lipgloss.NewStyle().Padding(0, 1)

	if isSelected {
		// Selected task - highlighted background
		taskStyle = taskStyle.
			Background(styles.Secondary).
			Foreground(styles.Background).
			Bold(true)
	} else if isCursor {
		// Cursor position - subtle highlight
		taskStyle = taskStyle.
			Background(styles.BackgroundLight).
			Foreground(styles.Foreground)
	}

	return taskStyle.Render(taskText)
}

// renderEmpty renders empty state
func (s *TaskScreen) renderEmpty() string {
	totalTasks := len(s.tasks)

	return lipgloss.NewStyle().
		Padding(2).
		Render(strings.Join([]string{
			styles.Title.Render("No tasks yet!"),
			"",
			styles.Dimmed.Render(fmt.Sprintf("Total tasks: %d", totalTasks)),
			"",
			styles.Dimmed.Render("Press 'n' to create your first task"),
		}, "\n"))
}

// renderShortcuts renders keyboard shortcuts
func (s *TaskScreen) renderShortcuts() string {
	var shortcuts []string

	if s.selectedTaskID != "" {
		shortcuts = []string{
			styles.Shortcut.Render("←/→") + styles.ShortcutText.Render(" move column"),
			styles.Shortcut.Render("del") + styles.ShortcutText.Render(" delete"),
			styles.Shortcut.Render("enter") + styles.ShortcutText.Render(" deselect"),
			styles.Shortcut.Render("e") + styles.ShortcutText.Render(" edit"),
		}
	} else {
		shortcuts = []string{
			styles.Shortcut.Render("tab") + styles.ShortcutText.Render(" next column"),
			styles.Shortcut.Render("enter") + styles.ShortcutText.Render(" select"),
			styles.Shortcut.Render("j/k") + styles.ShortcutText.Render(" navigate"),
			styles.Shortcut.Render("n") + styles.ShortcutText.Render(" new"),
			styles.Shortcut.Render("r") + styles.ShortcutText.Render(" refresh"),
		}
	}

	// Total count
	totalText := fmt.Sprintf("Total: %d tasks", len(s.tasks))
	shortcuts = append(shortcuts, styles.Dimmed.Render(totalText))

	return strings.Join(shortcuts, "  ")
}

// getTasksForColumn returns tasks for a specific column
func (s *TaskScreen) getTasksForColumn(column Column) []models.Task {
	var tasks []models.Task
	for _, task := range s.tasks {
		switch column {
		case ColumnTodo:
			if task.Status == models.TaskStatusPending {
				tasks = append(tasks, task)
			}
		case ColumnInProgress:
			if task.Status == models.TaskStatusInProgress {
				tasks = append(tasks, task)
			}
		case ColumnDone:
			if task.Status == models.TaskStatusCompleted {
				tasks = append(tasks, task)
			}
		}
	}
	return tasks
}

// getPreviousColumn returns the previous column
func (s *TaskScreen) getPreviousColumn() Column {
	switch s.activeColumn {
	case ColumnTodo:
		return ColumnTodo // Can't go before Todo
	case ColumnInProgress:
		return ColumnTodo
	case ColumnDone:
		return ColumnInProgress
	}
	return ColumnTodo
}

// getNextColumn returns the next column
func (s *TaskScreen) getNextColumn() Column {
	switch s.activeColumn {
	case ColumnTodo:
		return ColumnInProgress
	case ColumnInProgress:
		return ColumnDone
	case ColumnDone:
		return ColumnDone // Can't go after Done
	}
	return ColumnDone
}

// moveTaskToColumn moves a task to a different column (status)
func (s *TaskScreen) moveTaskToColumn(taskID string, targetColumn Column) tea.Cmd {
	return func() tea.Msg {
		// Find the task
		task, err := s.db.Tasks().FindByID(taskID)
		if err != nil {
			return taskMovedMsg{err: err}
		}

		// Update status based on target column
		switch targetColumn {
		case ColumnTodo:
			task.Status = models.TaskStatusPending
			task.CompletedAt = nil
		case ColumnInProgress:
			task.Status = models.TaskStatusInProgress
			task.CompletedAt = nil
		case ColumnDone:
			task.Status = models.TaskStatusCompleted
			// CompletedAt will be set by the repository
		}

		// Save to database
		err = s.db.Tasks().Update(task)
		return taskMovedMsg{err: err}
	}
}

// createTask creates a new task
func (s *TaskScreen) createTask(task *models.Task) tea.Cmd {
	return func() tea.Msg {
		err := s.db.Tasks().Create(task)
		return taskCreatedMsg{err: err}
	}
}

// deleteTask deletes a task by ID
func (s *TaskScreen) deleteTask(taskID string) tea.Cmd {
	return func() tea.Msg {
		err := s.db.Tasks().Delete(taskID)
		return taskDeletedMsg{err: err}
	}
}

// updateTask updates an existing task
func (s *TaskScreen) updateTask(task *models.Task) tea.Cmd {
	return func() tea.Msg {
		err := s.db.Tasks().Update(task)
		return taskUpdatedMsg{err: err}
	}
}

// loadTasks loads tasks from database
func (s *TaskScreen) loadTasks() tea.Cmd {
	return func() tea.Msg {
		s.loading = true // Set loading to true before loading tasks
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

type taskCreatedMsg struct {
	err error
}

type taskMovedMsg struct {
	err error
}

type taskDeletedMsg struct {
	err error
}

type taskUpdatedMsg struct {
	err error
}
