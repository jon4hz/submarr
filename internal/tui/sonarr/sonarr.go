package sonarr

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/subrr/internal/core/sonarr"
	"github.com/jon4hz/subrr/internal/tui/common"
)

type Model struct {
	common.EmbedableModel

	client *sonarr.Client
}

func New(c *sonarr.Client) *Model {
	return &Model{
		client: c,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (common.ClientModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.IsBack = true

		case "q":
			m.IsQuit = true
		}
	}

	return m, nil
}

func (m *Model) SetSize(width, height int) {
	m.Width = width
	m.Height = height
}

func (m Model) View() string {
	return "test"
}
