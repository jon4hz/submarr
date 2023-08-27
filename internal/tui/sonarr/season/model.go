package season

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/submarr/internal/core/sonarr"
	"github.com/jon4hz/submarr/internal/tui/common"
	sonarr_list "github.com/jon4hz/submarr/internal/tui/sonarr/list"
	"github.com/jon4hz/submarr/internal/tui/statusbar"
	zone "github.com/lrstanley/bubblezone"
)

type state int

const (
	stateFetchEpisodes state = iota + 1
	stateShowEpisodes
)

type Model struct {
	common.EmbedableModel

	client       *sonarr.Client
	state        state
	episodesList list.Model
	spinner      common.Spinner

	// make sure we only reload once at a time
	reloading bool
	mu        *sync.Mutex
}

func New(sonarr *sonarr.Client, width, height int) *Model {
	m := Model{
		client:       sonarr,
		state:        stateFetchEpisodes,
		spinner:      common.NewSpinner(),
		episodesList: sonarr_list.New(fmt.Sprintf("%s - Season %d", sonarr.GetSerie().Title, sonarr.GetSeason().SeasonNumber), nil, Delegate{}, width, height),
		mu:           &sync.Mutex{},
	}

	m.Width = width
	m.Height = height

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
		}

	case tea.MouseMsg:
		switch m.state {
		case stateShowEpisodes:
			switch msg.Type {
			case tea.MouseWheelUp:
				m.episodesList.CursorUp()
				return m, nil

			case tea.MouseWheelDown:
				m.episodesList.CursorDown()
				return m, nil

			case tea.MouseLeft:
				for i, listItem := range m.episodesList.VisibleItems() {
					item, _ := listItem.(EpisodeItem)
					if zone.Get(fmt.Sprintf("%d. %s", item.episode.EpisodeNumber, item.episode.Title)).InBounds(msg) {
						if i == m.episodesList.Index() {
							return m, func() tea.Msg { return nil } // TODO: open episode details
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
	}

	switch m.state {
	case stateShowEpisodes:
		var cmd tea.Cmd
		m.episodesList, cmd = m.episodesList.Update(msg)
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

func (m *Model) SetSize(width, height int) {
	m.Width = width
	m.Height = height
	m.episodesList.SetSize(width, height)
}

func (m *Model) View() string {
	switch m.state {
	case stateFetchEpisodes:
		return m.spinner.View()

	case stateShowEpisodes:
		return m.episodesList.View()
	}
	return ""
}
