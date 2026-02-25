package tui

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paulleonhardhellweg/bring-tui/internal/bring"
	"github.com/paulleonhardhellweg/bring-tui/internal/config"
)

type state int

const (
	stateLoading state = iota
	stateLogin
	stateList
	stateAddItem
	stateEditItem
	stateListPicker
)

// Messages
type authSuccessMsg struct {
	client *bring.Client
	stored *config.StoredAuth
}

type authErrorMsg struct{ err error }

type itemsLoadedMsg struct {
	items *bring.ItemsResponse
}

type itemsErrorMsg struct{ err error }

type itemAddedMsg struct{ item, spec string }
type itemDoneMsg struct{ item string }
type itemRemovedMsg struct{ item string }
type apiErrorMsg struct{ err error }

type listsLoadedMsg struct {
	lists []bring.List
}

type statusMsg struct {
	text    string
	isError bool
}

// App is the root TUI model
type App struct {
	state      state
	spinner    spinner.Model
	client     *bring.Client
	stored     *config.StoredAuth
	listView   *listViewModel
	addItem    *addItemModel
	editItem   *addItemModel
	loginView  *loginViewModel
	listPicker *listPickerModel
	status     string
	statusErr  bool
	width      int
	height     int
}

func NewApp() *App {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(colorPrimary)

	return &App{
		state:   stateLoading,
		spinner: s,
	}
}

func (a *App) Init() tea.Cmd {
	return tea.Batch(a.spinner.Tick, a.tryAuth())
}

func (a *App) tryAuth() tea.Cmd {
	return func() tea.Msg {
		client, stored, err := bring.Authenticate()
		if err != nil {
			return authErrorMsg{err: err}
		}
		return authSuccessMsg{client: client, stored: stored}
	}
}

func (a *App) loadItems() tea.Cmd {
	return func() tea.Msg {
		items, err := a.client.GetItems(a.stored.DefaultListUUID)
		if err != nil {
			if errors.Is(err, bring.ErrAuthExpired) {
				return authErrorMsg{err: err}
			}
			return itemsErrorMsg{err: err}
		}
		return itemsLoadedMsg{items: items}
	}
}

func (a *App) loadLists() tea.Cmd {
	return func() tea.Msg {
		lists, err := a.client.GetLists()
		if err != nil {
			if errors.Is(err, bring.ErrAuthExpired) {
				return authErrorMsg{err: err}
			}
			return itemsErrorMsg{err: err}
		}
		return listsLoadedMsg{lists: lists}
	}
}

// apiErr converts an API error to the appropriate tea.Msg,
// redirecting to the login screen if the session has expired.
func apiErr(err error, msg string) tea.Msg {
	if errors.Is(err, bring.ErrAuthExpired) {
		return authErrorMsg{err: err}
	}
	return statusMsg{text: msg, isError: true}
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		if a.listView != nil {
			a.listView.width = msg.Width
			a.listView.height = msg.Height
		}
		return a, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return a, tea.Quit
		}

	case authSuccessMsg:
		a.client = msg.client
		a.stored = msg.stored
		return a, a.loadItems()

	case authErrorMsg:
		a.state = stateLogin
		a.loginView = newLoginView()
		return a, a.loginView.Init()

	case itemsLoadedMsg:
		a.state = stateList
		a.listView = newListView(msg.items, a.stored.DefaultListName, a.width, a.height)
		a.status = ""
		return a, nil

	case itemsErrorMsg:
		a.status = fmt.Sprintf("Error: %v", msg.err)
		a.statusErr = true
		return a, nil

	case listsLoadedMsg:
		a.state = stateListPicker
		a.listPicker = newListPicker(msg.lists, a.stored.DefaultListUUID)
		return a, nil

	case statusMsg:
		a.status = msg.text
		a.statusErr = msg.isError
		return a, nil
	}

	// Delegate to sub-views
	switch a.state {
	case stateLoading:
		var cmd tea.Cmd
		a.spinner, cmd = a.spinner.Update(msg)
		return a, cmd

	case stateLogin:
		return a.updateLogin(msg)

	case stateList:
		return a.updateList(msg)

	case stateAddItem:
		return a.updateAddItem(msg)

	case stateEditItem:
		return a.updateEditItem(msg)

	case stateListPicker:
		return a.updateListPicker(msg)
	}

	return a, nil
}

func (a *App) updateLogin(msg tea.Msg) (tea.Model, tea.Cmd) {
	result, cmd := a.loginView.Update(msg)
	a.loginView = result

	if a.loginView.submitted {
		a.state = stateLoading
		email, password := a.loginView.email, a.loginView.password
		return a, tea.Batch(a.spinner.Tick, func() tea.Msg {
			client, stored, err := bring.LoginAndStore(email, password)
			if err != nil {
				return authErrorMsg{err: err}
			}
			return authSuccessMsg{client: client, stored: stored}
		})
	}

	if a.loginView.cancelled {
		return a, tea.Quit
	}

	return a, cmd
}

