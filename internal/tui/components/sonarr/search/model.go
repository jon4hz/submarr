package search

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/submarr/internal/core/sonarr"
	"github.com/jon4hz/submarr/internal/tui/common"
	"github.com/jon4hz/submarr/internal/tui/components/sonarr/addseries"
	sonarr_list "github.com/jon4hz/submarr/internal/tui/components/sonarr/list"
	"github.com/jon4hz/submarr/internal/tui/components/statusbar"
	"github.com/jon4hz/submarr/internal/tui/overlay"
	sonarrAPI "github.com/jon4hz/submarr/pkg/sonarr"
)

type SeriesAlreadyAddedMsg struct {
	Series *sonarrAPI.SeriesResource
}

type state int

const (
	stateInput state = iota + 1
	stateSearching
	stateShowResults
	stateAddSeries
)

type Model struct {
	common.EmbedableModel

	client  *sonarr.Client
	state   state
	spinner spinner.Model
	input   textinput.Model
	result  list.Model
	add     common.SubModel
	cancel  context.CancelFunc
}

func New(sonarr *sonarr.Client, width, height int) *Model {
	m := Model{
		client:  sonarr,
		state:   stateInput,
		spinner: spinner.New(spinner.WithSpinner(spinner.Points)),
		input:   textinput.New(),
		result:  sonarr_list.New("Search Results", nil, Delegate{}, width, height),
	}

	m.SetSize(width, height)

	m.input.Placeholder = "eg. Breaking Bad, tvdb:####"
	m.input.Width = width

	return &m
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		statusbar.NewHelpCmd(InputKeyMap.FullHelp()),
		m.input.Focus(),
	)
}

func (m *Model) Update(msg tea.Msg) (common.SubModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, InputKeyMap.Back):
			switch m.state {
			case stateInput:
				m.IsBack = true
				return m, nil
			case stateShowResults:
				if m.result.IsFiltered() || m.result.SettingFilter() {
					break
				}
				m.state = stateInput
				return m, tea.Sequence(
					statusbar.NewHelpCmd(InputKeyMap.FullHelp()),
					m.input.Focus(),
				)
			case stateSearching:
				m.state = stateInput
				return m, tea.Sequence(
					statusbar.NewHelpCmd(InputKeyMap.FullHelp()),
					m.input.Focus(),
				)
			}

		case key.Matches(msg, InputKeyMap.Quit):
			if m.state == stateAddSeries {
				break
			}
			m.IsQuit = true
			return m, nil
		}

	case sonarr.SearchSeriesResult:
		if m.state != stateSearching {
			break
		}
		if msg.Error != nil {
			m.state = stateInput
			return m, tea.Sequence(
				statusbar.NewErrCmd(msg.Error.Error()),
				statusbar.NewHelpCmd(InputKeyMap.FullHelp()),
				m.input.Focus(),
			)
		}

		m.state = stateShowResults

		return m, tea.Sequence(
			m.result.SetItems(msg.Items),
			statusbar.NewHelpCmd(ResultKeyMap.FullHelp()),
		)
	}

	switch m.state {
	case stateInput:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if key.Matches(msg, InputKeyMap.Select) {
				term := strings.TrimSpace(m.input.Value())
				return m, m.searchSeries(term)
			}
		}

		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd

	case stateSearching:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case stateShowResults:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if key.Matches(msg, InputKeyMap.Select) {
				item, ok := m.result.SelectedItem().(sonarr.SeriesItem)
				if !ok {
					break
				}
				if !item.Series.Added.IsZero() {
					return m, func() tea.Msg {
						return SeriesAlreadyAddedMsg{Series: item.Series}
					}
				}

				return m, m.addSeries(item.Series)
			}
		}
		var cmd tea.Cmd
		m.result, cmd = m.result.Update(msg)
		return m, cmd

	case stateAddSeries:
		var cmd tea.Cmd
		m.add, cmd = m.add.Update(msg)

		if m.add.Quit() {
			m.IsQuit = true
			return m, nil
		}

		if m.add.Back() {
			m.state = stateShowResults
			return m, statusbar.NewHelpCmd(ResultKeyMap.FullHelp())
		}

		return m, cmd
	}

	return m, nil
}

func (m *Model) searchSeries(term string) tea.Cmd {
	m.state = stateSearching
	m.input.Blur()
	cmd, cancel := m.client.SearchSeries(term)
	// cancel previous search
	if m.cancel != nil {
		m.cancel()
	}
	// set new cancel function
	m.cancel = cancel
	return tea.Batch(
		m.spinner.Tick,
		cmd,
	)
}

func (m *Model) addSeries(series *sonarrAPI.SeriesResource) tea.Cmd {
	m.state = stateAddSeries
	m.add = addseries.New(m.client, series, m.Width, m.Height)
	return m.add.Init()
}

func (m *Model) SetSize(width, height int) {
	width -= boxStyle.GetVerticalFrameSize()
	height -= boxStyle.GetHorizontalFrameSize()
	m.Width = width
	m.Height = height

	m.input.Width = width

	//inputHeight := lipgloss.Height(m.inputView())
	m.result.SetSize(width, height-1)

	if m.state == stateAddSeries {
		m.add.SetSize(min(width, 54), min(height, 34))
	}
}

var boxStyle = lipgloss.NewStyle().
	Padding(1, 2, 0, 2)

func (m Model) View() string {
	switch m.state {
	case stateInput:
		return boxStyle.Render(m.inputView())
	case stateSearching:
		return boxStyle.Render(m.searchView())
	case stateShowResults:
		return boxStyle.Render(m.resultView())
	case stateAddSeries:
		fg := m.add.View()
		x := ((m.Width - lipgloss.Width(fg)) / 2)
		y := ((m.Height - lipgloss.Height(fg)) / 2)
		// make sure background fills the whole screen
		bg := boxStyle.Render(m.resultView())
		return overlay.PlaceOverlay(x, y, fg, bg)
	}
	return "unknown"
}

func (m Model) inputView() string {
	var s strings.Builder
	s.WriteString("🔍 Search for new series:\n\n")
	s.WriteString(m.input.View())
	s.WriteByte('\n')
	s.WriteByte('\n')
	return s.String()
}

func (m Model) searchView() string {
	var s strings.Builder
	s.WriteString(m.inputView())
	s.WriteString(m.spinner.View())
	s.WriteString("  Searching...")
	return s.String()
}

func (m Model) resultView() string {
	var s strings.Builder
	s.WriteString(m.inputView())
	s.WriteString(m.result.View())
	return s.String()
}
