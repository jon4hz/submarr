package sonarr

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/subrr/internal/logging"
	"github.com/jon4hz/subrr/pkg/sonarr"
)

type Client struct {
	sonarr *sonarr.Client

	// is the client available?
	available bool

	// some client stats
	missing int
	queued  int32
}

func New(sonarr *sonarr.Client) *Client {
	if sonarr == nil {
		return nil
	}
	return &Client{
		sonarr: sonarr,
	}
}

// Init initializes the client and fetches some stats
func (c *Client) Init() error {
	ping, err := c.sonarr.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("failed to ping sonarr: %w", err)
	}
	if strings.ToLower(ping.Status) == "ok" {
		c.available = true
	} else {
		return nil
	}

	queue, err := c.sonarr.GetQueue(context.Background())
	if err != nil {
		c.available = false
		return fmt.Errorf("failed to get queue: %w", err)
	}
	c.queued = queue.TotalRecords

	return nil
}

func (c *Client) ClientListItem() ClientItem {
	return ClientItem{c}
}

type FetchSeriesResult struct {
	Items []list.Item
	Error error
}

type SeriesItem struct {
	Series sonarr.SeriesResource
}

func (s SeriesItem) FilterValue() string {
	return s.Series.Title
}

func (c *Client) FetchSeries() tea.Cmd {
	return func() tea.Msg {
		series, err := c.sonarr.GetSeries(context.Background())
		if err != nil {
			logging.Log.Error().Err(err).Msg("Failed to fetch series")
			return FetchSeriesResult{Error: err}
		}

		// sort series by title
		sort.Slice(series, func(i, j int) bool {
			return series[i].SortTitle < series[j].SortTitle
		})

		var items []list.Item
		for _, s := range series {
			items = append(items, SeriesItem{s})
		}
		return FetchSeriesResult{Items: items}
	}
}
