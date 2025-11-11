package screens

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

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
	err         error
	feedbackMsg string

	// Form state
	showForm          bool
	showDeleteConfirm bool
	showDetails       bool
	taskForm          components.TaskForm

	// Move mode state
	moveMode     bool
	targetColumn Column

	// Details view state
	subtaskCursor      int
	isCreatingSubtask  bool
	subtaskInput       components.Input
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
		showDetails:    false,
		moveMode:       false,
		subtaskCursor:  0,
		isCreatingSubtask: false,
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

	// Handle delete confirmation
	if s.showDeleteConfirm {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "y", "Y":
				s.showDeleteConfirm = false
				if s.selectedTaskID != "" {
					return s, s.deleteTask(s.selectedTaskID)
				}
				return s, nil
			case "n", "N", "esc":
				s.showDeleteConfirm = false
				return s, nil
			}
		}
		return s, nil
	}

	// Handle move mode
	if s.moveMode {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "left", "h":
				if s.targetColumn > ColumnTodo {
					s.targetColumn--
				}
				return s, nil
			case "right", "l":
				if s.targetColumn < ColumnDone {
					s.targetColumn++
				}
				return s, nil
			case "enter":
				s.moveMode = false
				if s.selectedTaskID != "" {
					return s, s.moveTaskToColumn(s.selectedTaskID, s.targetColumn)
				}
				return s, nil
			case "esc", "m", "q":
				s.moveMode = false
				return s, nil
			}
		}
		return s, nil
	}

	// Handle details view
	if s.showDetails {
		// Handle subtask creation mode
		if s.isCreatingSubtask {
			var cmd tea.Cmd
			cmd = s.subtaskInput.Update(msg)

			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.String() {
				case "enter":
					title := s.subtaskInput.Value()
					if title != "" {
						s.isCreatingSubtask = false
						return s, s.createSubtask(title)
					}
				case "esc":
					s.isCreatingSubtask = false
				}
			}
			return s, cmd
		}

		// Handle normal details view navigation
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter", "q":
				s.showDetails = false
				s.subtaskCursor = 0 // Reset cursor
				return s, nil
			case "j", "down":
				task := s.getTaskByID(s.selectedTaskID)
				if task != nil && s.subtaskCursor < len(task.Subtasks)-1 {
					s.subtaskCursor++
				}
				return s, nil
			case "k", "up":
				if s.subtaskCursor > 0 {
					s.subtaskCursor--
				}
				return s, nil
			case " ":
				task := s.getTaskByID(s.selectedTaskID)
				if task != nil && s.subtaskCursor < len(task.Subtasks) {
					return s, s.toggleSubtask()
				}
			case "t":
				s.isCreatingSubtask = true
				s.subtaskInput = components.NewInput("", "New subtask title...")
				return s, s.subtaskInput.Focus()
			}
		}
		return s, nil
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
			// If in details view, exit details view
			if s.showDetails {
				s.showDetails = false
				s.selectedTaskID = ""
				return s, nil
			}

			// Otherwise, enter details view for the task under the cursor
			tasks := s.getTasksForColumn(s.activeColumn)
			if s.cursors[s.activeColumn] < len(tasks) {
				task := tasks[s.cursors[s.activeColumn]]
				s.selectedTaskID = task.ID // Select the task to show its details
				s.showDetails = true
			}
			return s, nil

		case " ": // Spacebar for selection
			// In details or move mode, space does nothing
			if s.showDetails || s.moveMode {
				return s, nil
			}
			tasks := s.getTasksForColumn(s.activeColumn)
			if s.cursors[s.activeColumn] < len(tasks) {
				task := tasks[s.cursors[s.activeColumn]]
				// If it's already selected, deselect it
				if s.selectedTaskID == task.ID {
					s.selectedTaskID = ""
				} else { // Otherwise, select it
					s.selectedTaskID = task.ID
				}
			}
			return s, nil


		case "m":
			// Enter move mode
			if s.selectedTaskID != "" {
				s.moveMode = true
				s.targetColumn = s.activeColumn
			}
			return s, nil

		case "delete", "backspace":
			// Show delete confirmation
			if s.selectedTaskID != "" {
				s.showDeleteConfirm = true
			}
			return s, nil

		        case "r":
		            // Refresh tasks
		            return s, s.loadTasks()
		
		        case "x":
		            // Export tasks
		            return s, s.exportTasks()
		        
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

			// Sort tasks with custom logic
			priorityMap := map[models.TaskPriority]int{
				models.TaskPriorityUrgent: 4,
				models.TaskPriorityHigh:   3,
				models.TaskPriorityMedium: 2,
				models.TaskPriorityLow:    1,
			}
			sort.Slice(s.tasks, func(i, j int) bool {
				taskA := s.tasks[i]
				taskB := s.tasks[j]

				// Rule 1: Overdue status (overdue tasks first)
				isAOverdue := taskA.IsOverdue()
				isBOverdue := taskB.IsOverdue()
				if isAOverdue != isBOverdue {
					return isAOverdue
				}

				// Rule 2: Priority (higher number is higher priority)
				priorityA := priorityMap[taskA.Priority]
				priorityB := priorityMap[taskB.Priority]
				if priorityA != priorityB {
					return priorityA > priorityB
				}

				// Rule 3: Due Date (tasks with due dates first, then earlier dates first)
				if taskA.DueDate != nil && taskB.DueDate == nil {
					return true
				}
				if taskA.DueDate == nil && taskB.DueDate != nil {
					return false
				}
				if taskA.DueDate != nil && taskB.DueDate != nil {
					if !taskA.DueDate.Equal(*taskB.DueDate) {
						return taskA.DueDate.Before(*taskB.DueDate)
					}
				}

				// Rule 4: Creation Date (older first)
				return taskA.CreatedAt.Before(taskB.CreatedAt)
			})

		        // Adjust cursors if needed
		        for col := ColumnTodo; col <= ColumnDone; col++ {
		            tasks := s.getTasksForColumn(col)
		            if s.cursors[col] >= len(tasks) && len(tasks) > 0 {
		                s.cursors[col] = len(tasks) - 1
		            }
		        }
		        return s, nil
		
		    	case tasksExportedMsg:
		    		if msg.err != nil {
		    			s.feedbackMsg = lipgloss.NewStyle().Foreground(styles.Danger).Render(fmt.Sprintf("Export failed: %v", msg.err))
		    		} else {
		    			s.feedbackMsg = lipgloss.NewStyle().Foreground(styles.Success).Render(fmt.Sprintf("Tasks exported to %s", msg.filename))
		    		}
		    		return s, tea.Tick(3*time.Second, func(t time.Time) tea.Msg {		            return clearFeedbackMsg{}
		        })
		
		    case clearFeedbackMsg:
		        s.feedbackMsg = ""
		        return s, nil
		
		    case taskCreatedMsg:		if msg.err != nil {
			s.err = msg.err
		} else {
		}
		// Reload tasks after creation
		return s, s.loadTasks()

	case taskUpdatedMsg:
		s.showDetails = false
		if msg.err != nil {
			s.err = msg.err
		} else {
			// Deselect after successful update
			s.selectedTaskID = ""
		}
		// Reload tasks after update
		return s, s.loadTasks()

	case subtaskToggledMsg:
		if msg.err != nil {
			s.err = msg.err
		}
		// Reload tasks to reflect the change
		return s, s.loadTasks()

	case subtaskCreatedMsg:
		if msg.err != nil {
			s.err = msg.err
		}
		// Reload tasks to reflect the change
		return s, s.loadTasks()

	case taskMovedMsg:
		s.showDetails = false
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

	// Decide which main view to render
	var mainView string
	if s.showDetails {
		mainView = s.renderDetailsView()
	} else {
		mainView = s.renderKanban()
	}

	// If form is shown, overlay it on top
	if s.showForm {
		return s.overlayForm(mainView)
	}

	// If delete confirmation is shown, overlay it
	if s.showDeleteConfirm {
		return s.renderDeleteConfirmDialog(mainView)
	}

	return mainView
}

