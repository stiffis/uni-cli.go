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

// CoursesScreen represents the courses management screen
type CoursesScreen struct {
	db                *database.DB
	courses           []models.Course
	selectedIndex     int
	width             int
	height            int
	showForm          bool
	showDeleteConfirm bool
	courseForm        *components.CourseForm
	err               error
}

// NewCoursesScreen creates a new courses screen
func NewCoursesScreen(db *database.DB) *CoursesScreen {
	return &CoursesScreen{
		db:            db,
		courses:       []models.Course{},
		selectedIndex: 0,
	}
}

// Init initializes the courses screen
func (m CoursesScreen) Init() tea.Cmd {
	return m.fetchCoursesCmd()
}

// Update handles messages for the courses screen
func (m CoursesScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle form if active
	if m.showForm {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.showForm = false
				return m, m.fetchCoursesCmd()
			}
		case components.FormSubmitMsg:
			m.showForm = false
			return m, m.fetchCoursesCmd()
		}

		var newForm *components.CourseForm
		newForm, cmd = m.courseForm.Update(msg)
		m.courseForm = newForm
		return m, cmd
	}

	// Handle delete confirmation
	if m.showDeleteConfirm {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "y", "Y":
				m.showDeleteConfirm = false
				if m.selectedIndex >= 0 && m.selectedIndex < len(m.courses) {
					return m, m.deleteCourseCmd(m.courses[m.selectedIndex].ID)
				}
			case "n", "N", "esc":
				m.showDeleteConfirm = false
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.selectedIndex < len(m.courses)-1 {
				m.selectedIndex++
			}
		case "k", "up":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "g":
			m.selectedIndex = 0
		case "G":
			m.selectedIndex = len(m.courses) - 1
		case "n":
			// New course
			m.showForm = true
			m.courseForm = components.NewCourseForm(m.db, nil)
			cmd = m.courseForm.Init()
		case "e":
			// Edit course
			if m.selectedIndex >= 0 && m.selectedIndex < len(m.courses) {
				m.showForm = true
				m.courseForm = components.NewCourseForm(m.db, &m.courses[m.selectedIndex])
				cmd = m.courseForm.Init()
			}
		case "d":
			// Delete course
			if len(m.courses) > 0 {
				m.showDeleteConfirm = true
			}
		case "enter":
			// View course details (TODO: implement detail view)
			// For now, just edit
			if m.selectedIndex >= 0 && m.selectedIndex < len(m.courses) {
				m.showForm = true
				m.courseForm = components.NewCourseForm(m.db, &m.courses[m.selectedIndex])
				cmd = m.courseForm.Init()
			}
		}

	case fetchCoursesMsg:
		m.courses = msg.courses
		m.err = msg.err
		if m.selectedIndex >= len(m.courses) {
			m.selectedIndex = len(m.courses) - 1
			if m.selectedIndex < 0 {
				m.selectedIndex = 0
			}
		}
	}

	return m, cmd
}

// View renders the courses screen
func (m CoursesScreen) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Show form if active
	if m.showForm {
		return m.courseForm.View()
	}

	// Show delete confirmation if active
	if m.showDeleteConfirm {
		return m.renderDeleteConfirmation()
	}

	// Render main view
	title := m.renderTitle()
	list := m.renderCourseList()
	shortcuts := m.renderShortcuts()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		list,
		shortcuts,
	)
}

// renderTitle renders the title bar
func (m CoursesScreen) renderTitle() string {
	titleText := "󱉟 Courses"
	if len(m.courses) > 0 {
		titleText += fmt.Sprintf(" (%d)", len(m.courses))
	}

	return lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		Padding(1, 0).
		Render(titleText)
}

