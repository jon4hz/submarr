package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/subrr/internal/core"
	"github.com/jon4hz/subrr/internal/tui/clientslist"
	"github.com/jon4hz/subrr/internal/tui/common"
)

var docStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

type State int

const (
	StateUnknown State = iota
	StateLoading
	StateError
	StateReady
)

type Model struct {
	client     *core.Client
	clientList list.Model
	spinner    spinner.Model
	state      State
	err        error
}

func New(client *core.Client) *Model {
	m := &Model{
		state:      StateLoading,
		client:     client,
		spinner:    spinner.New(spinner.WithSpinner(spinner.Points)),
		clientList: list.New(nil, list.NewDefaultDelegate(), 0, 0),
	}
	return m
}

func (m Model) Run() error {
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	return err
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		clientslist.FetchClientsListItems(m.client),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)

	case spinner.TickMsg:
		if m.state == StateLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case clientslist.ItemsMsg:
		m.state = StateReady
		return m, m.clientList.SetItems(msg.Items)

	case common.ErrMsg:
		m.state = StateError
		m.err = msg.Err
	}

	switch m.state {
	case StateReady:
		var cmd tea.Cmd
		m.clientList, cmd = m.clientList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *Model) setSize(width, height int) {
	x, y := docStyle.GetFrameSize()
	width = width - x
	height = height - y

	docStyle.Width(width)
	docStyle.Height(height)

	m.clientList.SetSize(width, height)
}

func (m Model) View() string {
	switch m.state {
	case StateLoading:
		return docStyle.Render(m.spinner.View() + "  Loading...")
	case StateError:
		return docStyle.Render(fmt.Sprintf("Error: %s", m.err))
	case StateReady:
		return docStyle.Render(m.clientList.View())
	}

	return docStyle.Render("Unknown state")
}
