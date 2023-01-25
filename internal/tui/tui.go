package tui

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/subrr/internal/core"
	"github.com/jon4hz/subrr/internal/tui/clientslist"
	"github.com/jon4hz/subrr/internal/tui/common"
	zone "github.com/lrstanley/bubblezone"
)

var docStyle = lipgloss.NewStyle().Margin(1)

type State int

const (
	StateUnknown State = iota
	StateLoading
	StateError
	StateReady
)

type Model struct {
	client      *core.Client
	clientslist clientslist.Model
	spinner     spinner.Model
	state       State
	err         error
}

func New(client *core.Client) *Model {
	m := &Model{
		state:       StateLoading,
		client:      client,
		spinner:     spinner.New(spinner.WithSpinner(spinner.Points)),
		clientslist: clientslist.New(client),
	}
	return m
}

func (m Model) Run() error {
	_, err := tea.NewProgram(m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	).Run()
	return err
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.clientslist.Init(),
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

	case core.FetchClientsSuccessMsg:
		m.state = StateReady

	case core.FetchClientsErrorMsg:
		m.state = StateError
		m.err = errors.New(msg.Description)

	case common.ErrMsg:
		m.state = StateError
		m.err = msg.Err
	}

	switch m.state {
	case StateReady:
		var cmd tea.Cmd
		m.clientslist, cmd = m.clientslist.Update(msg)
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

	m.clientslist.SetSize(width-x, height)
}

func (m Model) View() string {
	switch m.state {
	case StateLoading:
		return docStyle.Render(m.spinner.View() + "  Loading...")
	case StateError:
		return docStyle.Render(fmt.Sprintf("Error: %s", m.err))
	case StateReady:
		return zone.Scan(docStyle.Render(m.clientslist.View()))
	}

	return docStyle.Render("Unknown state")
}
