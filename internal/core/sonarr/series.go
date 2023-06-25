package sonarr

import (
	"context"
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

	// Fetch all quality profiles
	profiles, err := c.sonarr.GetQualityProfiles(context.Background())
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to fetch quality profiles")
		return err
	}

	// Group profiles by ID
	profilesByID := make(map[int32]*sonarr.QualityProfileResource)
	for _, p := range profiles {
		profilesByID[p.ID] = p
	}

	// store the profiles in the client so we can use them later
	c.qualityProfiles = profilesByID

	// Add the quality profile name to the series
	// And sanitize the title
	for i := range series {
		series[i].ProfileName = profilesByID[series[i].QualityProfileID].Name
		series[i].Title = sanitizeTitle(series[i].Title)
	}

	c.series = series
	return nil
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
