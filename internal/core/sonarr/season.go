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

func (c *Client) ToggleMonitorSeries() tea.Cmd {
	return func() tea.Msg {
		if c.serie == nil {
			return FetchSerieResult{Error: ErrNoSerieSelected}
		}

		c.serie.Monitored = !c.serie.Monitored
		serie, err := c.sonarr.PutSerie(context.Background(), c.serie)
		if err != nil {
			logging.Log.Error().Err(err).Msg("Failed to toggle serie monitored state")
			return FetchSerieResult{Serie: c.serie, Error: fmt.Errorf("Failed to toggle serie monitored state")}
		}
		c.serie = serie
		return FetchSerieResult{Serie: serie}
	}
}

type FetchSeasonEpisodesResult struct {
	Episodes []*sonarr.EpisodeResource
	Error    error
}

func (c *Client) FetchSeasonEpisodes(season int32) tea.Cmd {
	return func() tea.Msg {
		if err := c.getSeasonEpisodes(season); err != nil {
			return FetchSeasonEpisodesResult{Error: err}
		}

		// refresh queue for current serie
		if err := c.getSeriesQueue(); err != nil {
			return FetchSeasonEpisodesResult{Error: err}
		}

		return FetchSeasonEpisodesResult{Episodes: c.seasonEpisodes}
	}
}

func (c *Client) getSeasonEpisodes(season int32) error {
	if c.serie == nil {
		return ErrNoSerieSelected
	}

	// Fetch all episodes of the selected season, without details
	var err error
	c.seasonEpisodes, err = c.sonarr.GetEpisodes(context.Background(), c.serie.ID, season)
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to get episodes")
		return err
	}

	// Fetch all episodes of the selected season, with details
	for i, episode := range c.seasonEpisodes {
		c.seasonEpisodes[i], err = c.sonarr.GetEpisode(context.Background(), episode.ID)
		if err != nil {
			logging.Log.Error().Err(err).Msg("Failed to get episode")
			continue
		}
	}

	return nil
}

func (c *Client) getSeriesQueue() error {
	if c.serie == nil {
		return ErrNoSerieSelected
	}
	queue, err := c.sonarr.GetQueueDetails(context.Background(), c.serie.ID)
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to get queue")
		return err
	}
	c.seriesQueue = queue
	return nil
}
