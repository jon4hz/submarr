package season

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/submarr/internal/core/sonarr"
	"github.com/jon4hz/submarr/internal/tui/common"
	"github.com/jon4hz/submarr/internal/tui/components/sonarr/episode"
	sonarr_list "github.com/jon4hz/submarr/internal/tui/components/sonarr/list"
	"github.com/jon4hz/submarr/internal/tui/components/statusbar"
	sonarrAPI "github.com/jon4hz/submarr/pkg/sonarr"
	zone "github.com/lrstanley/bubblezone"
)

type state int

const (
	stateFetchEpisodes state = iota + 1
	stateShowEpisodes
	stateEpisodeDetails
)

type Model struct {
	common.EmbedableModel

	client       *sonarr.Client
	state        state
	episodesList list.Model
	spinner      common.Spinner
	episode      common.SubModel

	// make sure we only reload once at a time
	reloading bool
	mu        *sync.Mutex
}

func New(sonarr *sonarr.Client, width, height int) *Model {
	m := Model{
		client:       sonarr,
		state:        stateFetchEpisodes,
		spinner:      common.NewSpinner(),
		episodesList: sonarr_list.New(fmt.Sprintf("%s ❯ Season %d", sonarr.GetSerie().Title, sonarr.GetSeason().SeasonNumber), nil, Delegate{}, width, height),
		mu:           &sync.Mutex{},
	}

	m.SetSize(width, height)

	return &m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		statusbar.NewHelpCmd(DefaultKeyMap.FullHelp()),
		m.spinner.Tick,
		m.client.FetchSeasonEpisodes(m.client.GetSeason().SeasonNumber),
	)
}

func (m *Model) Update(msg tea.Msg) (common.SubModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case stateFetchEpisodes, stateShowEpisodes:
			switch {
			case key.Matches(msg, DefaultKeyMap.Back):
				if !m.episodesList.SettingFilter() && !m.episodesList.IsFiltered() {
					m.IsBack = true
					return m, nil
				}

			case key.Matches(msg, DefaultKeyMap.Quit):
				if !m.episodesList.SettingFilter() {
					m.IsQuit = true
					return m, nil
				}

			case key.Matches(msg, DefaultKeyMap.Reload):
				if !m.episodesList.SettingFilter() && !m.GetReloading() {
					m.SetReloading(true)
					return m, tea.Batch(
						m.episodesList.StartSpinner(),
						m.client.FetchSeasonEpisodes(m.client.GetSeason().SeasonNumber),
						statusbar.NewMessageCmd("Reloading episodes...", statusbar.WithMessageTimeout(2)),
					)
				}

			case key.Matches(msg, DefaultKeyMap.AutomaticSearch):
				if !m.episodesList.SettingFilter() {
					item, _ := m.episodesList.SelectedItem().(EpisodeItem)
					if item.episode == nil {
						return m, nil
					}
					return m, tea.Batch(
						m.client.AutomaticSearchEpisode(item.episode.ID),
						statusbar.NewMessageCmd(fmt.Sprintf("Searching for episode %d...", item.episode.EpisodeNumber), statusbar.WithMessageTimeout(2)),
					)
				}

			case key.Matches(msg, DefaultKeyMap.Select):
				if !m.episodesList.SettingFilter() {
					item, _ := m.episodesList.SelectedItem().(EpisodeItem)
					if item.episode == nil {
						return m, nil
					}
					return m, tea.Batch(
						m.selectEpisode(item.episode),
						statusbar.NewMessageCmd(fmt.Sprintf("Loading episode %d...", item.episode.EpisodeNumber), statusbar.WithMessageTimeout(2)),
					)
				}
			}
		}

	case tea.MouseMsg:
		switch m.state {
		case stateShowEpisodes:
			switch msg.Button {
			case tea.MouseButtonWheelUp:
				m.episodesList.CursorUp()
				return m, nil

			case tea.MouseButtonWheelDown:
				m.episodesList.CursorDown()
				return m, nil

			case tea.MouseButtonLeft:
				for i, listItem := range m.episodesList.VisibleItems() {
					item, _ := listItem.(EpisodeItem)
					if zone.Get(fmt.Sprintf("%d. %s", item.episode.EpisodeNumber, item.episode.Title)).InBounds(msg) {
						if i == m.episodesList.Index() {
							return m, tea.Batch(
								m.selectEpisode(item.episode),
								statusbar.NewMessageCmd(fmt.Sprintf("Loading episode %d...", item.episode.EpisodeNumber), statusbar.WithMessageTimeout(2)),
							)
						}
						m.episodesList.Select(i)
						break
					}
				}
			}
		}

	case spinner.TickMsg:
		if m.state == stateFetchEpisodes {
			var cmd tea.Cmd
			m.spinner.Model, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case sonarr.FetchSeasonEpisodesResult:
		m.episodesList.StopSpinner()
		m.SetReloading(false)
		if msg.Error != nil {
			return m, statusbar.NewErrCmd("Error while fetching episodes!")
		}
		m.state = stateShowEpisodes
		return m, m.episodesList.SetItems(episodeToItems(msg.Episodes, m.client.GetSeriesQueue()))

	case sonarr.EpisodeHistoryResult:
		if msg.Error != nil {
			return m, statusbar.NewErrCmd("Error while fetching episode history!")
		}
		m.state = stateEpisodeDetails
		m.episode = episode.New(m.client, msg.Episode, m.Width, m.Height)
		return m, m.episode.Init()
	}

	switch m.state {
	case stateShowEpisodes:
		var cmd tea.Cmd
		m.episodesList, cmd = m.episodesList.Update(msg)
		return m, cmd

	case stateEpisodeDetails:
		var cmd tea.Cmd
		m.episode, cmd = m.episode.Update(msg)

		if m.episode.Back() {
			m.state = stateShowEpisodes
			return m, statusbar.NewHelpCmd(DefaultKeyMap.FullHelp())
		}
		if m.episode.Quit() {
			m.IsQuit = true
			return m, nil
		}

		return m, cmd
	}

	return m, nil
}

func (m *Model) GetReloading() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.reloading
}

func (m *Model) SetReloading(reloading bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.reloading = reloading
}

func (m *Model) selectEpisode(episode *sonarrAPI.EpisodeResource) tea.Cmd {
	return m.client.GetEpisodeHistory(episode)
}

func (m *Model) SetSize(width, height int) {
	width -= boxStyle.GetHorizontalFrameSize()
	height -= boxStyle.GetVerticalFrameSize()

	m.Width = width
	m.Height = height

	m.episodesList.SetSize(width, height)

	if m.episode != nil {
		m.episode.SetSize(width, height+boxStyle.GetVerticalFrameSize())
	}
}

var boxStyle = lipgloss.NewStyle().
	Padding(1, 0, 0, 0)

func (m *Model) View() string {
	switch m.state {
	case stateFetchEpisodes:
		return boxStyle.Render(m.spinner.View())

	case stateShowEpisodes:
		return boxStyle.Render(m.episodesList.View())

	case stateEpisodeDetails:
		return m.episode.View()
	}

	return ""
}
