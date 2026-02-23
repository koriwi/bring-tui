package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/paulleonhardhellweg/bring-tui/internal/bring"
)

type listPickerModel struct {
	lists     []bring.List
	cursor    int
	selected  *bring.List
	cancelled bool
	currentID string
}

func newListPicker(lists []bring.List, currentListUUID string) *listPickerModel {
	return &listPickerModel{
		lists:     lists,
		cursor:    0,
		currentID: currentListUUID,
	}
}

func (m *listPickerModel) Update(msg tea.Msg) (*listPickerModel, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(keyMsg, keys.Down):
			if m.cursor < len(m.lists)-1 {
				m.cursor++
			}
		case key.Matches(keyMsg, keys.Enter):
			if len(m.lists) > 0 {
				m.selected = &m.lists[m.cursor]
			}
		case key.Matches(keyMsg, keys.Escape), keyMsg.String() == "q":
			m.cancelled = true
		}
	}
	return m, nil
}

func (m *listPickerModel) View() string {
	s := "\n" + titleStyle.Render("📋 Select a list") + "\n"
	s += dividerStyle.Render("────────────────────────────────") + "\n"

	for i, list := range m.lists {
		cursor := "  "
		if i == m.cursor {
			cursor = "▸ "
		}

		current := ""
		if list.ListUUID == m.currentID {
			current = " (current)"
		}

		style := itemStyle
		if i == m.cursor {
			style = selectedItemStyle
		}

		s += style.Render(fmt.Sprintf("%s%s%s", cursor, list.Name, current)) + "\n"
	}

	s += "\n" + helpDescStyle.Render("  Enter: select · Esc: cancel") + "\n"
	return s
}
