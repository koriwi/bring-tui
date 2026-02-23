package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Add      key.Binding
	Done     key.Binding
	Remove   key.Binding
	Lists    key.Binding
	Reload   key.Binding
	Quit     key.Binding
	Enter    key.Binding
	Escape   key.Binding
	Tab      key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	Done: key.NewBinding(
		key.WithKeys("d", "enter"),
		key.WithHelp("d", "done"),
	),
	Remove: key.NewBinding(
		key.WithKeys("x", "delete"),
		key.WithHelp("x", "remove"),
	),
	Lists: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "lists"),
	),
	Reload: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reload"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
	),
}

func helpBar() string {
	items := []struct{ key, desc string }{
		{"a", "add"},
		{"d", "done"},
		{"x", "remove"},
		{"enter", "re-add"},
		{"l", "lists"},
		{"r", "reload"},
		{"q", "quit"},
	}
	var s string
	for i, item := range items {
		if i > 0 {
			s += helpDescStyle.Render(" · ")
		}
		s += helpKeyStyle.Render(item.key) + " " + helpDescStyle.Render(item.desc)
	}
	return s
}
