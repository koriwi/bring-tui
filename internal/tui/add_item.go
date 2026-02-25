package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type addItemModel struct {
	nameInput    textinput.Model
	specInput    textinput.Model
	focusName    bool
	itemName     string
	spec         string
	submitted    bool
	cancelled    bool
	title        string
	originalName string
}

func newAddItem() *addItemModel {
	name := textinput.New()
	name.Placeholder = "Item name (e.g. Milch)"
	name.Focus()
	name.CharLimit = 100

	spec := textinput.New()
	spec.Placeholder = "Description (optional, e.g. 1.5%)"
	spec.CharLimit = 200

	return &addItemModel{
		nameInput: name,
		specInput: spec,
		focusName: true,
		title:     "Add Item",
	}
}

func newEditItem(name, spec string) *addItemModel {
	nameInput := textinput.New()
	nameInput.Placeholder = "Item name (e.g. Milch)"
	nameInput.SetValue(name)
	nameInput.Focus()
	nameInput.CharLimit = 100

	specInput := textinput.New()
	specInput.Placeholder = "Description (optional, e.g. 1.5%)"
	specInput.SetValue(spec)
	specInput.CharLimit = 200

	return &addItemModel{
		nameInput:    nameInput,
		specInput:    specInput,
		focusName:    true,
		title:        "Edit Item",
		originalName: name,
	}
}

func (m *addItemModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *addItemModel) Update(msg tea.Msg) (*addItemModel, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, keys.Escape):
			m.cancelled = true
			return m, nil

		case key.Matches(keyMsg, keys.Tab):
			if m.focusName {
				m.focusName = false
				m.nameInput.Blur()
				m.specInput.Focus()
			} else {
				m.focusName = true
				m.specInput.Blur()
				m.nameInput.Focus()
			}
			return m, nil

		case key.Matches(keyMsg, keys.Enter):
			if m.focusName && m.nameInput.Value() != "" {
				// Move to spec field
				m.focusName = false
				m.nameInput.Blur()
				m.specInput.Focus()
				return m, nil
			}
			if !m.focusName && m.nameInput.Value() != "" {
				m.itemName = m.nameInput.Value()
				m.spec = m.specInput.Value()
				m.submitted = true
				return m, nil
			}
			// If name is filled and user presses enter on spec
			if m.nameInput.Value() != "" {
				m.itemName = m.nameInput.Value()
				m.spec = m.specInput.Value()
				m.submitted = true
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	if m.focusName {
		m.nameInput, cmd = m.nameInput.Update(msg)
	} else {
		m.specInput, cmd = m.specInput.Update(msg)
	}
	return m, cmd
}

func (m *addItemModel) View() string {
	s := "\n" + titleStyle.Render(m.title) + "\n\n"
	s += "  Item:  " + m.nameInput.View() + "\n"
	s += "  Desc:  " + m.specInput.View() + "\n\n"
	s += helpDescStyle.Render("  Tab: switch field · Enter: confirm · Esc: cancel") + "\n"
	return s
}