func (a *App) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case keyMsg.String() == "q":
			return a, tea.Quit

		case keyMsg.String() == "a":
			a.state = stateAddItem
			a.addItem = newAddItem()
			return a, a.addItem.Init()

		case keyMsg.String() == "e":
			if item := a.listView.selectedItem(); item != nil {
				a.state = stateEditItem
				a.editItem = newEditItem(item.ItemID, item.Spec)
				return a, a.editItem.Init()
			}

		case keyMsg.String() == "x":
			if item := a.listView.selectedItem(); item != nil {
				itemName := item.ItemID
				a.listView.removeItem(itemName)
				a.status = fmt.Sprintf("Removed: %s", itemName)
				a.statusErr = false
				return a, func() tea.Msg {
					if err := a.client.RemoveItem(a.stored.DefaultListUUID, itemName); err != nil {
						return apiErr(err, fmt.Sprintf("Error: %v", err))
					}
					return nil
				}
			}

		case keyMsg.String() == "enter":
			if item := a.listView.selectedItem(); item != nil {
				itemName, itemSpec := item.ItemID, item.Spec
				a.listView.completeItem(itemName)
				a.status = fmt.Sprintf("✓ %s", itemName)
				a.statusErr = false
				return a, func() tea.Msg {
					if err := a.client.CompleteItem(a.stored.DefaultListUUID, itemName, itemSpec); err != nil {
						return apiErr(err, fmt.Sprintf("Error: %v", err))
					}
					return nil
				}
			} else if item := a.listView.selectedRecentlyItem(); item != nil {
				itemName, spec := item.ItemID, item.Spec
				a.listView.readdItem(itemName, spec)
				a.status = fmt.Sprintf("Re-added: %s", itemName)
				a.statusErr = false
				return a, func() tea.Msg {
					if err := a.client.AddItem(a.stored.DefaultListUUID, itemName, spec); err != nil {
						return apiErr(err, fmt.Sprintf("Error re-adding %s: %v", itemName, err))
					}
					return nil
				}
			}

		case keyMsg.String() == "l":
			a.state = stateLoading
			return a, tea.Batch(a.spinner.Tick, a.loadLists())

		case keyMsg.String() == "r":
			a.state = stateLoading
			a.status = ""
			return a, tea.Batch(a.spinner.Tick, a.loadItems())
		}
	}

	a.listView.Update(msg)
	return a, nil
}

func (a *App) updateAddItem(msg tea.Msg) (tea.Model, tea.Cmd) {
	result, cmd := a.addItem.Update(msg)
	a.addItem = result

	if a.addItem.submitted {
		item, spec := a.addItem.itemName, a.addItem.spec
		a.state = stateList
		a.listView.addItem(item, spec)
		a.status = fmt.Sprintf("Added: %s", item)
		a.statusErr = false
		return a, func() tea.Msg {
			if err := a.client.AddItem(a.stored.DefaultListUUID, item, spec); err != nil {
				return apiErr(err, fmt.Sprintf("Error adding %s: %v", item, err))
			}
			return nil
		}
	}

	if a.addItem.cancelled {
		a.state = stateList
		return a, nil
	}

	return a, cmd
}

func (a *App) updateEditItem(msg tea.Msg) (tea.Model, tea.Cmd) {
	result, cmd := a.editItem.Update(msg)
	a.editItem = result

	if a.editItem.submitted {
		oldName := a.editItem.originalName
		newName, newSpec := a.editItem.itemName, a.editItem.spec
		a.state = stateList
		a.listView.updateItem(oldName, newName, newSpec)
		a.status = fmt.Sprintf("Updated: %s", newName)
		a.statusErr = false
		return a, func() tea.Msg {
			if err := a.client.EditItem(a.stored.DefaultListUUID, oldName, newName, newSpec); err != nil {
				return apiErr(err, fmt.Sprintf("Error editing %s: %v", newName, err))
			}
			return nil
		}
	}

	if a.editItem.cancelled {
		a.state = stateList
		return a, nil
	}

	return a, cmd
}

func (a *App) updateListPicker(msg tea.Msg) (tea.Model, tea.Cmd) {
	result, cmd := a.listPicker.Update(msg)
	a.listPicker = result

	if a.listPicker.selected != nil {
		a.stored.DefaultListUUID = a.listPicker.selected.ListUUID
		a.stored.DefaultListName = a.listPicker.selected.Name
		go func() {
			_ = config.Save(a.stored)
		}()
		a.state = stateLoading
		a.status = fmt.Sprintf("Switched to: %s", a.listPicker.selected.Name)
		a.statusErr = false
		return a, tea.Batch(a.spinner.Tick, a.loadItems())
	}

	if a.listPicker.cancelled {
		a.state = stateList
		return a, nil
	}

	return a, cmd
}

func (a *App) View() string {
	switch a.state {
	case stateLoading:
		return "\n " + a.spinner.View() + " Loading...\n"

	case stateLogin:
		return a.loginView.View()

	case stateList:
		view := a.listView.View()
		if a.status != "" {
			if a.statusErr {
				view += "\n" + errorStyle.Render(a.status)
			} else {
				view += "\n" + successStyle.Render(a.status)
			}
		}
		view += "\n" + statusBarStyle.Render(helpBar())
		return view

	case stateAddItem:
		return a.addItem.View()

	case stateEditItem:
		return a.editItem.View()

	case stateListPicker:
		return a.listPicker.View()
	}

	return ""
}
