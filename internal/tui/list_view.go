package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/key"
	"github.com/paulleonhardhellweg/bring-tui/internal/bring"
)

type listViewModel struct {
	items    []bring.Item
	recently []bring.Item
	listName string
	cursor   int
	width    int
	height   int
}

func newListView(resp *bring.ItemsResponse, listName string, width, height int) *listViewModel {
	if listName == "" {
		listName = "Shopping List"
	}
	return &listViewModel{
		items:    resp.Purchase,
		recently: resp.Recently,
		listName: listName,
		cursor:   0,
		width:    width,
		height:   height,
	}
}

func (m *listViewModel) selectedItem() *bring.Item {
	if len(m.items) == 0 || m.cursor >= len(m.items) {
		return nil
	}
	return &m.items[m.cursor]
}

func (m *listViewModel) removeItem(itemID string) {
	for i, item := range m.items {
		if item.ItemID == itemID {
			m.items = append(m.items[:i], m.items[i+1:]...)
			if m.cursor >= len(m.items) && m.cursor > 0 {
				m.cursor--
			}
			return
		}
	}
}

func (m *listViewModel) addItem(itemID, spec string) {
	m.items = append([]bring.Item{{ItemID: itemID, Spec: spec}}, m.items...)
	m.cursor = 0
}

func (m *listViewModel) Update(msg tea.Msg) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(keyMsg, keys.Down):
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		}
	}
}

func (m *listViewModel) View() string {
	s := "\n" + titleStyle.Render(fmt.Sprintf("🛒 %s", m.listName)) + "\n"
	s += dividerStyle.Render("────────────────────────────────") + "\n"

	if len(m.items) == 0 {
		s += sectionStyle.Render("  List is empty. Press 'a' to add items.") + "\n"
	}

	for i, item := range m.items {
		cursor := "  "
		style := itemStyle
		if i == m.cursor {
			cursor = "▸ "
			style = selectedItemStyle
		}

		line := fmt.Sprintf("%s● %s", cursor, item.ItemID)
		if item.Spec != "" {
			line += " " + specStyle.Render("— "+item.Spec)
		}
		s += style.Render(line) + "\n"
	}

	if len(m.recently) > 0 {
		s += "\n" + sectionStyle.Render("  ── recently bought ──") + "\n"
		max := 5
		if len(m.recently) < max {
			max = len(m.recently)
		}
		for _, item := range m.recently[:max] {
			s += doneItemStyle.Render(fmt.Sprintf("  ✓ %s", item.ItemID)) + "\n"
		}
		if len(m.recently) > 5 {
			s += sectionStyle.Render(fmt.Sprintf("  ... and %d more", len(m.recently)-5)) + "\n"
		}
	}

	return s
}
