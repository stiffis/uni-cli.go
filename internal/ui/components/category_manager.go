package components

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stiffis/UniCLI/internal/database"
	"github.com/stiffis/UniCLI/internal/models"
	"github.com/stiffis/UniCLI/internal/ui/styles"
)

var (
	docStyle = lipgloss.NewStyle()
)

type managerMode int

const (
	modeList managerMode = iota
	modeForm
)

type keyMap struct {
	New    key.Binding
	Edit   key.Binding
	Delete key.Binding
	Quit   key.Binding
	Esc    key.Binding
}

var keys = keyMap{
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "quit"),
	),
}

// itemDelegate implements list.ItemDelegate for our custom category rendering.
type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	cat, ok := listItem.(models.Category)
	if !ok {
		return
	}

	// "University "
	icon := ""
	name := cat.Name

	// Style the icon with the category's color
	iconStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(cat.Color))
	styledIcon := iconStyle.Render(icon)

	str := fmt.Sprintf("%s %s", name, styledIcon)

	if index == m.Index() {
		// Selected item
		fmt.Fprint(w, lipgloss.NewStyle().
			Background(styles.SelectedBackground).
			Foreground(styles.SelectedForeground).
			Render(str))
	} else {
		// Normal item
		fmt.Fprint(w, str)
	}
}

type CategoryManager struct {
	list     list.Model
	form     CategoryForm
	db       *database.DB
	mode     managerMode
	quitting bool
	delegate itemDelegate
}

func NewCategoryManager(db *database.DB) *CategoryManager {
	items := []list.Item{}

	delegate := itemDelegate{}
	l := list.New(items, delegate, 0, 0)
	l.Title = "Category Manager"
	l.SetShowHelp(false)

	return &CategoryManager{
		list:     l,
		db:       db,
		mode:     modeList,
		delegate: delegate,
	}
}

func (m *CategoryManager) Init() tea.Cmd {
	return m.fetchCategories
}

type CategoriesMsg []models.Category
type CategoryMsg models.Category
type errMsg error

func (m *CategoryManager) fetchCategories() tea.Msg {
	cats, err := m.db.Categories().FindAll()
	if err != nil {
		return errMsg(err)
	}
	return CategoriesMsg(cats)
}

func (m *CategoryManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case modeForm:
		return m.updateForm(msg)
	default: // modeList
		return m.updateList(msg)
	}
}

func (m *CategoryManager) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-1)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit), key.Matches(msg, keys.Esc):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, keys.New):
			m.mode = modeForm
			m.form = NewCategoryForm(nil)
			return m, m.form.Init()
		case key.Matches(msg, keys.Edit):
			if item, ok := m.list.SelectedItem().(models.Category); ok {
				m.mode = modeForm
				m.form = NewCategoryForm(&item)
				return m, m.form.Init()
			}
		case key.Matches(msg, keys.Delete):
			if item, ok := m.list.SelectedItem().(models.Category); ok {
				err := m.db.Categories().Delete(item.ID)
				if err != nil {
					// handle error
				}
				return m, m.fetchCategories
			}
		}

	case CategoriesMsg:
		items := make([]list.Item, len(msg))
		for i, cat := range msg {
			items[i] = cat
		}
		m.list.SetItems(items)
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *CategoryManager) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.form, cmd = m.form.Update(msg)

	if m.form.IsCancelled() {
		m.mode = modeList
		return m, nil
	}

	if m.form.IsSubmitted() {
		category := m.form.GetCategory()
		var err error
		if m.form.IsNewCategory() {
			err = m.db.Categories().Create(category)
		} else {
			err = m.db.Categories().Update(category)
		}

		if err != nil {
			// handle error
		}

		m.mode = modeList
		return m, m.fetchCategories
	}

	return m, cmd
}

func (m *CategoryManager) renderShortcuts() string {
	var shortcuts []string
	shortcuts = []string{
		styles.Shortcut.Render("n") + styles.ShortcutText.Render(" new"),
		styles.Shortcut.Render("e") + styles.ShortcutText.Render(" edit"),
		styles.Shortcut.Render("d") + styles.ShortcutText.Render(" delete"),
		styles.Shortcut.Render("q") + styles.ShortcutText.Render(" quit"),
	}
	shortcutLine := strings.Join(shortcuts, "  ")
	return lipgloss.NewStyle().Render(shortcutLine)
}

func (m *CategoryManager) View() string {
	if m.quitting {
		return ""
	}

	switch m.mode {
	case modeForm:
		return m.form.View()
	default:
		var b strings.Builder
		for i, item := range m.list.Items() {
			m.delegate.Render(&b, m.list, i, item)
			if i < len(m.list.Items())-1 {
				b.WriteString("\n")
			}
		}

		boxStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Primary).
			Padding(3, 6)

		return lipgloss.JoinVertical(lipgloss.Left, boxStyle.Render(b.String()), m.renderShortcuts())
	}
}

func (m *CategoryManager) IsQuitting() bool {
	return m.quitting
}

func (m *CategoryManager) SetSize(width, height int) {
	m.list.SetSize(width, height)
}

func (m *CategoryManager) Reset() {
	m.quitting = false
}

