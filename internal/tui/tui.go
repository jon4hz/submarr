package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/subrr/internal/core"
	"github.com/jon4hz/subrr/internal/logging"
	"github.com/jon4hz/subrr/internal/tui/clientslist"
	"github.com/jon4hz/subrr/internal/tui/common"
	"github.com/jon4hz/subrr/internal/tui/sonarr"
	"github.com/jon4hz/subrr/internal/tui/statusbar"
	zone "github.com/lrstanley/bubblezone"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1)
	errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
)

type state int

const (
	stateUnknown state = iota
	stateLoading
	stateReady
	stateClient
)

type Model struct {
	// totalWidth and totalHeight are the width and height of the entire terminal.
	totalWidth  int
	totalHeight int

	// availableWidth and availableHeight are the width and height of the available
	// space in the terminal after the statusbar has been drawn.
	availableWidth  int
	availableHeight int

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
	state state

	// clientModel is the model of the active client.
	clientModel common.ClientModel
}

func New(client *core.Client) *Model {
	m := &Model{
		state:          stateLoading,
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
		// handle some special keys here.
		// these keys are not handled per state.
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

		// handle keys per state
		switch m.state {
		case stateReady:
			switch {
			case key.Matches(msg, clientslist.DefaultKeyMap.Select):
				selected := m.clientslist.SelectedItem()
				if selected != nil {
					item, ok := selected.(clientslist.ClientsItem)
					if !ok {
						logging.Log.Error().Msg("could not convert selected item to clientsitem")
						return m, nil
					}
					cmd := m.enterClient(item)
					return m, cmd
				}
			}
		}

	case tea.MouseMsg:
		// handle mouse events for all states
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

		// handle mouse events per state
		switch m.state {
		case stateReady:
			switch msg.Type {
			case tea.MouseLeft:
				for i, listItem := range m.clientslist.VisibleItems() {
					item, _ := listItem.(clientslist.ClientsItem)
					// Check each item to see if it's in bounds.
					if zone.Get(item.String()).InBounds(msg) {
						// if so, check if it's the selected item.
						if i == m.clientslist.Index() {
							// if so, enter the client.
							cmd := m.enterClient(item)
							return m, cmd
						}
						break
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)

	case spinner.TickMsg:
		if m.state == stateLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case core.FetchClientsMsg:
		if m.state != stateReady {
			m.state = stateReady
			cmds = append(cmds,
				statusbar.NewMessageCmd("Welcome to Subrr!", statusbar.WithMessageTimeout(2)),
				statusbar.NewHelpCmd(m.clientslist.Help()),
			)
		}

	// also hijack the statusbar.SetHelpMsg because that can have an impact on the layout
	case statusbar.SetHelpMsg:
		var cmd tea.Cmd
		m.statusbar, cmd = m.statusbar.Update(msg)
		cmds = append(cmds, cmd)
		m.setSize(m.totalWidth, m.totalHeight)
		return m, tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	m.statusbar, cmd = m.statusbar.Update(msg)
	cmds = append(cmds, cmd)

	switch m.state {
	case stateReady:
		var cmd tea.Cmd
		m.clientslist, cmd = m.clientslist.Update(msg)
		cmds = append(cmds, cmd)

	case stateClient:
		var cmd tea.Cmd
		m.clientModel, cmd = m.clientModel.Update(msg)
		cmds = append(cmds, cmd)

		if m.clientModel.Quit() {
			return m, tea.Quit
		}

		if m.clientModel.Back() {
			m.state = stateReady
			cmds = append(cmds,
				// reset the title of the statusbar
				statusbar.NewTitleCmd("Subrr", statusbar.WithTitleForeground(lipgloss.Color("#39FF14"))),
				// reset the help of the statusbar
				statusbar.NewHelpCmd(m.clientslist.Help()),
			)
		}
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

	m.availableWidth = width
	m.availableHeight = height

	m.clientslist.SetSize(width, height)

	if m.state == stateClient {
		m.clientModel.SetSize(width, height)
	}
}

func (m *Model) enterClient(item clientslist.ClientsItem) tea.Cmd {
	if !item.Available() {
		return nil
	}

	switch strings.ToLower(item.String()) {
	case "sonarr":
		m.state = stateClient
		m.clientModel = sonarr.New(m.client.Sonarr, m.availableWidth, m.availableHeight)
		return m.clientModel.Init()
	}

	return nil
}

func (m Model) View() string {
	switch m.state {
	case stateLoading:
		return docStyle.Render(m.spinner.View() + "  " + m.loadingMessage)

	case stateReady:
		return zone.Scan(
			lipgloss.JoinVertical(lipgloss.Top,
				docStyle.Render(m.clientslist.View()),
				m.statusbar.View(),
			),
		)

	case stateClient:
		return zone.Scan(
			lipgloss.JoinVertical(lipgloss.Top,
				docStyle.Render(m.clientModel.View()),
				m.statusbar.View(),
			),
		)
	}

	// this should never happen
	return docStyle.Render("Unknown state")
}
