package seasons

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/jon4hz/submarr/internal/tui/common"
	"github.com/jon4hz/submarr/internal/tui/styles"
	zone "github.com/lrstanley/bubblezone"
	"github.com/muesli/reflow/truncate"
)

type Delegate struct{}

var (
	defaultStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true).
			Padding(0, 2).
			Margin(0, 1)

	selectedStyle = defaultStyle.Copy().
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(styles.SonarrBlue)
)

func (d Delegate) Height() int { return 6 }

func (d Delegate) Spacing() int { return 0 }

func (d Delegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d Delegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var season string

	x, _ := defaultStyle.GetFrameSize()
	itemWidth := m.Width() - x
	width := itemWidth + defaultStyle.GetHorizontalPadding()

	i, ok := item.(SeasonItem)
	if ok {
		season = renderItem(i, itemWidth)
	} else {
		return
	}

	if itemWidth-2 <= 0 {
		// short-circuit
		return
	}

	if index == m.Index() {
		season = selectedStyle.Width(width).Render(season)
	} else {
		season = defaultStyle.Width(width).Render(season)
	}

	fmt.Fprintf(w, "%s", season)
}

var (
	selectedForeground = lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			MaxHeight(1)

	completeSeasonMonitoredStyle = lipgloss.NewStyle().
					Foreground(lipgloss.AdaptiveColor{Light: "#4ECCA3", Dark: "#4ECCA3"})

	incompleteSeasonMonitoredStyle = lipgloss.NewStyle().
					Foreground(lipgloss.AdaptiveColor{Light: "#F71735", Dark: "#F71735"})

	completeSeasonUnmonitoredStyle = lipgloss.NewStyle().
					Foreground(lipgloss.AdaptiveColor{Light: "#36665E", Dark: "#36665E"})

	incompleteSeasonUnmonitoredStyle = lipgloss.NewStyle().
						Foreground(lipgloss.AdaptiveColor{Light: "#6B2334", Dark: "#6B2334"})
)

func renderItem(item SeasonItem, itemWidth int) string {
	symbol := common.Unselected
	if item.Season.Monitored {
		symbol = common.Selected
	}
	seasonTitle := fmt.Sprintf("%s Season %d", symbol, item.Season.SeasonNumber)

	title := zone.Mark(fmt.Sprintf("Season %d", item.Season.SeasonNumber),
		truncate.StringWithTail(seasonTitle, uint(itemWidth), common.Ellipsis),
	)

	prevAiring := "---"
	diskSize := "---"
	var seasonStats string
	if item.Season.Statistics != nil {
		if !item.Season.Statistics.PreviousAiring.IsZero() {
			prevAiring = item.Season.Statistics.PreviousAiring.Local().Format("02.01.2006 @ 15:04")
		}
		prevAiring = truncate.StringWithTail(prevAiring, uint(itemWidth), common.Ellipsis)

		diskSize = humanize.Bytes(uint64(item.Season.Statistics.SizeOnDisk))
		diskSize = truncate.StringWithTail(diskSize, uint(itemWidth), common.Ellipsis)

		seasonStats = fmt.Sprintf("%.0f%% • %d/%d Episodes Available", item.Season.Statistics.PercentOfEpisodes, item.Season.Statistics.EpisodeFileCount, item.Season.Statistics.EpisodeCount)
		seasonStats = truncate.StringWithTail(seasonStats, uint(itemWidth), common.Ellipsis)
	}

	if item.Season.Monitored {
		return renderMonitored(item, title, prevAiring, diskSize, seasonStats)
	}
	return renderUnmonitored(item, title, diskSize, seasonStats)
}

func renderMonitored(item SeasonItem, title, prevAiring, diskSize, seasonStats string) string {
	textColor := selectedForeground
	title = titleStyle.Foreground(textColor).Render(title)

	style := incompleteSeasonMonitoredStyle
	if item.Season.Statistics.PercentOfEpisodes == 100 {
		style = completeSeasonMonitoredStyle
	}
	seasonStats = style.Render(seasonStats)

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		prevAiring,
		diskSize,
		seasonStats,
	)
}

func renderUnmonitored(item SeasonItem, title, diskSize, seasonStats string) string {
	textColor := styles.SubtleColor
	title = titleStyle.Foreground(textColor).Render(title)
	prevAiring := lipgloss.NewStyle().Foreground(textColor).Render("---")
	diskSize = lipgloss.NewStyle().Foreground(textColor).Render(diskSize)

	style := incompleteSeasonUnmonitoredStyle
	if item.Season.Statistics.PercentOfEpisodes == 100 {
		style = completeSeasonUnmonitoredStyle
	}
	seasonStats = style.Render(seasonStats)

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		prevAiring,
		diskSize,
		seasonStats,
	)
}
