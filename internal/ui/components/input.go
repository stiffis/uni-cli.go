package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stiffis/UniCLI/internal/ui/styles"
)

// Input is a styled text input component
type Input struct {
	textInput textinput.Model
	label     string
	width     int
}

// NewInput creates a new input field
func NewInput(label string, placeholder string) Input {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 100
	ti.Width = 50

	return Input{
		textInput: ti,
		label:     label,
		width:     50,
	}
}

// Focus focuses the input
func (i *Input) Focus() tea.Cmd {
	return i.textInput.Focus()
}

// Blur removes focus from the input
func (i *Input) Blur() {
	i.textInput.Blur()
}

// SetValue sets the input value
func (i *Input) SetValue(value string) {
	i.textInput.SetValue(value)
}

// Value returns the input value
func (i *Input) Value() string {
	return i.textInput.Value()
}

func (i *Input) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	i.textInput, cmd = i.textInput.Update(msg)
	return cmd
}

func (i Input) View() string {
	labelStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true)

	label := labelStyle.Render(i.label)
	input := i.textInput.View()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		label,
		input,
	)
}
