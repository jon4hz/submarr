package sonarr

import (
	"context"
	"errors"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/submarr/internal/logging"
)

type SearchSeriesResult struct {
	Items []list.Item
	Error error
}

func (c *Client) SearchSeries(term string) (tea.Cmd, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	return func() tea.Msg {
		res, err := c.sonarr.GetSeriesLookup(ctx, term)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			logging.Log.Error().Err(err).Msg("Failed to search series")
			return SearchSeriesResult{Error: err}
		}

		// Sanitize series
		sanitizeSeriesResources(res)

		var items []list.Item
		for _, s := range res {
			items = append(items, SeriesItem{s})
		}
		return SearchSeriesResult{Items: items}
	}, cancel
}
