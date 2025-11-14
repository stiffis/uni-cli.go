package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stiffis/UniCLI/internal/models"
	"github.com/stiffis/UniCLI/internal/ui/styles"
)

type CategoryForm struct {
	originalCategory *models.Category
	categoryID       string
	nameInput        Input
	colorInput       Input
	focusedField     int
	submitted        bool
	cancelled        bool
	width            int
	height           int
}

const (
	categoryFieldName = iota
	categoryFieldColor
	categoryFieldButtons
)

func NewCategoryForm(category *models.Category) CategoryForm {
	nameInput := NewInput("Name:", "Enter category name...")
	colorInput := NewInput("Color:", "Enter hex color (e.g., #FF00FF)...")

	form := CategoryForm{
		nameInput:    nameInput,
		colorInput:   colorInput,
		focusedField: categoryFieldName,
		width:        60,
		height:       10,
	}

	if category != nil {
		form.originalCategory = category
		form.categoryID = category.ID
		form.nameInput.SetValue(category.Name)
		form.colorInput.SetValue(category.Color)
	}

	form.nameInput.Focus()

	return form
}

func (f CategoryForm) Init() tea.Cmd {
	return nil
}

func (f CategoryForm) Update(msg tea.Msg) (CategoryForm, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			f.cancelled = true
			return f, nil

		case "tab", "down":
			f.blurAll()
			f.focusedField = (f.focusedField + 1) % 3
			cmd = f.focusField(f.focusedField)
			return f, cmd

		case "shift+tab", "up":
			f.blurAll()
			f.focusedField = (f.focusedField + 2) % 3
			cmd = f.focusField(f.focusedField)
			return f, cmd

		case "enter":
			if f.focusedField == categoryFieldButtons {
				if f.nameInput.Value() != "" && f.colorInput.Value() != "" {
					f.submitted = true
				}
				return f, nil
			}
		}
	}

	switch f.focusedField {
	case categoryFieldName:
		cmd = f.nameInput.Update(msg)
	case categoryFieldColor:
		cmd = f.colorInput.Update(msg)
	}

	return f, cmd
}

func (f CategoryForm) View() string {
	var sections []string

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		Align(lipgloss.Center).
		Width(f.width).
		Render(" New Category")
	sections = append(sections, title)
	sections = append(sections, "")

	sections = append(sections, f.nameInput.View())
	sections = append(sections, "")

	sections = append(sections, f.colorInput.View())
	sections = append(sections, "")

	sections = append(sections, f.renderButtons())
	sections = append(sections, "")

	helpStyle := lipgloss.NewStyle().
		Foreground(styles.Muted).
		Italic(true)

	help := helpStyle.Render("Tab: next field  |  Esc: cancel  |  Enter: submit")
	sections = append(sections, help)

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Primary).
		Padding(1, 2).
		Width(f.width)

	return modalStyle.Render(content)
}

func (f CategoryForm) renderButtons() string {
	var submitText string
	if f.categoryID != "" {
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

	if f.focusedField == categoryFieldButtons {
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

func (f *CategoryForm) blurAll() {
	f.nameInput.Blur()
	f.colorInput.Blur()
}

func (f *CategoryForm) focusField(field int) tea.Cmd {
	switch field {
	case categoryFieldName:
		return f.nameInput.Focus()
	case categoryFieldColor:
		return f.colorInput.Focus()
	}
	return nil
}

func (f CategoryForm) GetCategory() *models.Category {
	var category *models.Category
	if f.originalCategory != nil {
		category = f.originalCategory
	} else {
		category = models.NewCategory("", "")
	}

	category.Name = strings.TrimSpace(f.nameInput.Value())
	category.Color = strings.TrimSpace(f.colorInput.Value())

	return category
}

func (f CategoryForm) IsSubmitted() bool {
	return f.submitted
}

func (f CategoryForm) IsCancelled() bool {
	return f.cancelled
}

func (f CategoryForm) IsNewCategory() bool {
	return f.categoryID == ""
}
