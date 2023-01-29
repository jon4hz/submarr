package series

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/jon4hz/subrr/internal/core/sonarr"
	"github.com/jon4hz/subrr/internal/tui/common"
	sonarrAPI "github.com/jon4hz/subrr/pkg/sonarr"
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

func (d Delegate) Height() int { return 6 }

func (d Delegate) Spacing() int { return 1 }

func (d Delegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d Delegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var serie string

	x, _ := defaultStyle.GetFrameSize()
	itemWidth := m.Width() - x
	width := itemWidth + defaultStyle.GetHorizontalPadding()

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
		serie = selectedStyle.Width(width).Render(serie)
	} else {
		serie = defaultStyle.Width(width).Render(serie)
	}

	fmt.Fprintf(w, "%s", serie)
}

var (
	selectedForeground = lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}
	subtileForeground  = lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}

	titleStyle = lipgloss.NewStyle().
			Underline(true).
			Bold(true).
			MaxHeight(1)

	separator = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
			Padding(0, 1).
			Render("•")
)

func renderItem(item sonarr.SeriesItem, itemWidth int, isSelected bool) string {
	textColor := selectedForeground
	if !isSelected {
		textColor = subtileForeground
	}

	title := titleStyle.Foreground(textColor).Render(sanitizeTitle(item.Series.Title))
	title = zone.Mark(item.Series.Title,
		truncate.StringWithTail(title, uint(itemWidth), common.Ellipsis),
	)

	var episodeStats string
	if item.Series.Statistics != nil {
		episodeCount := fmt.Sprintf("%d/%d (%.0f%%)", item.Series.Statistics.EpisodeFileCount, item.Series.Statistics.EpisodeCount, item.Series.Statistics.PercentOfEpisodes)
		seasonCount := fmt.Sprintf("%d Seasons", item.Series.Statistics.SeasonCount)
		diskSize := fmt.Sprintf("%s", humanize.Bytes(uint64(item.Series.Statistics.SizeOnDisk)))

		episodeStats = lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Foreground(textColor).Render(episodeCount),
			separator,
			lipgloss.NewStyle().Foreground(textColor).Render(seasonCount),
			separator,
			lipgloss.NewStyle().Foreground(textColor).Render(diskSize),
		)
		episodeStats = truncate.StringWithTail(episodeStats, uint(itemWidth), common.Ellipsis)
	}

	seriesType := common.Title(string(item.Series.SeriesType))
	profile := item.Series.ProfileName
	profileStats := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Foreground(textColor).Render(seriesType),
		separator,
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
		separator,
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

// sanitizeTitle replaces all unicode whitespace characters with a single space.
// For some weird reason, some titles contain characters like U+00A0 (NO-BREAK SPACE)
func sanitizeTitle(s string) string {
	for _, r := range s {
		if unicode.IsSpace(r) {
			s = strings.Replace(s, string(r), " ", -1)
		}
	}
	return s
}
