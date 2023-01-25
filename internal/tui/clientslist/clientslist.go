package clientslist

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/subrr/internal/core"
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
	m.clientList.SetShowStatusBar(false)
	m.clientList.Title = "Available Clients"
	return m
}

func (m Model) Init() tea.Cmd {
	return m.client.FetchClients()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case core.FetchClientsSuccessMsg:
		return m, m.clientList.SetItems(msg.Items)

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

func (m Model) View() string {
	return m.clientList.View()
}
