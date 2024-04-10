package sonarr

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/submarr/internal/core/sonarr"
	"github.com/jon4hz/submarr/internal/tui/common"
	"github.com/jon4hz/submarr/internal/tui/components/sonarr/overview"
	"github.com/jon4hz/submarr/internal/tui/components/statusbar"
	"github.com/jon4hz/submarr/internal/tui/styles"
)

type Model struct {
	common.EmbedableModel

	client *sonarr.Client

	submodel common.SubModel
}

func New(c *sonarr.Client, width, height int) *Model {
	m := Model{
		client:   c,
		submodel: overview.New(c, width, height),
	}

	m.Width = width
	m.Height = height

	return &m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		statusbar.NewTitleCmd("Sonarr", statusbar.WithTitleForeground(styles.SonarrBlue)),
		statusbar.NewHelpCmd(DefaultKeyMap.FullHelp()),
		m.submodel.Init(),
	)
}

func (m *Model) Update(msg tea.Msg) (common.SubModel, tea.Cmd) {
	var cmd tea.Cmd
	m.submodel, cmd = m.submodel.Update(msg)
	if m.submodel.Back() {
		m.IsBack = true
	}
	if m.submodel.Quit() {
		m.IsQuit = true
	}
	return m, cmd
}

func (m *Model) SetSize(width, height int) {
	m.submodel.SetSize(width, height)
}

func (m Model) View() string {
	return m.submodel.View()
}
