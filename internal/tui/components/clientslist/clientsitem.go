package clientslist

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/submarr/internal/tui/common"
	"github.com/jon4hz/submarr/internal/tui/styles"
	zone "github.com/lrstanley/bubblezone"
	"github.com/muesli/reflow/truncate"
)

type ClientsItem interface {
	fmt.Stringer
	// just to fulfill the list.Item interface
	FilterValue() string
	// Return the title of the client
	Title() string
	// Whether the client is available
	Available() bool
	// Some stats about the client. Will be displayed next to each other separated by a dot
	Stats() []string
}

var (
	selectedForeground = lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}

	titleStyle = lipgloss.NewStyle().
			Underline(true).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(styles.SubtileColor).
			Padding(0, 0, 0, 1).
			Align(lipgloss.Right)

	statsStyle = lipgloss.NewStyle()

	separator = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
			Padding(0, 1).
			Render("•")
)

func renderItem(item ClientsItem, itemWidth int, isSelected bool) string {
	textColor := selectedForeground
	if !isSelected {
		textColor = styles.SubtileColor
	}

	title := titleStyle.Foreground(textColor).Render(item.Title() + "\n")

	width := itemWidth - lipgloss.Width(title)
	if width < 2 {
		return zone.Mark(item.String(),
			truncate.StringWithTail(title, uint(itemWidth), common.Ellipsis),
		)
	}
	title = zone.Mark(item.String(), title)

	status := "available " + common.Available
	if !item.Available() {
		status = "unavailable " + common.Unavailable
	}
	status = truncate.StringWithTail(status, uint(width-statusStyle.GetHorizontalPadding()), common.Ellipsis)
	status = statusStyle.Width(itemWidth - lipgloss.Width(title)).Render(status)

	var statsBuilder strings.Builder
	for i, stat := range item.Stats() {
		statsBuilder.WriteString(statsStyle.Foreground(textColor).Render(stat))
		if i < len(item.Stats())-1 {
			statsBuilder.WriteString(separator)
		}
	}
	stats := truncate.StringWithTail(statsBuilder.String(), uint(itemWidth), common.Ellipsis)

	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top,
			title, status,
		),
		stats,
	)
}