// renderDeleteConfirmDialog renders the delete confirmation dialog over the base view
func (s *TaskScreen) renderDeleteConfirmDialog(baseView string) string {
	// Find the task to get its title
	var taskTitle string
	for _, task := range s.tasks {
		if task.ID == s.selectedTaskID {
			taskTitle = task.Title
			break
		}
	}

	question := fmt.Sprintf("Delete task \"%s\"?", taskTitle)

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
				styles.Shortcut.Render("y") + styles.ShortcutText.Render(" delete"),
				"  ",
				styles.Shortcut.Render("n") + styles.ShortcutText.Render(" cancel"),
			),
		))

	return lipgloss.Place(
		s.width,
		s.height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)
}

// getTaskByID finds a task by its ID
func (s *TaskScreen) getTaskByID(id string) *models.Task {
	for _, task := range s.tasks {
		if task.ID == id {
			return &task
		}
	}
	return nil
}

// renderDetailsView renders the detailed view of a single task
func (s *TaskScreen) renderDetailsView() string {
	task := s.getTaskByID(s.selectedTaskID)
	if task == nil {
		return "Task not found."
	}

	// Styles
	titleStyle := styles.Title.Copy().Bold(true)
	labelStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.Muted)
	descStyle := lipgloss.NewStyle().Width(s.width - 10).Padding(0, 2)

	// --- Content ---
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render(task.Title))
	b.WriteString("\n\n")

	// Description
	b.WriteString(labelStyle.Render("Description"))
	b.WriteString("\n")
	if task.Description != "" {
		b.WriteString(descStyle.Render(task.Description))
	} else {
		b.WriteString(descStyle.Render(styles.Dimmed.Render("No description.")))
	}
	b.WriteString("\n\n")

	// Details grid
	details := lipgloss.NewStyle().Width(s.width / 2).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			labelStyle.Render("Status")+": "+task.Status.String(),
			labelStyle.Render("Priority")+": "+task.Priority.String(),
		),
	)
	b.WriteString(details)
	b.WriteString("\n")

	// Due Date
	if task.DueDate != nil {
		dueStr := task.DueDate.Format("Mon, 02 Jan 2006")
		if task.IsOverdue() {
			dueStr += " (Overdue)"
		}
		b.WriteString(labelStyle.Render("Due")+": "+dueStr)
		b.WriteString("\n")
	}

	// Tags
	if len(task.Tags) > 0 {
		var tagStrings []string
		for _, tag := range task.Tags {
			tagStrings = append(tagStrings, styles.Tag.Render(tag))
		}
		b.WriteString(labelStyle.Render("Tags")+": "+strings.Join(tagStrings, " "))
		b.WriteString("\n")
	}

	// Subtasks
	if len(task.Subtasks) > 0 {
		b.WriteString("\n")
		b.WriteString(labelStyle.Render("Checklist"))
		b.WriteString("\n")
		for i, st := range task.Subtasks {
			checkbox := "[ ]"
			if st.IsCompleted {
				checkbox = "[x]"
			}

			subtaskStyle := lipgloss.NewStyle()
			if i == s.subtaskCursor {
				subtaskStyle = subtaskStyle.Foreground(styles.Primary)
			}

			title := st.Title
			if st.IsCompleted {
				title = lipgloss.NewStyle().Strikethrough(true).Render(title)
			}

			b.WriteString(subtaskStyle.Render(fmt.Sprintf("  %s %s\n", checkbox, title)))
		}
	}

	// Render subtask input if creating
	if s.isCreatingSubtask {
		b.WriteString("\n" + s.subtaskInput.View())
	}

	// --- Layout ---
	containerStyle := styles.Panel.Copy().
		BorderForeground(styles.Primary).
		Width(s.width - 2).
		Height(s.height - 8).
		Padding(2, 4)

	content := containerStyle.Render(b.String())

	// Shortcuts
	shortcuts := s.renderShortcuts()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		"", // Spacer
		shortcuts,
	)
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

	if s.moveMode && s.targetColumn == column {
		columnStyle = styles.PanelTarget.
			Width(width).
			Height(s.height - 8)
	} else if s.activeColumn == column {
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

// renderShortcuts renders keyboard shortcuts or a feedback message
func (s *TaskScreen) renderShortcuts() string {
	// If there's a feedback message, show it instead of shortcuts
	if s.feedbackMsg != "" {
		return s.feedbackMsg
	}

	var shortcuts []string

	if s.moveMode {
		shortcuts = []string{
			styles.Shortcut.Render("←/→") + styles.ShortcutText.Render(" select column"),
			styles.Shortcut.Render("enter") + styles.ShortcutText.Render(" confirm"),
			styles.Shortcut.Render("esc") + styles.ShortcutText.Render(" cancel"),
		}
	} else if s.showDetails {
		shortcuts = []string{
			styles.Shortcut.Render("enter") + styles.ShortcutText.Render(" close"),
			styles.Shortcut.Render("j/k") + styles.ShortcutText.Render(" nav"),
			styles.Shortcut.Render("space") + styles.ShortcutText.Render(" toggle"),
			styles.Shortcut.Render("t") + styles.ShortcutText.Render(" new subtask"),
		}
	} else if s.selectedTaskID != "" {
		shortcuts = []string{
			styles.Shortcut.Render("space") + styles.ShortcutText.Render(" deselect"),
			styles.Shortcut.Render("m") + styles.ShortcutText.Render(" move"),
			styles.Shortcut.Render("del") + styles.ShortcutText.Render(" delete"),
			styles.Shortcut.Render("enter") + styles.ShortcutText.Render(" details"),
			styles.Shortcut.Render("e") + styles.ShortcutText.Render(" edit"),
		}
	} else {
		shortcuts = []string{
			styles.Shortcut.Render("space") + styles.ShortcutText.Render(" select"),
			styles.Shortcut.Render("tab") + styles.ShortcutText.Render(" next column"),
			styles.Shortcut.Render("enter") + styles.ShortcutText.Render(" details"),
			styles.Shortcut.Render("j/k") + styles.ShortcutText.Render(" navigate"),
			styles.Shortcut.Render("n") + styles.ShortcutText.Render(" new"),
			styles.Shortcut.Render("r") + styles.ShortcutText.Render(" refresh"),
			styles.Shortcut.Render("x") + styles.ShortcutText.Render(" export"),
		}
	}

	shortcutLine := strings.Join(shortcuts, "  ")

	// Total count
	totalText := fmt.Sprintf("Total: %d tasks", len(s.tasks))
	totalLine := styles.Dimmed.Render(totalText)

	return lipgloss.JoinVertical(lipgloss.Left, shortcutLine, totalLine)
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

func (s *TaskScreen) createSubtask(title string) tea.Cmd {
	return func() tea.Msg {
		subtask := models.Subtask{
			TaskID:      s.selectedTaskID,
			Title:       title,
			IsCompleted: false,
		}
		err := s.db.Tasks().CreateSubtask(&subtask)
		return subtaskCreatedMsg{err: err}
	}
}

// deleteTask deletes a task by ID
func (s *TaskScreen) deleteTask(taskID string) tea.Cmd {
	return func() tea.Msg {
		err := s.db.Tasks().Delete(taskID)
		return taskDeletedMsg{err: err}
	}
}

func (s *TaskScreen) toggleSubtask() tea.Cmd {
	return func() tea.Msg {
		task := s.getTaskByID(s.selectedTaskID)
		if task == nil || s.subtaskCursor >= len(task.Subtasks) {
			return subtaskToggledMsg{err: fmt.Errorf("subtask not found")}
		}

		subtask := &task.Subtasks[s.subtaskCursor]
		subtask.IsCompleted = !subtask.IsCompleted

		err := s.db.Tasks().UpdateSubtask(subtask)
		return subtaskToggledMsg{err: err}
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

type subtaskToggledMsg struct {
	err error
}

type subtaskCreatedMsg struct {
	err error
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

// exportTasks exports all tasks to a JSON file
func (s *TaskScreen) exportTasks() tea.Cmd {
	return func() tea.Msg {
		if len(s.tasks) == 0 {
			return tasksExportedMsg{err: fmt.Errorf("no tasks to export")}
		}

		// Marshal the tasks to pretty-printed JSON
		data, err := json.MarshalIndent(s.tasks, "", "  ")
		if err != nil {
			return tasksExportedMsg{err: err}
		}

		// Create a timestamped filename
		filename := fmt.Sprintf("unicli_export_%s.json", time.Now().Format("20060102_150405"))

		// Write the file
		if err := os.WriteFile(filename, data, 0644); err != nil {
			return tasksExportedMsg{err: err}
		}

		return tasksExportedMsg{filename: filename, err: nil}
	}
}

// tasksExportedMsg is sent when tasks have been exported
type tasksExportedMsg struct {
	filename string
	err      error
}

// clearFeedbackMsg is sent to clear the feedback message after a delay
type clearFeedbackMsg struct{}
