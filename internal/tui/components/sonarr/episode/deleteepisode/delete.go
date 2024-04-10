package deleteepisode

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	confirm bool
}

func New(client *sonarr.Client, episode *sonarrAPI.EpisodeResource, width, height int) common.SubModel {
	m := Model{
		client:  client,
		episode: episode,
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
		case "enter":
			if m.confirm {
				return m, m.client.DeleteEpisodeFile(m.episode)
			} else {
				m.IsBack = true
			}
		case "tab", "right", "left":
			m.confirm = !m.confirm
		}
	}
	return m, nil
}

func (m *Model) SetSize(width, height int) {
	m.Width = width
	m.Height = height
}

var (
	buttonStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Margin(0, 1)

	boxStyle = lipgloss.NewStyle()
)

func (m Model) View() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Are you sure you want to delete %q?\n\n", m.episode.EpisodeFile.Path))

	var confirmCol lipgloss.Color
	var cancelCol lipgloss.Color
	if m.confirm {
		confirmCol = styles.SonarrBlue
	} else {
		cancelCol = styles.SonarrBlue
	}
	s.WriteString(lipgloss.Place(
		m.Width/2, 1,
		lipgloss.Center, lipgloss.Top,
		buttonStyle.Background(cancelCol).Underline(!m.confirm).Render("Cancel")+
			buttonStyle.Background(confirmCol).Underline(m.confirm).Render("Confirm"),
	))
	return boxStyle.Width(m.Width / 2).Render(s.String())
}
