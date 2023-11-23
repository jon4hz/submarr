package episode

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/jon4hz/submarr/internal/core/sonarr"
	"github.com/jon4hz/submarr/internal/tui/common"
	"github.com/jon4hz/submarr/internal/tui/components/statusbar"
	"github.com/jon4hz/submarr/internal/tui/styles"
	sonarrAPI "github.com/jon4hz/submarr/pkg/sonarr"
)

type Model struct {
	common.EmbedableModel

	client  *sonarr.Client
	episode *sonarrAPI.EpisodeResource
	table   table.Model
}

func New(client *sonarr.Client, episode *sonarrAPI.EpisodeResource, width, height int) common.SubModel {
	m := Model{
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
		switch msg.String() {
		case "esc":
			m.IsBack = true
			return m, nil
		}
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

func (m *Model) SetSize(width, height int) {
	height = height - boxStyle.GetVerticalFrameSize()

	m.Width = width
	m.Height = height

	m.table.SetWidth(width - boxStyle.GetHorizontalFrameSize())
	m.resizeTable(width - boxStyle.GetHorizontalFrameSize())
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
