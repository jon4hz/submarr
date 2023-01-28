package clientslist

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/subrr/internal/core"
	"github.com/jon4hz/subrr/internal/tui/common"
	zone "github.com/lrstanley/bubblezone"
)

type Model struct {
	client     *core.Client
	clientList list.Model
	width      int
}

func New(client *core.Client) Model {
	m := Model{
		client:     client,
		clientList: list.New(nil, newClientDelegate(), 0, 0),
	}

	// list options
	m.clientList.SetShowStatusBar(false)
	m.clientList.SetFilteringEnabled(false)
	m.clientList.Title = "Available Clients"
	m.clientList.Styles.Title = m.clientList.Styles.Title.Copy().
		Background(lipgloss.Color("#7B61FF"))
	m.clientList.SetShowHelp(false)

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.client.FetchClients(),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case core.FetchClientsMsg:
		cmds = append(cmds, m.clientList.SetItems(msg.Items))
		if len(msg.Errors) > 0 {
			cmds = append(cmds, common.NewErrCmds(msg.Errors...)...)
		}
		return m, tea.Batch(cmds...)

	case tea.MouseMsg:
		if msg.Type == tea.MouseWheelUp {
			m.clientList.CursorUp()
			return m, nil
		}

		if msg.Type == tea.MouseWheelDown {
			m.clientList.CursorDown()
			return m, nil
		}

		if msg.Type == tea.MouseLeft {
			for i, listItem := range m.clientList.VisibleItems() {
				item, _ := listItem.(ClientsItem)
				// Check each item to see if it's in bounds.
				if zone.Get(item.String()).InBounds(msg) {
					// If so, select it in the list.
					m.clientList.Select(i)
					break
				}
			}
		}
	}

	var cmd tea.Cmd
	m.clientList, cmd = m.clientList.Update(msg)
	return m, cmd
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.clientList.SetSize(width, height)
}

func (m Model) Help() [][]key.Binding {
	return m.clientList.FullHelp()
}

func (m Model) View() string {
	return m.clientList.View()
}
