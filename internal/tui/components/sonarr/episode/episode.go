package episode

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/jon4hz/submarr/internal/core/sonarr"
	"github.com/jon4hz/submarr/internal/tui/common"
	"github.com/jon4hz/submarr/internal/tui/components/sonarr/episode/mediainfo"
	"github.com/jon4hz/submarr/internal/tui/components/statusbar"
	"github.com/jon4hz/submarr/internal/tui/overlay"
	"github.com/jon4hz/submarr/internal/tui/styles"
	sonarrAPI "github.com/jon4hz/submarr/pkg/sonarr"
)

type state int

const (
	stateEpisode state = iota + 1
	stateDetails
	stateConfirmDelete
)

type Model struct {
	common.EmbedableModel

	client        *sonarr.Client
	state         state
	episode       *sonarrAPI.EpisodeResource
	table         table.Model
	mediaInfo     common.SubModel
	confirmDelete common.SubModel
}

func New(client *sonarr.Client, episode *sonarrAPI.EpisodeResource, width, height int) common.SubModel {
	m := Model{
		state:   stateEpisode,
		client:  client,
		episode: episode,
	}

	if episode.HasFile {
		rows := []table.Row{
			{
				episode.EpisodeFile.Path,
				humanize.IBytes(uint64(episode.EpisodeFile.Size)),
				client.GetLanguageProfileByID(episode.Series.LanguageProfileID).Name,
				episode.EpisodeFile.Quality.Quality.Name,
			},
		}

		columns := []table.Column{
			{Title: "Path"},
			{Title: "Size"},
			{Title: "Language"},
			{Title: "Quality"},
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(2),
			table.WithWidth(width),
		)

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(styles.SubtleColor).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Bold(false)
		t.SetStyles(s)

		m.table = t
	}

	m.SetSize(width, height)

	return &m
}

func (m Model) Init() tea.Cmd {
	return statusbar.NewHelpCmd(DefaultKeyMap.FullHelp())
}

func (m *Model) Update(msg tea.Msg) (common.SubModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case stateEpisode:
			switch {
			case key.Matches(msg, DefaultKeyMap.Back):
				m.IsBack = true
				return m, nil
			case key.Matches(msg, DefaultKeyMap.Quit):
				m.IsQuit = true
				return m, nil
			case key.Matches(msg, DefaultKeyMap.Select):
				m.state = stateDetails
				m.mediaInfo = mediainfo.New(m.episode, m.Width, m.Height)
				return m, m.mediaInfo.Init()
			case key.Matches(msg, DefaultKeyMap.Delete):
				m.state = stateConfirmDelete
				return m, nil
			}

		case stateDetails:
			switch {
			case key.Matches(msg, DefaultKeyMap.Back):
				m.state = stateEpisode
				return m, statusbar.NewHelpCmd(DefaultKeyMap.FullHelp())
			case key.Matches(msg, DefaultKeyMap.Quit):
				m.IsQuit = true
				return m, nil
			}

		case stateConfirmDelete:
			switch {
			case key.Matches(msg, DefaultKeyMap.Back):
				m.state = stateEpisode
				return m, statusbar.NewHelpCmd(DefaultKeyMap.FullHelp())
			case key.Matches(msg, DefaultKeyMap.Quit):
				m.IsQuit = true
				return m, nil
			}
		}
	}

	switch m.state {
	case stateEpisode:
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	case stateDetails:
		if m.mediaInfo == nil {
			break
		}
		var cmd tea.Cmd
		m.mediaInfo, cmd = m.mediaInfo.Update(msg)
		return m, cmd
	case stateConfirmDelete:
		if m.confirmDelete == nil {
			break
		}
		var cmd tea.Cmd
		m.confirmDelete, cmd = m.confirmDelete.Update(msg)
		return m, cmd
	}

	return m, nil
}

var (
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			Padding(0, 1, 1, 1)

	titleStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Background(styles.PurpleColor)

	networkStyle = lipgloss.NewStyle().
			Background(styles.SonarrBlue).
			Padding(0, 1)

	qualityStyle = lipgloss.NewStyle().
			Background(styles.BlueColor).
			Padding(0, 1)
)

func (m Model) View() string {
	switch m.state {
	case stateEpisode:
		return m.episodeView()
	case stateDetails:
		return m.episodeDetailsView()
	case stateConfirmDelete:
		return m.episodeConfirmDeleteView()
	}
	return ":("
}

func (m Model) episodeView() string {
	airDate := "TBA"
	if !m.episode.AirDateUTC.IsZero() {
		airDate = m.episode.AirDateUTC.Local().Format("02 Jan 2006 at 15:04")
	}

	var s strings.Builder
	s.WriteString(titleStyle.Render(fmt.Sprintf("%s ❯ Season %d ❯ %d. %s", m.episode.Series.Title, m.episode.SeasonNumber, m.episode.EpisodeNumber, m.episode.Title)))
	s.WriteString("\n\n")

	s.WriteString("Airs:            ")
	s.WriteString(
		fmt.Sprintf("%s on %s\n",
			airDate,
			networkStyle.Render(m.episode.Series.Network),
		),
	)
	s.WriteString("Quality Profile: ")
	s.WriteString(
		qualityStyle.Render(m.client.GetQualityProfileByID(m.episode.Series.QualityProfileID).Name))
	s.WriteString("\n\n")

	overview := "No episode overview"
	if m.episode.Overview != "" {
		overview = m.episode.Overview
	}
	s.WriteString(overview)
	s.WriteString("\n\n")

	if m.episode.HasFile {
		s.WriteString(m.table.View())
		s.WriteString("")
	}

	return boxStyle.Width(m.Width).Height(m.Height).Render(s.String())
}

var overlayStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(1, 2, 1, 2)

func (m Model) episodeDetailsView() string {
	fg := overlayStyle.Render(m.mediaInfo.View())
	x := ((m.Width - lipgloss.Width(fg)) / 2)
	y := ((m.Height - lipgloss.Height(fg)) / 2)
	// make sure background fills the whole screen
	bg := m.episodeView()
	return overlay.PlaceOverlay(x, y, fg, bg)
}

func (m Model) episodeConfirmDeleteView() string {
	return "confirm delete"
}

func (m *Model) SetSize(width, height int) {
	height = height - boxStyle.GetVerticalFrameSize()

	m.Width = width
	m.Height = height

	m.table.SetWidth(width - boxStyle.GetHorizontalFrameSize())
	m.resizeTable(width - boxStyle.GetHorizontalFrameSize())

	if m.confirmDelete != nil {
		m.confirmDelete.SetSize(width, height)
	}
}

func (m *Model) resizeTable(width int) {
	if !m.episode.HasFile {
		return
	}

	rows := m.table.Rows()

	qualityWidth := max(lipgloss.Width(rows[0][3]), 7)
	languageWidth := max(lipgloss.Width(rows[0][2]), 8)
	sizeWidth := max(lipgloss.Width(rows[0][1]), 5)
	pathWidth := max(width-(qualityWidth+languageWidth+sizeWidth+10), 4)

	columns := []table.Column{
		{Title: "Path", Width: pathWidth},
		{Title: "Size", Width: sizeWidth},
		{Title: "Language", Width: languageWidth},
		{Title: "Quality", Width: qualityWidth},
	}

	m.table.SetColumns(columns)
}
