package sonarr

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/subrr/internal/tui/common"
	"github.com/muesli/reflow/truncate"
)

type ClientItem struct {
	c *Client
}

var queueStyle = lipgloss.NewStyle().Padding(0, 0, 0, 1).Align(lipgloss.Right)

func (i ClientItem) String() string { return "sonarr" }

func (i ClientItem) FilterValue() string { return "" }

func (i ClientItem) Render(itemWidth int) string {
	title := strings.Title(i.String())
	width := itemWidth - lipgloss.Width(title)
	if width < 2 {
		return truncate.StringWithTail(title, uint(itemWidth), common.Ellipsis)
	}
	queue := fmt.Sprintf("%d queued", i.c.queued)
	queue = truncate.StringWithTail(queue, uint(width-queueStyle.GetHorizontalPadding()), common.Ellipsis)
	queue = queueStyle.Width(itemWidth - lipgloss.Width(title)).Render(queue)
	return lipgloss.JoinHorizontal(lipgloss.Top, title, queue)
}
