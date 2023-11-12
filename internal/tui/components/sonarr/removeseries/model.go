package removeseries

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/jon4hz/submarr/internal/core/sonarr"
	"github.com/jon4hz/submarr/internal/tui/common"
	"github.com/jon4hz/submarr/internal/tui/components/statusbar"
	"github.com/jon4hz/submarr/internal/tui/components/toggle"
	"github.com/jon4hz/submarr/internal/tui/styles"
	sonarrAPI "github.com/jon4hz/submarr/pkg/sonarr"
)

type Model struct {
	common.EmbedableModel

	client *sonarr.Client
	series *sonarrAPI.SeriesResource

	deleteFiles  toggle.Model
	addExclusion toggle.Model

	selectedOption     rmOption
	longestOptionWidth int
}

type rmOption int

const (
	rmOptionDeleteFiles rmOption = iota + 1
	rmOptionAddExclusion
	rmOptionRemoveSeries
)

var rmOptions = map[rmOption]string{
	rmOptionDeleteFiles:  "Delete files",
	rmOptionAddExclusion: "Add to exclusion list",
	rmOptionRemoveSeries: "",
}

func New(client *sonarr.Client, series *sonarrAPI.SeriesResource, width, height int) common.SubModel {
	m := Model{
		client:             client,
		series:             series,
		selectedOption:     1,
		deleteFiles:        toggle.New(),
		addExclusion:       toggle.New(),
		longestOptionWidth: getLongestOptionWidth(),
	}

	m.Width = width
	m.Height = height

	rmOptions[rmOptionDeleteFiles] = fmt.Sprintf("Delete %d files", series.Statistics.EpisodeFileCount)

	return &m
}

func getLongestOptionWidth() int {
	var longest int
	for _, option := range rmOptions {
		if len(option) > longest {
			longest = len(option)
		}
	}
	return longest
}

func (m Model) Init() tea.Cmd {
	return statusbar.NewHelpCmd(DefaultKeyMap.FullHelp())
}

func (m *Model) Update(msg tea.Msg) (common.SubModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Down):
			m.nextOption()
		case key.Matches(msg, DefaultKeyMap.Up):
			m.previousOption()
		case key.Matches(msg, DefaultKeyMap.Select):
			if m.selectedOption == rmOptionRemoveSeries {
				return m, m.rmSeries()
			}
		case key.Matches(msg, DefaultKeyMap.Delete):
			return m, m.rmSeries()
		case key.Matches(msg, DefaultKeyMap.Back):
			m.IsBack = true
		case key.Matches(msg, DefaultKeyMap.Quit):
			m.IsQuit = true
		}

		switch m.selectedOption {
		case rmOptionDeleteFiles:
			var cmd tea.Cmd
			m.deleteFiles, cmd = m.deleteFiles.Update(msg)
			return m, cmd
		case rmOptionAddExclusion:
			var cmd tea.Cmd
			m.addExclusion, cmd = m.addExclusion.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m *Model) nextOption() {
	m.selectedOption++
	if int(m.selectedOption) > len(rmOptions) {
		m.selectedOption = 1
	}
}

func (m *Model) previousOption() {
	m.selectedOption--
	if m.selectedOption < 1 {
		m.selectedOption = rmOption(len(rmOptions))
	}
}

func (m Model) rmSeries() tea.Cmd {
	deleteFiles := m.deleteFiles.Toggled()
	addExclusion := m.addExclusion.Toggled()
	return m.client.DeleteSeries(m.series, deleteFiles, addExclusion)
}

var (
	boxStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).Padding(1, 2, 1, 2)
	titleStyle  = lipgloss.NewStyle().Align(lipgloss.Center).Bold(true).Underline(true)
	keyStyle    = lipgloss.NewStyle().Align(lipgloss.Right).Margin(1, 2, 1, 0)
	valueStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).Padding(0, 1, 0)
	buttonStyle = lipgloss.NewStyle().Align(lipgloss.Center).Border(lipgloss.RoundedBorder(), true).Padding(0, 1, 0)
)

func (m Model) View() string {
	var s strings.Builder

	kvs := [][]string{
		{
			rmOptions[rmOptionDeleteFiles],
			m.deleteFiles.View(),
		},
		{
			rmOptions[rmOptionAddExclusion],
			m.addExclusion.View(),
		},
	}

	lines := make([]string, len(kvs))
	for i, kv := range kvs {
		var color lipgloss.TerminalColor = styles.SubtleColor
		if i == int(m.selectedOption)-1 {
			color = styles.SonarrBlue
		}
		lines[i] = lipgloss.JoinHorizontal(lipgloss.Left,
			keyStyle.Width(m.longestOptionWidth).Render(kv[0]),
			valueStyle.Width(m.longestOptionWidth).BorderForeground(color).Render(kv[1]),
		)
	}

	options := lipgloss.JoinVertical(lipgloss.Right,
		lines...,
	)

	width := lipgloss.Width(options)
	s.WriteString(
		titleStyle.Width(width).Render(fmt.Sprintf("%s (%d)", m.series.Title, m.series.Year)),
	)
	s.WriteString("\n\n")

	s.WriteString(options)
	s.WriteString("\n\n")

	if m.deleteFiles.Toggled() {
		errStyle := lipgloss.NewStyle().Foreground(styles.ErrorColor)
		s.WriteString(
			errStyle.Width(width).Render(fmt.Sprintf("The series folder %q and all of its content will be deleted.", m.series.Path)),
		)
		s.WriteByte('\n')
		s.WriteString(
			errStyle.Width(width).Render(fmt.Sprintf("%d episode files totaling %s", m.series.Statistics.EpisodeFileCount, humanize.IBytes(uint64(m.series.Statistics.SizeOnDisk)))),
		)
		s.WriteString("\n\n")
	}

	var color lipgloss.TerminalColor = styles.SubtleColor
	if m.selectedOption == rmOptionRemoveSeries {
		color = styles.SonarrBlue
	}
	s.WriteString(
		lipgloss.Place(width, 1, lipgloss.Center,
			lipgloss.Top, buttonStyle.BorderForeground(color).Render("Delete Series")),
	)

	return boxStyle.MaxWidth(m.Width).Render(s.String())
}

func (m *Model) SetSize(width, height int) {
	width -= boxStyle.GetHorizontalFrameSize()
	height -= boxStyle.GetVerticalFrameSize()

	m.Width = width
	m.Height = height
}
