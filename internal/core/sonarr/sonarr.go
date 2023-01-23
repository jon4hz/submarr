package sonarr

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/subrr/pkg/sonarr"
	"github.com/muesli/reflow/truncate"
)

const ellipsis = "…"

type Client struct {
	sonarr *sonarr.Client

	// is the client available?
	available bool

	// some client stats
	missing int
	queued  int
}

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
		return truncate.StringWithTail(title, uint(itemWidth), ellipsis)
	}
	queue := fmt.Sprintf("%d queued", i.c.queued)
	queue = truncate.StringWithTail(queue, uint(width-queueStyle.GetHorizontalPadding()), ellipsis)
	queue = queueStyle.Width(itemWidth - lipgloss.Width(title)).Render(queue)
	return lipgloss.JoinHorizontal(lipgloss.Top, title, queue)
}

func (c *Client) ListItem() ClientItem {
	return ClientItem{c}
}

func New(sonarr *sonarr.Client) *Client {
	return &Client{
		sonarr: sonarr,
	}
}

func (c *Client) Init() error {
	ping, err := c.sonarr.Ping(context.Background())
	if err != nil {
		return err
	}
	if strings.ToLower(ping.Status) == "ok" {
		c.available = true
	} else {
		return nil
	}

	queue, err := c.sonarr.GetQueue(context.Background())
	if err != nil {
		return err
	}
	c.queued = int(queue.TotalRecords)

	return nil
}
