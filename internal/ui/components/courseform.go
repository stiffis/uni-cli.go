package components

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stiffis/UniCLI/internal/database"
	"github.com/stiffis/UniCLI/internal/models"
	"github.com/stiffis/UniCLI/internal/ui/styles"
)

// FormSubmitMsg is sent when a form is successfully submitted
type FormSubmitMsg struct{}

// CourseForm represents the course creation/edit form
type CourseForm struct {
	db            *database.DB
	course        *models.Course
	inputs        []textinput.Model
	focusedInput  int
	isEdit        bool
	width         int
	height        int
	err           string
	scheduleInput string // Temporary storage for schedule input
}

const (
	courseInputName = iota
	courseInputCode
	courseInputProfessor
	courseInputLocation
	courseInputSemester
	courseInputCredits
	courseInputColor
	courseInputSchedule
	courseInputDescription
)

// NewCourseForm creates a new course form
func NewCourseForm(db *database.DB, course *models.Course) *CourseForm {
	inputs := make([]textinput.Model, 9)

	// Name
	inputs[courseInputName] = textinput.New()
	inputs[courseInputName].Placeholder = "e.g., Calculus I"
	inputs[courseInputName].Focus()
	inputs[courseInputName].Width = 50

	// Code
	inputs[courseInputCode] = textinput.New()
	inputs[courseInputCode].Placeholder = "e.g., MATH 101"
	inputs[courseInputCode].Width = 50

	// Professor
	inputs[courseInputProfessor] = textinput.New()
	inputs[courseInputProfessor].Placeholder = "e.g., Dr. Smith"
	inputs[courseInputProfessor].Width = 50

	// Location
	inputs[courseInputLocation] = textinput.New()
	inputs[courseInputLocation].Placeholder = "e.g., Room 301"
	inputs[courseInputLocation].Width = 50

	// Semester
	inputs[courseInputSemester] = textinput.New()
	inputs[courseInputSemester].Placeholder = "e.g., Fall 2025"
	inputs[courseInputSemester].Width = 50

	// Credits
	inputs[courseInputCredits] = textinput.New()
	inputs[courseInputCredits].Placeholder = "e.g., 3"
	inputs[courseInputCredits].Width = 50

	// Color
	inputs[courseInputColor] = textinput.New()
	inputs[courseInputColor].Placeholder = "e.g., #7E9CD8"
	inputs[courseInputColor].Width = 50

	// Schedule (simplified input)
	inputs[courseInputSchedule] = textinput.New()
	inputs[courseInputSchedule].Placeholder = "e.g., Mon/Wed/Fri 09:00-10:30"
	inputs[courseInputSchedule].Width = 50

	// Description
	inputs[courseInputDescription] = textinput.New()
	inputs[courseInputDescription].Placeholder = "Brief description"
	inputs[courseInputDescription].Width = 50

	isEdit := course != nil

	if isEdit {
		inputs[courseInputName].SetValue(course.Name)
		inputs[courseInputCode].SetValue(course.Code)
		inputs[courseInputProfessor].SetValue(course.Professor)
		inputs[courseInputLocation].SetValue(course.Location)
		inputs[courseInputSemester].SetValue(course.Semester)
		if course.Credits > 0 {
			inputs[courseInputCredits].SetValue(strconv.Itoa(course.Credits))
		}
		inputs[courseInputColor].SetValue(course.Color)
		inputs[courseInputDescription].SetValue(course.Description)

		if len(course.Schedule) > 0 {
			schedStr := formatScheduleForDisplay(course.Schedule)
			inputs[courseInputSchedule].SetValue(schedStr)
		}
	}

	return &CourseForm{
		db:           db,
		course:       course,
		inputs:       inputs,
		focusedInput: 0,
		isEdit:       isEdit,
	}
}

// Init initializes the form
func (f *CourseForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *CourseForm) Update(msg tea.Msg) (*CourseForm, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.width = msg.Width
		f.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return f, tea.Quit
		case "esc":
			return f, nil
		case "tab", "down":
			f.nextInput()
		case "shift+tab", "up":
			f.prevInput()
		case "enter":
			if f.focusedInput == len(f.inputs)-1 {
				// Submit form
				return f, f.submitForm()
			}
			f.nextInput()
		case "ctrl+s":
			// Submit form
			return f, f.submitForm()
		}
	}

	var cmd tea.Cmd
	f.inputs[f.focusedInput], cmd = f.inputs[f.focusedInput].Update(msg)
	cmds = append(cmds, cmd)

	return f, tea.Batch(cmds...)
}

