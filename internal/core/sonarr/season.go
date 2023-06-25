package sonarr

import (
	"context"
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/subrr/internal/logging"
	"github.com/jon4hz/subrr/pkg/sonarr"
)

var ErrNoSerieSelected = errors.New("no serie selected")

type FetchSerieResult struct {
	Serie *sonarr.SeriesResource
	Error error
}

func (c *Client) ReloadSerie() tea.Cmd {
	return func() tea.Msg {
		// if currently no serie is selected, return error
		if c.serie == nil {
			return FetchSerieResult{Error: ErrNoSerieSelected}
		}

		s, err := c.sonarr.GetSerie(context.Background(), c.serie.TVDBID)
		// only update serie if there was no error
		if err != nil {
			logging.Log.Error().Err(err).Msg("Failed to reload serie")
			return FetchSerieResult{Serie: s, Error: fmt.Errorf("Failed to reload serie")}
		}
		c.serie = s
		return FetchSerieResult{Serie: s}
	}
}

func (c *Client) ToggleMonitorSeason(id int) tea.Cmd {
	return func() tea.Msg {
		if c.serie == nil {
			return FetchSerieResult{Error: ErrNoSerieSelected}
		}
		if len(c.serie.Seasons) <= id {
			return FetchSerieResult{Error: fmt.Errorf("Season not found")}
		}
		// toggle season monitored state
		c.serie.Seasons[id].Monitored = !c.serie.Seasons[id].Monitored
		serie, err := c.sonarr.PutSerie(context.Background(), c.serie)
		if err != nil {
			logging.Log.Error().Err(err).Msg("Failed to toggle season monitored state")
			return FetchSerieResult{Serie: c.serie, Error: fmt.Errorf("Failed to toggle season monitored state")}
		}
		c.serie = serie
		return FetchSerieResult{Serie: serie}
	}
}

func (c *Client) getSeasonEpisodes(season int32) error {
	if c.serie == nil {
		return ErrNoSerieSelected
	}

	var err error
	c.episodes, err = c.sonarr.GetEpisodes(context.Background(), c.serie.ID, season)
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to get episodes")
		return err
	}
	return nil
}
