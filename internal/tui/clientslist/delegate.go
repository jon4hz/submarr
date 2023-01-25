package clientslist

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	ellipsis = "…"
)

type DefaultItemStyles struct {
	DefaultClient  lipgloss.Style
	SelectedSonarr lipgloss.Style
	SelectedRadarr lipgloss.Style
	SelectedLidarr lipgloss.Style
}

func NewDefaultItemStyles() (s DefaultItemStyles) {
	s.DefaultClient = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		Padding(0, 2).
		Margin(0, 1)

	s.SelectedSonarr = s.DefaultClient.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		BorderForeground(lipgloss.Color("#00CCFF"))

	s.SelectedRadarr = s.DefaultClient.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		BorderForeground(lipgloss.Color("#FFA500"))

	return s
}

type clientDelegate struct {
	Styles        DefaultItemStyles
	UpdateFunc    func(tea.Msg, *list.Model) tea.Cmd
	ShortHelpFunc func() []key.Binding
	FullHelpFunc  func() [][]key.Binding
}

func newClientDelegate() clientDelegate {
	return clientDelegate{
		Styles: NewDefaultItemStyles(),
	}
}

func (d clientDelegate) Height() int { return 5 }

func (d clientDelegate) Spacing() int { return 1 }

func (d clientDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	if d.UpdateFunc == nil {
		return nil
	}
	return d.UpdateFunc(msg, m)
}

func (d clientDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var (
		client string
		s      = &d.Styles
	)

	x, _ := s.DefaultClient.GetFrameSize()
	itemWidth := m.Width() - x
	width := itemWidth + s.DefaultClient.GetHorizontalPadding()

	i, ok := item.(ClientsItem)
	if ok {
		client = renderItem(i, itemWidth, index == m.Index())
	} else {
		return
	}

	if itemWidth-2 <= 0 {
		// short-circuit
		return
	}

	if index == m.Index() {
		switch strings.ToLower(i.String()) {
		case "sonarr":
			client = s.SelectedSonarr.Width(width).Render(client)

		case "radarr":
			client = s.SelectedRadarr.Width(width).Render(client)

		default:
			client = s.DefaultClient.Width(width).Render(client)
		}
	} else {
		client = s.DefaultClient.Width(width).Render(client)
	}

	fmt.Fprintf(w, "%s", client)
}

func (d clientDelegate) ShortHelp() []key.Binding {
	if d.ShortHelpFunc != nil {
		return d.ShortHelpFunc()
	}
	return nil
}

func (d clientDelegate) FullHelp() [][]key.Binding {
	if d.FullHelpFunc != nil {
		return d.FullHelpFunc()
	}
	return nil
}