func (f *CourseForm) View() string {
	title := "New Course"
	if f.isEdit {
		title = "Edit Course"
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		Padding(1, 0)

	labelStyle := lipgloss.NewStyle().
		Foreground(styles.Accent).
		Width(15)

	var fields []string

	// Name
	fields = append(fields, fmt.Sprintf("%s %s",
		labelStyle.Render("Name:"),
		f.inputs[courseInputName].View(),
	))

	// Code
	fields = append(fields, fmt.Sprintf("%s %s",
		labelStyle.Render("Code:"),
		f.inputs[courseInputCode].View(),
	))

	// Professor
	fields = append(fields, fmt.Sprintf("%s %s",
		labelStyle.Render("Professor:"),
		f.inputs[courseInputProfessor].View(),
	))

	// Location
	fields = append(fields, fmt.Sprintf("%s %s",
		labelStyle.Render("Location:"),
		f.inputs[courseInputLocation].View(),
	))

	// Semester
	fields = append(fields, fmt.Sprintf("%s %s",
		labelStyle.Render("Semester:"),
		f.inputs[courseInputSemester].View(),
	))

	// Credits
	fields = append(fields, fmt.Sprintf("%s %s",
		labelStyle.Render("Credits:"),
		f.inputs[courseInputCredits].View(),
	))

	// Color
	fields = append(fields, fmt.Sprintf("%s %s",
		labelStyle.Render("Color:"),
		f.inputs[courseInputColor].View(),
	))

	// Schedule
	fields = append(fields, fmt.Sprintf("%s %s",
		labelStyle.Render("Schedule:"),
		f.inputs[courseInputSchedule].View(),
	))

	// Description
	fields = append(fields, fmt.Sprintf("%s %s",
		labelStyle.Render("Description:"),
		f.inputs[courseInputDescription].View(),
	))

	formContent := strings.Join(fields, "\n\n")

	// Error message
	errorMsg := ""
	if f.err != "" {
		errorMsg = lipgloss.NewStyle().
			Foreground(styles.Warning).
			Render("\n⚠ " + f.err)
	}

	// Help text
	help := lipgloss.NewStyle().
		Foreground(styles.Muted).
		Render("\nCtrl+S to save • Esc to cancel • Tab/Shift+Tab to navigate")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(title),
		formContent,
		errorMsg,
		help,
	)
}

// nextInput focuses the next input
func (f *CourseForm) nextInput() {
	f.inputs[f.focusedInput].Blur()
	f.focusedInput = (f.focusedInput + 1) % len(f.inputs)
	f.inputs[f.focusedInput].Focus()
}

// prevInput focuses the previous input
func (f *CourseForm) prevInput() {
	f.inputs[f.focusedInput].Blur()
	f.focusedInput--
	if f.focusedInput < 0 {
		f.focusedInput = len(f.inputs) - 1
	}
	f.inputs[f.focusedInput].Focus()
}

// submitForm saves the course
func (f *CourseForm) submitForm() tea.Cmd {
	return func() tea.Msg {
		// Validate
		if f.inputs[courseInputName].Value() == "" {
			f.err = "Name is required"
			return nil
		}

		credits := 0
		if f.inputs[courseInputCredits].Value() != "" {
			var err error
			credits, err = strconv.Atoi(f.inputs[courseInputCredits].Value())
			if err != nil {
				f.err = "Credits must be a number"
				return nil
			}
		}

		schedules, err := parseScheduleInput(f.inputs[courseInputSchedule].Value())
		if err != nil {
			f.err = "Invalid schedule format. Use: Mon/Wed/Fri 09:00-10:30"
			return nil
		}

		var course *models.Course
		if f.isEdit {
			course = f.course
		} else {
			course = models.NewCourse(f.inputs[courseInputName].Value())
		}

		course.Name = f.inputs[courseInputName].Value()
		course.Code = f.inputs[courseInputCode].Value()
		course.Professor = f.inputs[courseInputProfessor].Value()
		course.Location = f.inputs[courseInputLocation].Value()
		course.Semester = f.inputs[courseInputSemester].Value()
		course.Credits = credits
		course.Color = f.inputs[courseInputColor].Value()
		course.Description = f.inputs[courseInputDescription].Value()
		course.Schedule = schedules

		var saveErr error
		if f.isEdit {
			saveErr = f.db.Courses().Update(course)
		} else {
			saveErr = f.db.Courses().Create(course)
		}

		if saveErr != nil {
			f.err = saveErr.Error()
			return nil
		}

		return FormSubmitMsg{}
	}
}

// Helper functions

func parseScheduleInput(input string) ([]models.CourseSchedule, error) {
	if input == "" {
		return []models.CourseSchedule{}, nil
	}

	// Expected format: "Mon/Wed/Fri 09:00-10:30"
	parts := strings.Fields(input)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format")
	}

	days := strings.Split(parts[0], "/")
	times := strings.Split(parts[1], "-")
	if len(times) != 2 {
		return nil, fmt.Errorf("invalid time format")
	}

	startTime := times[0]
	endTime := times[1]

	dayMap := map[string]int{
		"Mon": 1, "Monday": 1,
		"Tue": 2, "Tuesday": 2,
		"Wed": 3, "Wednesday": 3,
		"Thu": 4, "Thursday": 4,
		"Fri": 5, "Friday": 5,
		"Sat": 6, "Saturday": 6,
		"Sun": 7, "Sunday": 7,
	}

	var schedules []models.CourseSchedule
	for _, day := range days {
		dayNum, ok := dayMap[day]
		if !ok {
			return nil, fmt.Errorf("invalid day: %s", day)
		}

		schedule := models.NewCourseSchedule(
			"", // Will be set when saving
			dayNum,
			startTime,
			endTime,
		)
		schedules = append(schedules, *schedule)
	}

	return schedules, nil
}

func formatScheduleForDisplay(schedules []models.CourseSchedule) string {
	if len(schedules) == 0 {
		return ""
	}

	var days []string
	for _, sched := range schedules {
		days = append(days, sched.DayOfWeekShort())
	}

	return fmt.Sprintf("%s %s-%s",
		strings.Join(days, "/"),
		schedules[0].StartTime,
		schedules[0].EndTime,
	)
}
