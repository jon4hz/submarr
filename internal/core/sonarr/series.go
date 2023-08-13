package sonarr

import (
	"context"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/subrr/internal/logging"
	"github.com/jon4hz/subrr/pkg/sonarr"
)

type FetchSeriesResult struct {
	Items []list.Item
	Error error
}

type SeriesItem struct {
	Series *sonarr.SeriesResource
}

func (s SeriesItem) FilterValue() string {
	return s.Series.Title
}

func (c *Client) FetchSeries() tea.Cmd {
	return func() tea.Msg {
		if err := c.fetchSeries(); err != nil {
			return FetchSeriesResult{Error: err}
		}

		// Create a list item for each series
		var items []list.Item
		for _, s := range c.series {
			items = append(items, SeriesItem{s})
		}
		return FetchSeriesResult{Items: items}
	}
}

func (c *Client) fetchSeries() error {
	series, err := c.sonarr.GetSeries(context.Background())
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to fetch series")
		return err
	}

	// sort series by title
	sort.Slice(series, func(i, j int) bool {
		return series[i].SortTitle < series[j].SortTitle
	})

	// Add the quality profile name to the series
	for i := range series {
		qp := c.GetQualityProfileByID(series[i].QualityProfileID)
		if qp != nil {
			series[i].ProfileName = qp.Name
		}
	}

	// Sanitize series
	sanitizeSeriesResources(series)

	c.series = series
	return nil
}

func sanitizeSeriesResources(series []*sonarr.SeriesResource) {
	for i := range series {
		series[i].Title = sanitizeTitle(series[i].Title)
		series[i].Overview = sanitizeOverview(series[i].Overview)
	}
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

var punctuationRe = regexp.MustCompile(`([,.:])(\S)`)

// sanitizeOverview removes all newline and tab characters from the overview.
func sanitizeOverview(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.ReplaceAll(s, "\r", "")

	res := punctuationRe.ReplaceAllString(s, "$1 $2")
	return strings.TrimSpace(res)
}
