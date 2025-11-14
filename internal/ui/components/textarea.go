package components

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stiffis/UniCLI/internal/ui/styles"
)

// TextArea is a multi-line text input
type TextArea struct {
	textarea textarea.Model
	label    string
}

// NewTextArea creates a new text area
func NewTextArea(label string, placeholder string) TextArea {
	ta := textarea.New()
	ta.Placeholder = placeholder
	ta.CharLimit = 500
	ta.SetWidth(50)
	ta.SetHeight(3)
	ta.ShowLineNumbers = false

	return TextArea{
		textarea: ta,
		label:    label,
	}
}

// Focus focuses the textarea
func (t *TextArea) Focus() tea.Cmd {
	return t.textarea.Focus()
}

// Blur removes focus
func (t *TextArea) Blur() {
	t.textarea.Blur()
}

// SetValue sets the textarea value
func (t *TextArea) SetValue(value string) {
	t.textarea.SetValue(value)
}

// Value returns the textarea value
func (t *TextArea) Value() string {
	return t.textarea.Value()
}

func (t *TextArea) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	t.textarea, cmd = t.textarea.Update(msg)
	return cmd
}

func (t TextArea) View() string {
	labelStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true)

	label := labelStyle.Render(t.label)
	area := t.textarea.View()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		label,
		area,
	)
}
