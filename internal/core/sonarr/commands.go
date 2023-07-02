package sonarr

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/subrr/internal/logging"
	"github.com/jon4hz/subrr/pkg/sonarr"
)

func (c *Client) AutomaticSearchSeries() tea.Cmd {
	return func() tea.Msg {
		if c.serie == nil {
			return ErrNoSerieSelected
		}

		req := sonarr.CommandRequest{
			Name:     "SeriesSearch",
			SeriesID: c.serie.ID,
		}
		_, err := c.sonarr.PostCommand(context.Background(), &req)
		if err != nil {
			logging.Log.Err(err).Msg("failed to send command")
			return err
		}
		return nil
	}
}

func (c *Client) AutomaticSearchSeason(seasonNumber int32) tea.Cmd {
	return func() tea.Msg {
		if c.serie == nil {
			return ErrNoSerieSelected
		}

		req := sonarr.CommandRequest{
			Name:         "SeasonSearch",
			SeasonNumber: seasonNumber,
			SeriesID:     c.serie.ID,
		}
		_, err := c.sonarr.PostCommand(context.Background(), &req)
		if err != nil {
			logging.Log.Err(err).Msg("failed to send command")
			return err
		}
		return nil
	}
}

func (c *Client) RefreshSeries() tea.Cmd {
	return func() tea.Msg {
		if c.serie == nil {
			return ErrNoSerieSelected
		}

		req := sonarr.CommandRequest{
			Name:     "RefreshSeries",
			SeriesID: c.serie.ID,
		}
		_, err := c.sonarr.PostCommand(context.Background(), &req)
		if err != nil {
			logging.Log.Err(err).Msg("failed to send command")
			return err
		}
		return nil
	}
}
