package sonarr

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/subrr/internal/core/sonarr"
	"github.com/jon4hz/subrr/internal/tui/common"
	sonarr_list "github.com/jon4hz/subrr/internal/tui/sonarr/list"
	"github.com/jon4hz/subrr/internal/tui/sonarr/search"
	"github.com/jon4hz/subrr/internal/tui/sonarr/season"
	"github.com/jon4hz/subrr/internal/tui/sonarr/series"
	"github.com/jon4hz/subrr/internal/tui/statusbar"
	sonarrAPI "github.com/jon4hz/subrr/pkg/sonarr"
	zone "github.com/lrstanley/bubblezone"
)

type state int

const (
	stateUnknown state = iota
	stateLoading
	stateSeries
	stateSeriesLoading
	stateSeriesDetails
	stateSeason
	stateSearch
)

type Model struct {
	common.EmbedableModel

	client *sonarr.Client

	seriesList list.Model

	submodel common.SubModel

	spinner common.Spinner

	state state
}

func New(c *sonarr.Client, width, height int) *Model {
	m := Model{
		state:      stateLoading,
		client:     c,
		seriesList: sonarr_list.New("Series", nil, series.Delegate{}, width, height),
		spinner:    common.NewSpinner(),
	}

	m.Width = width
	m.Height = height

	m.seriesList.InfiniteScrolling = true
	m.seriesList.FilterInput.Prompt = "Search: "

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

func (m *Model) Update(msg tea.Msg) (common.SubModel, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// handle keybindings per state
		switch m.state {
		case stateLoading:
			switch {
			case key.Matches(msg, DefaultKeyMap.Back):
				m.IsBack = true

			case key.Matches(msg, DefaultKeyMap.Quit):
				m.IsQuit = true
			}

		case stateSeries:
			switch {
			case key.Matches(msg, DefaultKeyMap.Back):
				if !m.seriesList.SettingFilter() && !m.seriesList.IsFiltered() {
					m.IsBack = true
				}

			case key.Matches(msg, DefaultKeyMap.Quit):
				if !m.seriesList.SettingFilter() {
					m.IsQuit = true
				}

			case key.Matches(msg, DefaultKeyMap.Reload):
				if !m.seriesList.SettingFilter() {
					cmds = append(cmds,
						m.client.FetchSeries(),
						m.seriesList.StartSpinner(),
						statusbar.NewMessageCmd("Reloading...", statusbar.WithMessageTimeout(2)),
					)
				}

			case key.Matches(msg, DefaultKeyMap.Select):
				item, _ := m.seriesList.SelectedItem().(sonarr.SeriesItem)
				if !m.seriesList.SettingFilter() {
					cmd := m.loadSeries(item.Series)
					return m, cmd
				}

			case key.Matches(msg, DefaultKeyMap.AddNew):
				return m, m.addNewSeries()
			}

		case stateSeriesLoading:
			switch {
			case key.Matches(msg, DefaultKeyMap.Back):
				m.state = stateSeries
				return m, nil

			case key.Matches(msg, DefaultKeyMap.Quit):
				if !m.seriesList.SettingFilter() {
					m.IsQuit = true
				}
				return m, nil
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
							cmd := m.loadSeries(item.Series)
							return m, cmd
						}
						// else select the item
						m.seriesList.Select(i)
						break
					}
				}
			}
		}

	case sonarr.FetchSerieResult:
		switch m.state {
		case stateSeriesLoading:
			if msg.Error != nil {
				return m, statusbar.NewErrCmd("Failed to fetch series")
			}
			m.state = stateSeriesDetails
			m.client.SetSerie(msg.Serie)
			m.submodel = series.New(m.client, m.Width, m.Height)

			return m, m.submodel.Init()
		}

	case sonarr.FetchSeriesResult:
		switch m.state {
		case stateSeriesDetails:
			break
		default:
			m.seriesList.StopSpinner()

			m.state = stateSeries
			if msg.Error != nil {
				cmds = append(cmds, statusbar.NewErrCmd("Failed to fetch series"))
			} else {
				cmds = append(cmds, m.seriesList.SetItems(msg.Items))
			}
			cmds = append(cmds, statusbar.NewHelpCmd(DefaultKeyMap.FullHelp()))

			return m, tea.Batch(cmds...)
		}

	case search.SeriesAlreadyAddedMsg:
		switch m.state {
		case stateSearch:
			return m, m.loadSeries(msg.Series)
		}

	case sonarr.AddSeriesResult:
		switch m.state {
		case stateSearch:
			m.state = stateSeries
			if msg.Error != nil {
				cmds = append(cmds, statusbar.NewErrCmd("Failed to add series"))
			} else {
				cmds = append(cmds,
					m.seriesList.SetItems(msg.Items),
					statusbar.NewMessageCmd(fmt.Sprintf("Added Series: %s", msg.AddedTitle)),
					statusbar.NewHelpCmd(DefaultKeyMap.FullHelp()),
				)
			}
			return m, tea.Batch(cmds...)
		}

	case series.SelectSeasonMsg:
		return m, m.selectSeason(m.client.GetSerie().Seasons[msg])
	}

	switch m.state {
	case stateLoading:
		var cmd tea.Cmd
		m.spinner.Model, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case stateSeries:
		var cmd tea.Cmd
		m.seriesList, cmd = m.seriesList.Update(msg)
		cmds = append(cmds, cmd)

	case stateSeriesLoading:
		var cmd tea.Cmd
		m.spinner.Model, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	default:
		var cmd tea.Cmd
		m.submodel, cmd = m.submodel.Update(msg)
		cmds = append(cmds, cmd)

		if m.submodel.Quit() {
			return m, tea.Quit
		}

		if m.submodel.Back() {
			switch m.state {
			case stateSeriesLoading, stateSeriesDetails, stateSearch:
				m.state = stateSeries
				cmds = append(cmds,
					// reset the help of the statusbar
					statusbar.NewHelpCmd(DefaultKeyMap.FullHelp()),
				)

			case stateSeason:
				cmds = append(cmds,
					m.loadSeries(m.client.GetSerie()),
				)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) loadSeries(seriesResource *sonarrAPI.SeriesResource) tea.Cmd {
	m.state = stateSeriesLoading
	m.spinner.Message = common.GetRandomLoadingMessage()
	m.client.SetSerie(seriesResource)
	return tea.Batch(
		m.client.ReloadSerie(),
		m.spinner.Tick,
	)
}

func (m *Model) selectSeason(seasonResource *sonarrAPI.SeasonResource) tea.Cmd {
	m.state = stateSeason
	m.client.SetSeason(seasonResource)
	m.submodel = season.New(m.client, m.Width, m.Height)

	return m.submodel.Init()
}

func (m *Model) addNewSeries() tea.Cmd {
	m.state = stateSearch
	m.submodel = search.New(m.client, m.Width, m.Height)

	return m.submodel.Init()
}

func (m *Model) SetSize(width, height int) {
	m.Width = width
	m.Height = height

	m.seriesList.SetSize(width, height)

	if m.submodel != nil {
		m.submodel.SetSize(width, height)
	}
}

func (m Model) View() string {
	switch m.state {
	case stateLoading, stateSeriesLoading:
		return m.spinner.View()

	case stateSeries:
		return m.seriesList.View()

	default:
		if m.submodel != nil {
			return m.submodel.View()
		}
		return "unknown state"
	}
}
