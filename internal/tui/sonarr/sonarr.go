package sonarr

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/subrr/internal/core/sonarr"
	"github.com/jon4hz/subrr/internal/tui/common"
	"github.com/jon4hz/subrr/internal/tui/sonarr/series"
	"github.com/jon4hz/subrr/internal/tui/statusbar"
	zone "github.com/lrstanley/bubblezone"
)

type state int

const (
	stateUnknown state = iota
	stateLoading
	stateSeries
	stateSeriesDetails
)

type Model struct {
	common.EmbedableModel

	client *sonarr.Client

	seriesList list.Model

	spinner        spinner.Model
	loadingMessage string

	state state
}

func New(c *sonarr.Client, width, height int) *Model {
	m := Model{
		state:          stateLoading,
		client:         c,
		seriesList:     list.NewModel(nil, series.Delegate{}, width, height),
		spinner:        spinner.New(spinner.WithSpinner(spinner.Points)),
		loadingMessage: common.GetRandomLoadingMessage(),
	}
	m.Width = width
	m.Height = height

	m.seriesList.DisableQuitKeybindings()
	m.seriesList.Title = "Series"
	m.seriesList.Styles.Title = m.seriesList.Styles.Title.Copy().
		Background(lipgloss.Color("#7B61FF"))
	m.seriesList.SetShowHelp(false)

	m.seriesList.FilterInput.Prompt = "Search: "
	m.seriesList.FilterInput.CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00CCFF"))

	return &m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		statusbar.NewTitleCmd("Sonarr", statusbar.WithTitleForeground(lipgloss.Color("#00CCFF"))),
		statusbar.NewHelpCmd(DefaultKeyMap.FullHelp()),
		m.spinner.Tick,
		m.client.FetchSeries(),
	)
}

func (m *Model) Update(msg tea.Msg) (common.ClientModel, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// general keybindings for all states
		switch msg.String() {
		case "esc":
			if !m.seriesList.SettingFilter() && !m.seriesList.IsFiltered() {
				m.IsBack = true
			}

		case "q":
			if !m.seriesList.SettingFilter() {
				m.IsQuit = true
			}
		}

		// keybindings for specific states
		switch m.state {
		case stateSeries:
			switch {
			case key.Matches(msg, DefaultKeyMap.Reload):
				cmds = append(cmds,
					m.client.FetchSeries(),
					m.seriesList.StartSpinner(),
					statusbar.NewMessageCmd("Reloading...", statusbar.WithMessageTimeout(2)),
				)

			case key.Matches(msg, DefaultKeyMap.Select):
				if !m.seriesList.SettingFilter() {
				}
			}
		}

	case tea.MouseMsg:
		switch m.state {
		case stateSeries:
			switch msg.Type {
			case tea.MouseWheelUp:
				m.seriesList.CursorUp()
				return m, nil

			case tea.MouseWheelDown:
				m.seriesList.CursorDown()
				return m, nil

			case tea.MouseLeft:
				for i, listItem := range m.seriesList.VisibleItems() {
					item, _ := listItem.(sonarr.SeriesItem)
					if zone.Get(item.Series.Title).InBounds(msg) {
						// if we click on an already selected item, open the details
						if i == m.seriesList.Index() {
							// TODO: open series details

						}
						// else select the item
						m.seriesList.Select(i)
						break
					}
				}
			}
		}

	case sonarr.FetchSeriesResult:
		m.seriesList.StopSpinner()

		m.state = stateSeries
		if msg.Error != nil {
			cmds = append(cmds, statusbar.NewErrCmd("Failed to fetch series"))
		}
		cmds = append(cmds, m.seriesList.SetItems(msg.Items))

		return m, tea.Batch(cmds...)
	}

	switch m.state {
	case stateLoading:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case stateSeries:
		var cmd tea.Cmd
		m.seriesList, cmd = m.seriesList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) SetSize(width, height int) {
	m.Width = width
	m.Height = height

	m.seriesList.SetSize(width, height)
}

func (m Model) View() string {
	switch m.state {
	case stateLoading:
		return m.spinner.View() + "  " + m.loadingMessage
	case stateSeries:
		return m.seriesList.View()
	}

	return "unknown state"
}
