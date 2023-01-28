package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/subrr/internal/core"
	"github.com/jon4hz/subrr/internal/tui/clientslist"
	"github.com/jon4hz/subrr/internal/tui/common"
	"github.com/jon4hz/subrr/internal/tui/statusbar"
	zone "github.com/lrstanley/bubblezone"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1)
	errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
)

type State int

const (
	StateUnknown State = iota
	StateLoading
	StateError
	StateReady
)

type Model struct {
	// totalWidth and totalHeight are the width and height of the entire terminal.
	totalWidth  int
	totalHeight int

	// client is the core client used to fetch data from the API.
	client *core.Client

	// clientslist is the startview of the application and shows all clients.
	clientslist clientslist.Model

	// spinner is the loading spinner.
	spinner spinner.Model

	// loadingMessage is the message shown while loading.
	loadingMessage string

	// statusbar is the statusbar of the application.
	statusbar statusbar.Model

	// state is the current state of the application.
	state State
}

func New(client *core.Client) *Model {
	m := &Model{
		state:          StateLoading,
		client:         client,
		spinner:        spinner.New(spinner.WithSpinner(spinner.Points)),
		clientslist:    clientslist.New(client),
		statusbar:      statusbar.New("Subrr"),
		loadingMessage: common.GetRandomLoadingMessage(),
	}

	// statusbar options
	m.statusbar.TitleForeground = lipgloss.Color("#39FF14")
	m.statusbar.TitleBackground = lipgloss.Color("#313041")

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
		m.statusbar.Init(),
		m.clientslist.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		// handle the statusbar gracefully here.
		// Because after toggeling the help view, the other views must be resized.
		case "?":
			var cmd tea.Cmd
			m.statusbar, cmd = m.statusbar.Update(msg)
			cmds = append(cmds, cmd)
			m.setSize(m.totalWidth, m.totalHeight)
			return m, tea.Batch(cmds...)
		}

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseLeft:
			// handle the statusbar gracefully here.
			// Because after toggeling the help view, the other views must be resized.
			if zone.Get("toggle-help").InBounds(msg) {
				var cmd tea.Cmd
				m.statusbar, cmd = m.statusbar.Update(msg)
				cmds = append(cmds, cmd)
				m.setSize(m.totalWidth, m.totalHeight)
				return m, tea.Batch(cmds...)
			}
		}

	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)

	case spinner.TickMsg:
		if m.state == StateLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case core.FetchClientsMsg:
		m.state = StateReady
		cmds = append(cmds,
			func() tea.Msg {
				return statusbar.NewMessageMsg("Welcome to Subrr!", statusbar.WithMessageTimeout(2))
			},
			func() tea.Msg {
				return statusbar.NewHelpMsg(m.clientslist.Help())
			},
		)
	}

	var cmd tea.Cmd
	m.statusbar, cmd = m.statusbar.Update(msg)
	cmds = append(cmds, cmd)

	switch m.state {
	case StateReady:
		var cmd tea.Cmd
		m.clientslist, cmd = m.clientslist.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) setSize(width, height int) {
	m.totalWidth = width
	m.totalHeight = height

	m.statusbar.SetWidth(width)

	// this will render the statusbar and return the height.
	// it is a bit suboptimal, but since setHeight is not called very often, it is ok.
	statusHeight := m.statusbar.GetHeight()

	x, y := docStyle.GetFrameSize()
	width = width - x
	height = height - y

	// deduct the statusbar height
	height = height - statusHeight

	docStyle.Width(width)
	docStyle.Height(height)
	width = width - x

	m.clientslist.SetSize(width, height)
}

func (m Model) View() string {
	switch m.state {
	case StateLoading:
		return docStyle.Render(m.spinner.View() + "  " + m.loadingMessage)

	case StateReady:
		return zone.Scan(
			lipgloss.JoinVertical(lipgloss.Top,
				docStyle.Render(m.clientslist.View()),
				m.statusbar.View(),
			),
		)
	}

	// this should never happen
	return docStyle.Render("Unknown state")
}
