package season

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/jon4hz/subrr/internal/tui/common"
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
			BorderForeground(lipgloss.Color("#00CCFF"))
)

func (d Delegate) Height() int { return 5 }

func (d Delegate) Spacing() int { return 0 }

func (d Delegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
func (d Delegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var season string

	x, _ := defaultStyle.GetFrameSize()
	itemWidth := m.Width() - x
	width := itemWidth + defaultStyle.GetHorizontalPadding()

	i, ok := item.(EpisodeItem)
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
	subtileForeground  = lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			MaxHeight(1)

	downloadedMonitoredStyle = lipgloss.NewStyle().
					Foreground(lipgloss.AdaptiveColor{Light: "#4ECCA3", Dark: "#4ECCA3"})

	unmetCutoffMonitoredStyle = lipgloss.NewStyle().
					Foreground(lipgloss.AdaptiveColor{Light: "#FF9000", Dark: "#FF9000"})

	missingMonitoredStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#F71735", Dark: "#F71735"})

	downloadedUnmonitoredStyle = lipgloss.NewStyle().
					Foreground(lipgloss.AdaptiveColor{Light: "#36665E", Dark: "#36665E"})

	unmetCutoffUnmonitoredStyle = lipgloss.NewStyle().
					Foreground(lipgloss.AdaptiveColor{Light: "#945C1A", Dark: "#945C1A"})

	missingUnmonitoredStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#6B2334", Dark: "#6B2334"})
)

func renderItem(item EpisodeItem, itemWidth int) string {
	episodeTitle := fmt.Sprintf("%d. %s", item.episode.EpisodeNumber, item.episode.Title)

	title := zone.Mark(episodeTitle,
		truncate.StringWithTail(episodeTitle, uint(itemWidth), common.Ellipsis),
	)

	airDate := "---"
	if !item.episode.AirDateUTC.IsZero() {
		airDate = item.episode.AirDateUTC.Local().Format("02.01.2006")
	}
	airDate = truncate.StringWithTail(airDate, uint(itemWidth), common.Ellipsis)

	episodeStats := "Missing"
	if item.episode.EpisodeFile != nil {
		quality := "unknown"
		if item.episode.EpisodeFile.Quality != nil {
			quality = item.episode.EpisodeFile.Quality.Quality.Name
		}
		episodeStats = fmt.Sprintf("%s • %s", quality, humanize.IBytes(uint64(item.episode.EpisodeFile.Size)))
	}
	episodeStats = truncate.StringWithTail(episodeStats, uint(itemWidth), common.Ellipsis)

	if item.episode.Monitored {
		return renderMonitored(item, title, airDate, episodeStats)
	}
	return renderUnmonitored(item, title, airDate, episodeStats)
}

func renderMonitored(item EpisodeItem, title, airDate, episodeStats string) string {
	textColor := selectedForeground
	title = titleStyle.Foreground(textColor).Render(title)

	style := missingMonitoredStyle
	if item.episode.HasFile {
		if item.episode.EpisodeFile.QualityCutoffNotMet {
			style = unmetCutoffMonitoredStyle
		} else {
			style = downloadedMonitoredStyle
		}
	}
	episodeStats = style.Render(episodeStats)

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		airDate,
		episodeStats,
	)
}

func renderUnmonitored(item EpisodeItem, title, airDate, episodeStats string) string {
	textColor := subtileForeground
	title = titleStyle.Foreground(textColor).Render(title)
	airDate = lipgloss.NewStyle().Foreground(textColor).Render(airDate)

	style := missingUnmonitoredStyle
	if item.episode.HasFile {
		if item.episode.EpisodeFile.QualityCutoffNotMet {
			style = unmetCutoffUnmonitoredStyle
		} else {
			style = downloadedUnmonitoredStyle
		}
	}
	episodeStats = style.Render(episodeStats)

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		airDate,
		episodeStats,
	)
}