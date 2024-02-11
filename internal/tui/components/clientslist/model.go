package clientslist

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/submarr/internal/core"
	"github.com/jon4hz/submarr/internal/tui/components/statusbar"
	"github.com/jon4hz/submarr/internal/tui/styles"
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
		clientList: list.New(nil, clientDelegate{}, 0, 0),
	}

	// list options
	m.clientList.SetShowStatusBar(false)
	m.clientList.SetFilteringEnabled(false)
	m.clientList.Title = "Available Clients"
	m.clientList.Styles.Title = m.clientList.Styles.Title.Copy().
		Background(styles.PurpleColor)
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
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Reload):
			// reload the clients list and start the spinner
			cmds = append(cmds,
				m.client.FetchClients(),
				m.clientList.StartSpinner(),
			)
		}

	case core.FetchClientsMsg:
		m.clientList.StopSpinner()
		cmds = append(cmds, m.clientList.SetItems(msg.Items))
		if len(msg.Errors) > 0 {
			cmds = append(cmds, statusbar.NewErrCmds(msg.Errors...)...)
		}
		return m, tea.Batch(cmds...)

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			m.clientList.CursorUp()
			return m, nil

		case tea.MouseWheelDown:
			m.clientList.CursorDown()
			return m, nil

		case tea.MouseLeft:
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
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.clientList.SetSize(width, height)
}

func (m Model) Help() [][]key.Binding {
	return DefaultKeyMap.FullHelp()
}

func (m Model) SelectedItem() list.Item {
	return m.clientList.SelectedItem()
}

func (m Model) VisibleItems() []list.Item {
	return m.clientList.VisibleItems()
}

func (m Model) Index() int {
	return m.clientList.Index()
}

func (m Model) View() string {
	return m.clientList.View()
}