// renderCourseList renders the list of courses
func (m CoursesScreen) renderCourseList() string {
	if len(m.courses) == 0 {
		return lipgloss.NewStyle().
			Foreground(styles.Muted).
			Padding(2, 4).
			Render("No courses yet. Press 'n' to create one.")
	}

	var items []string
	maxHeight := m.height - 10 // Reserve space for title and shortcuts

	for i, course := range m.courses {
		if len(items) >= maxHeight {
			break
		}

		item := m.renderCourseItem(course, i == m.selectedIndex)
		items = append(items, item)
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

// renderCourseItem renders a single course item
func (m CoursesScreen) renderCourseItem(course models.Course, isSelected bool) string {
	// Cursor
	cursor := "  "
	if isSelected {
		cursor = "► "
	}

	// Course code and name
	var codeName string
	if course.Code != "" {
		codeName = fmt.Sprintf("%s - %s", course.Code, course.Name)
	} else {
		codeName = course.Name
	}

	// Course color indicator
	colorIndicator := "●"
	if course.Color != "" {
		colorIndicator = lipgloss.NewStyle().
			Foreground(lipgloss.Color(course.Color)).
			Render("●")
	}

	// Details line
	var details []string
	if course.Professor != "" {
		details = append(details, course.Professor)
	}
	if course.Location != "" {
		details = append(details, course.Location)
	}
	if course.Credits > 0 {
		details = append(details, fmt.Sprintf("%d credits", course.Credits))
	}
	detailsStr := strings.Join(details, " • ")

	// Schedule summary
	scheduleSummary := ""
	if len(course.Schedule) > 0 {
		var days []string
		for _, sched := range course.Schedule {
			days = append(days, sched.DayOfWeekShort())
		}
		scheduleSummary = fmt.Sprintf("📅 %s %s-%s",
			strings.Join(days, "/"),
			course.Schedule[0].StartTime,
			course.Schedule[0].EndTime,
		)
	}

	// Style
	nameStyle := lipgloss.NewStyle().Bold(true)
	detailStyle := lipgloss.NewStyle().Foreground(styles.Muted)

	if isSelected {
		nameStyle = nameStyle.Foreground(styles.Primary)
	}

	// Build the item
	nameLineLine := cursor + colorIndicator + " " + nameStyle.Render(codeName)
	detailsLine := "   " + detailStyle.Render(detailsStr)
	scheduleLine := "   " + detailStyle.Render(scheduleSummary)

	var lines []string
	lines = append(lines, nameLineLine)
	if detailsStr != "" {
		lines = append(lines, detailsLine)
	}
	if scheduleSummary != "" {
		lines = append(lines, scheduleLine)
	}

	return lipgloss.NewStyle().
		Padding(0, 2).
		Render(strings.Join(lines, "\n"))
}

// renderShortcuts renders the keyboard shortcuts
func (m CoursesScreen) renderShortcuts() string {
	shortcuts := []string{
		styles.Shortcut.Render("n") + styles.ShortcutText.Render(" new"),
		styles.Shortcut.Render("e") + styles.ShortcutText.Render(" edit"),
		styles.Shortcut.Render("d") + styles.ShortcutText.Render(" delete"),
		styles.Shortcut.Render("enter") + styles.ShortcutText.Render(" view"),
		styles.Shortcut.Render("esc") + styles.ShortcutText.Render(" back"),
	}

	return lipgloss.NewStyle().
		Padding(1, 0).
		Render(strings.Join(shortcuts, "  "))
}

// renderDeleteConfirmation renders the delete confirmation dialog
func (m CoursesScreen) renderDeleteConfirmation() string {
	if m.selectedIndex < 0 || m.selectedIndex >= len(m.courses) {
		return ""
	}

	course := m.courses[m.selectedIndex]
	question := fmt.Sprintf("Delete course \"%s\"?", course.Name)

	dialog := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Danger).
		Padding(1, 2).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
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

	// Build the complete view with base content
	baseView := m.renderCourseList()
	
	// Overlay the dialog on top
	return baseView + "\n\n" + dialog
}

// Commands

type fetchCoursesMsg struct {
	courses []models.Course
	err     error
}

func (m CoursesScreen) fetchCoursesCmd() tea.Cmd {
	return func() tea.Msg {
		courses, err := m.db.Courses().GetAll()
		return fetchCoursesMsg{
			courses: courses,
			err:     err,
		}
	}
}

func (m CoursesScreen) deleteCourseCmd(id string) tea.Cmd {
	return func() tea.Msg {
		if err := m.db.Courses().Delete(id); err != nil {
			return fetchCoursesMsg{err: err}
		}
		courses, err := m.db.Courses().GetAll()
		return fetchCoursesMsg{
			courses: courses,
			err:     err,
		}
	}
}

// IsCourseFormActive returns true if the course form is currently active
func (m CoursesScreen) IsCourseFormActive() bool {
	return m.showForm
}
