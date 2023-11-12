package series

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/jon4hz/submarr/internal/core/sonarr"
	"github.com/jon4hz/submarr/internal/tui/common"
	"github.com/jon4hz/submarr/internal/tui/styles"
	sonarrAPI "github.com/jon4hz/submarr/pkg/sonarr"
	zone "github.com/lrstanley/bubblezone"
	"github.com/muesli/reflow/truncate"
)

type Delegate struct{}

var (
	DefaultStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true).
			Padding(0, 2).
			Margin(0, 1)

	SelectedStyle = DefaultStyle.Copy().
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(styles.SonarrBlue)
)

func (d Delegate) Height() int { return 6 }

func (d Delegate) Spacing() int { return 0 }

func (d Delegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d Delegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var serie string

	x, _ := DefaultStyle.GetFrameSize()
	itemWidth := m.Width() - x
	width := itemWidth + DefaultStyle.GetHorizontalPadding()

	i, ok := item.(sonarr.SeriesItem)
	if ok {
		serie = renderItem(i, itemWidth, index == m.Index())
	} else {
		return
	}

	if itemWidth-2 <= 0 {
		// short-circuit
		return
	}

	if index == m.Index() {
		serie = SelectedStyle.Width(width).Render(serie)
	} else {
		serie = DefaultStyle.Width(width).Render(serie)
	}

	fmt.Fprintf(w, "%s", serie)
}

var (
	SelectedForeground = lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}

	TitleStyle = lipgloss.NewStyle().
			Underline(true).
			Bold(true).
			MaxHeight(1)

	Separator = lipgloss.NewStyle().
			Foreground(styles.SubtleColor).
			Padding(0, 1).
			Render("•")
)

func renderItem(item sonarr.SeriesItem, itemWidth int, isSelected bool) string {
	textColor := SelectedForeground
	if !isSelected {
		textColor = styles.SubtleColor
	}

	title := TitleStyle.Foreground(textColor).Render(item.Series.Title)
	title = zone.Mark(item.Series.Title,
		truncate.StringWithTail(title, uint(itemWidth), common.Ellipsis),
	)

	var episodeStats string
	if item.Series.Statistics != nil {
		episodeCount := fmt.Sprintf("%d/%d (%.0f%%)", item.Series.Statistics.EpisodeFileCount, item.Series.Statistics.EpisodeCount, item.Series.Statistics.PercentOfEpisodes)
		seasonCount := fmt.Sprintf("%d Seasons", item.Series.Statistics.SeasonCount)
		diskSize := humanize.IBytes(uint64(item.Series.Statistics.SizeOnDisk))

		episodeStats = lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Foreground(textColor).Render(episodeCount),
			Separator,
			lipgloss.NewStyle().Foreground(textColor).Render(seasonCount),
			Separator,
			lipgloss.NewStyle().Foreground(textColor).Render(diskSize),
		)
		episodeStats = truncate.StringWithTail(episodeStats, uint(itemWidth), common.Ellipsis)
	}

	seriesType := common.Title(string(item.Series.SeriesType))
	profile := item.Series.ProfileName
	profileStats := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Foreground(textColor).Render(seriesType),
		Separator,
		lipgloss.NewStyle().Foreground(textColor).Render(profile),
	)
	profileStats = truncate.StringWithTail(profileStats, uint(itemWidth), common.Ellipsis)

	network := item.Series.Network
	var seriesStatus string
	switch item.Series.Status {
	case sonarrAPI.Continuing:
		if !item.Series.NextAiring.IsZero() {
			seriesStatus = item.Series.NextAiring.Local().Format("02.01.2006 @ 15:04")
		} else {
			seriesStatus = "Continuing"
		}
	default:
		seriesStatus = common.Title(string(item.Series.Status))
	}
	networkStats := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Foreground(textColor).Render(network),
		Separator,
		lipgloss.NewStyle().Foreground(textColor).Render(seriesStatus),
	)
	networkStats = truncate.StringWithTail(networkStats, uint(itemWidth), common.Ellipsis)

	s := lipgloss.JoinVertical(lipgloss.Top,
		title,
		episodeStats,
		profileStats,
		networkStats,
	)

	return s
}
