package sonarr

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/subrr/internal/logging"
	"github.com/jon4hz/subrr/pkg/sonarr"
)

func (c *Client) AutomaticSearchSeason(seasonNumber int32) tea.Cmd {
	return func() tea.Msg {
		if c.serie == nil {
			return nil
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
