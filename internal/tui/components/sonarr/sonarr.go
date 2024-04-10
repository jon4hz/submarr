package sonarr

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/submarr/internal/core/sonarr"
	"github.com/jon4hz/submarr/internal/tui/common"
	"github.com/jon4hz/submarr/internal/tui/components/sonarr/overview"
	"github.com/jon4hz/submarr/internal/tui/components/statusbar"
	"github.com/jon4hz/submarr/internal/tui/components/tabs"
	"github.com/jon4hz/submarr/internal/tui/styles"
)

type Model struct {
	common.EmbedableModel

	client *sonarr.Client

	submodels []common.TabModel
	tabs      *tabs.Tabs
	activeTab int
}

func New(c *sonarr.Client, width, height int) *Model {
	submodels := []common.TabModel{
		overview.New(c, width, height-tabStyle.GetVerticalFrameSize()),
	}
	tabTitles := make([]string, len(submodels))
	for i, submodel := range submodels {
		tabTitles[i] = submodel.Title()
	}

	m := Model{
		client:    c,
		tabs:      tabs.New(tabTitles, styles.SonarrBlue, width, 1),
		submodels: submodels,
	}

	m.Width = width
	m.Height = height

	return &m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		statusbar.NewTitleCmd("Sonarr", statusbar.WithTitleForeground(styles.SonarrBlue)),
		statusbar.NewHelpCmd(DefaultKeyMap.FullHelp()),
		m.tabs.Init(),
		m.submodels[m.activeTab].Init(),
	)
}

func (m *Model) Update(msg tea.Msg) (common.SubModel, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			var cmd tea.Cmd
			m.tabs, cmd = m.tabs.Update(msg)
			return m, cmd
		}

	case tea.MouseEvent:
		if msg.Type == tea.MouseLeft {
			var cmd tea.Cmd
			m.tabs, cmd = m.tabs.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tabs.ActiveTabMsg:
		m.activeTab = int(msg)
		return m, m.submodels[m.activeTab].Init()
	}

	var cmd tea.Cmd
	m.submodels[m.activeTab], cmd = m.submodels[m.activeTab].Update(msg)
	cmds = append(cmds, cmd)

	if m.submodels[m.activeTab].Quit() {
		m.IsQuit = true
		return m, nil
	}
	if m.submodels[m.activeTab].Back() {
		m.IsBack = true
		return m, nil
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) SetSize(width, height int) {
	height = height - tabStyle.GetVerticalFrameSize() - 2
	m.Width = width

	m.tabs.SetSize(width, 1)

	m.Height = height

	for _, submodel := range m.submodels {
		submodel.SetSize(m.Width, m.Height)
	}
}

var tabStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(styles.SubtleColor).
	Padding(0, 2, 0, 2)

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		tabStyle.Width(m.Width).Render(m.tabs.View()),
		m.submodels[m.activeTab].View(),
	)
}
