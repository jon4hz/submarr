package search

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/submarr/internal/core/sonarr"
	"github.com/jon4hz/submarr/internal/tui/common"
	"github.com/jon4hz/submarr/internal/tui/components/sonarr/series"
	"github.com/jon4hz/submarr/internal/tui/styles"
	zone "github.com/lrstanley/bubblezone"
	"github.com/muesli/reflow/truncate"
)

type Delegate struct{}

var (
	defaultStyle = series.DefaultStyle.Copy()

	selectedStyle = series.SelectedStyle.Copy()

	statusStyle = lipgloss.NewStyle().
			Padding(0, 0, 0, 1).
			Align(lipgloss.Right)
)

func (d Delegate) Height() int { return 6 }

func (d Delegate) Spacing() int { return 0 }

func (d Delegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

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
	SelectedForeground = series.SelectedForeground

	TitleStyle = series.TitleStyle.Copy()

	Separator = series.Separator
)

func renderItem(item sonarr.SeriesItem, itemWidth int, isSelected bool) string {
	textColor := SelectedForeground
	if !isSelected {
		textColor = styles.SubtleColor
	}

	status := ""
	if !item.Series.Added.IsZero() {
		status = common.Available
	}
	status = statusStyle.Render(status)
	width := itemWidth - lipgloss.Width(status)

	title := TitleStyle.Foreground(textColor).Render(item.Series.Title)
	title = zone.Mark(item.Series.Title,
		truncate.StringWithTail(title, uint(width), common.Ellipsis),
	)

	title = lipgloss.JoinHorizontal(lipgloss.Left,
		title, lipgloss.PlaceHorizontal(itemWidth-lipgloss.Width(title), lipgloss.Right, status),
	)

	var episodeStats string
	if item.Series.Statistics != nil {
		seasonCount := fmt.Sprintf("%d Seasons", item.Series.Statistics.SeasonCount)

		episodeStats = lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Foreground(textColor).Render(seasonCount),
			Separator,
			lipgloss.NewStyle().Foreground(textColor).Render(fmt.Sprint(item.Series.Year)),
			Separator,
			lipgloss.NewStyle().Foreground(textColor).Render(item.Series.Network),
		)
		episodeStats = truncate.StringWithTail(episodeStats, uint(itemWidth), common.Ellipsis)
	}

	desc := truncate.StringWithTail(item.Series.Overview, uint(itemWidth)*2, common.Ellipsis)
	desc = lipgloss.NewStyle().Foreground(textColor).Width(itemWidth).Height(2).MaxHeight(2).Render(desc)

	s := lipgloss.JoinVertical(lipgloss.Top,
		title,
		episodeStats,
		desc,
	)

	return s
}
