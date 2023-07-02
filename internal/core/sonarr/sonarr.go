package sonarr

import (
	"context"
	"fmt"
	"strings"

	"github.com/jon4hz/subrr/pkg/sonarr"
)

type Client struct {
	sonarr *sonarr.Client

	// is the client available?
	available bool

	// some client stats
	missing int
	queued  int32

	// quality profiles by id
	qualityProfiles map[int32]*sonarr.QualityProfileResource

	// all available series
	series []*sonarr.SeriesResource

	// currently selected serie
	serie *sonarr.SeriesResource

	// currently selected season
	season *sonarr.SeasonResource

	// episodes of the currently selected season
	seasonEpisodes []*sonarr.EpisodeResource

	// queue of the currently selected serie
	seriesQueue []*sonarr.QueueResource
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

func (c *Client) SetSerie(serie *sonarr.SeriesResource) {
	c.serie = serie
}

func (c *Client) GetSerie() *sonarr.SeriesResource {
	return c.serie
}

func (c *Client) GetSerieQualityProfile() *sonarr.QualityProfileResource {
	if c.qualityProfiles == nil {
		return nil
	}
	return c.qualityProfiles[c.serie.QualityProfileID]
}

func (c *Client) SetSeason(season *sonarr.SeasonResource) {
	c.season = season
}

func (c *Client) GetSeason() *sonarr.SeasonResource {
	return c.season
}

func (c *Client) GetSeasonEpisodes() []*sonarr.EpisodeResource {
	return c.seasonEpisodes
}

func (c *Client) GetSeriesQueue() []*sonarr.QueueResource {
	return c.seriesQueue
}
