package sonarr

import (
	"context"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/submarr/internal/logging"
	"github.com/jon4hz/submarr/pkg/sonarr"
)

var ErrNoEpisodes = errors.New("no episodes provided")

func (c *Client) doCommandRequest(req *sonarr.CommandRequest) (*sonarr.CommandResource, error) { // nolint:unparam
	res, err := c.sonarr.PostCommand(context.Background(), req)
	if err != nil {
		logging.Log.Error("Failed to send command", "err", err)
		return nil, err
	}
	return res, nil
}

func (c *Client) AutomaticSearchEpisode(epiodeIDs ...int32) tea.Cmd {
	return func() tea.Msg {
		if len(epiodeIDs) == 0 {
			logging.Log.Error(ErrNoEpisodes)
			return ErrNoEpisodes
		}
		req := sonarr.CommandRequest{
			Name:       "EpisodeSearch",
			EpisodeIDs: epiodeIDs,
		}
		_, err := c.doCommandRequest(&req)
		return err
	}
}

func (c *Client) AutomaticSearchSeries() tea.Cmd {
	return func() tea.Msg {
		if c.serie == nil {
			logging.Log.Error(ErrNoSerieSelected)
			return ErrNoSerieSelected
		}

		req := sonarr.CommandRequest{
			Name:     "SeriesSearch",
			SeriesID: c.serie.ID,
		}
		_, err := c.doCommandRequest(&req)
		return err
	}
}

func (c *Client) AutomaticSearchSeason(seasonNumber int32) tea.Cmd {
	return func() tea.Msg {
		if c.serie == nil {
			logging.Log.Error(ErrNoSerieSelected)
			return ErrNoSerieSelected
		}

		req := sonarr.CommandRequest{
			Name:         "SeasonSearch",
			SeasonNumber: seasonNumber,
			SeriesID:     c.serie.ID,
		}
		_, err := c.doCommandRequest(&req)
		return err
	}
}

func (c *Client) RefreshSeries() tea.Cmd {
	return func() tea.Msg {
		if c.serie == nil {
			logging.Log.Error(ErrNoSerieSelected)
			return ErrNoSerieSelected
		}

		req := sonarr.CommandRequest{
			Name:     "RefreshSeries",
			SeriesID: c.serie.ID,
		}
		_, err := c.doCommandRequest(&req)
		return err
	}
}
