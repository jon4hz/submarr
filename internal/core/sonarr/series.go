package sonarr

import (
	"context"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/submarr/internal/logging"
	"github.com/jon4hz/submarr/pkg/sonarr"
)

type FetchSeriesResult struct {
	Items []list.Item
	Error error
}

type AddSeriesResult struct {
	AddedTitle string
	Items      []list.Item
	Error      error
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
			logging.Log.Error("Failed to fetch series", "err", err)
			return FetchSeriesResult{Error: err}
		}
		return FetchSeriesResult{Items: c.newSeriesItems()}
	}
}

func (c *Client) fetchSeries() error {
	series, err := c.sonarr.GetSeries(context.Background())
	if err != nil {
		logging.Log.Error("Failed to fetch series", "err", err)
		return err
	}

	// Add the quality profile name to the series
	for i := range series {
		qp := c.GetQualityProfileByID(series[i].QualityProfileID)
		if qp != nil {
			series[i].ProfileName = qp.Name
		}
	}

	// Sanitize series
	sanitizeSeriesResources(series)

	sortSeries(series)
	c.series = series

	return nil
}

// newSeriesItems creates a list item for each series
func (c *Client) newSeriesItems() []list.Item {
	var items []list.Item
	for _, s := range c.series {
		items = append(items, SeriesItem{s})
	}
	return items
}

// sortSeries sorts the series by their sort title
func sortSeries(series []*sonarr.SeriesResource) {
	sort.Slice(series, func(i, j int) bool {
		return series[i].SortTitle < series[j].SortTitle
	})
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

func (c *Client) PostSeries(series *sonarr.SeriesResource) tea.Cmd {
	return func() tea.Msg {
		resp, err := c.sonarr.PostSerie(context.Background(), series)
		if err != nil {
			logging.Log.Error("Failed to add series", "err", err)
			return AddSeriesResult{Error: err}
		}
		c.series = append(c.series, resp)
		sortSeries(c.series)
		return AddSeriesResult{
			AddedTitle: series.Title,
			Items:      c.newSeriesItems(),
		}
	}
}
