package clientslist

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/submarr/internal/tui/styles"
)

type DefaultItemStyles struct {
	DefaultClient  lipgloss.Style
	SelectedSonarr lipgloss.Style
	SelectedRadarr lipgloss.Style
}

var (
	defaultClient = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			Padding(0, 2).
			Margin(0, 1)

	selectedStyle = defaultClient.Copy().
			BorderStyle(lipgloss.ThickBorder())

	selectedSonarr = selectedStyle.Copy().
			Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
			BorderForeground(styles.SonarrBlue)

	selectedRadarr = selectedStyle.Copy().
			Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
			BorderForeground(lipgloss.Color("#FFA500"))
)

type clientDelegate struct {
	Styles DefaultItemStyles
}

func (d clientDelegate) Height() int { return 5 }

func (d clientDelegate) Spacing() int { return 1 }

func (d clientDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d clientDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var client string

	x, _ := defaultClient.GetFrameSize()
	itemWidth := m.Width() - x
	width := itemWidth + defaultClient.GetHorizontalPadding()

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
			client = selectedSonarr.Width(width).Render(client)

		case "radarr":
			client = selectedRadarr.Width(width).Render(client)

		default:
			client = defaultClient.Width(width).Render(client)
		}
	} else {
		client = defaultClient.Width(width).Render(client)
	}

	fmt.Fprintf(w, "%s", client)
}
