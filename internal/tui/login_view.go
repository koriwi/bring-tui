package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type loginViewModel struct {
	emailInput    textinput.Model
	passwordInput textinput.Model
	focusEmail    bool
	email         string
	password      string
	submitted     bool
	cancelled     bool
	err           string
}

func newLoginView() *loginViewModel {
	emailInput := textinput.New()
	emailInput.Placeholder = "Email"
	emailInput.Focus()
	emailInput.CharLimit = 100

	passwordInput := textinput.New()
	passwordInput.Placeholder = "Password"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.CharLimit = 100

	return &loginViewModel{
		emailInput:    emailInput,
		passwordInput: passwordInput,
		focusEmail:    true,
	}
}

func (m *loginViewModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *loginViewModel) Update(msg tea.Msg) (*loginViewModel, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, keys.Escape):
			m.cancelled = true
			return m, nil

		case key.Matches(keyMsg, keys.Tab), key.Matches(keyMsg, keys.Enter):
			if m.focusEmail {
				m.focusEmail = false
				m.emailInput.Blur()
				m.passwordInput.Focus()
				return m, nil
			}
			// On password field, enter submits
			if key.Matches(keyMsg, keys.Enter) && m.emailInput.Value() != "" && m.passwordInput.Value() != "" {
				m.email = m.emailInput.Value()
				m.password = m.passwordInput.Value()
				m.submitted = true
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	if m.focusEmail {
		m.emailInput, cmd = m.emailInput.Update(msg)
	} else {
		m.passwordInput, cmd = m.passwordInput.Update(msg)
	}
	return m, cmd
}

func (m *loginViewModel) View() string {
	s := "\n" + titleStyle.Render("🔐 Bring! Login") + "\n\n"
	s += "  Email:    " + m.emailInput.View() + "\n"
	s += "  Password: " + m.passwordInput.View() + "\n\n"
	if m.err != "" {
		s += errorStyle.Render("  "+m.err) + "\n\n"
	}
	s += helpDescStyle.Render("  Tab/Enter: next field · Enter: login · Esc: quit") + "\n"
	return s
}
