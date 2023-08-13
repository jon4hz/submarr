package sonarr

import (
	"context"
	"strings"

	"github.com/jon4hz/subrr/internal/logging"
	"github.com/jon4hz/subrr/pkg/sonarr"
)

type Client struct {
	sonarr *sonarr.Client
	// is the client available?
	available bool
	// number of totalMissing episodes
	totalMissing int32
	// number of items in the download queue
	totalQueued int32
	// quality profiles by id
	qualityProfiles []*sonarr.QualityProfileResource
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
	// all available root folders
	rootFolders []*sonarr.RootFolderResource
	// all available languageProfiles
	languageProfiles []*sonarr.LanguageProfileResource
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
// TODO: gather only stats here, move the rest to the fetchers to another function which is called when the actual client is selected
func (c *Client) Init() error {
	ping, err := c.sonarr.Ping(context.Background())
	if err != nil {
		logging.Log.Error().Err(err).Msg("failed to ping sonarr")
		return err
	}
	if strings.ToLower(ping.Status) == "ok" {
		c.available = true
	} else {
		return nil
	}

	collectors := []func() error{
		c.FetchQueueNumber,
		c.FetchMissingNumber,
		c.FetchQualityProfiles,
		c.FetchRootFolders,
		c.FetchLanguageProfiles,
	}
	for _, collector := range collectors {
		if err := collector(); err != nil {
			return err
		}
	}

	return nil
}

// FetchQueueNumber fetches the number of items in the download queue
func (c *Client) FetchQueueNumber() error {
	queue, err := c.sonarr.GetQueue(context.Background())
	if err != nil {
		c.available = false
		logging.Log.Error().Err(err).Msg("failed to get queue")
		return err
	}

	if queue != nil {
		c.totalQueued = queue.TotalRecords
	}

	return nil
}

// FetchMissingNumber fetches the number of missing episodes
func (c *Client) FetchMissingNumber() error {
	totalMissings, err := c.sonarr.GetMissings(context.Background())
	if err != nil {
		c.available = false
		logging.Log.Error().Err(err).Msg("Failed to get totalMissing episodes")
		return err
	}

	if totalMissings != nil {
		c.totalMissing = totalMissings.TotalRecords
	}

	return nil
}

// FetchQualityProfiles fetches all quality profiles
func (c *Client) FetchQualityProfiles() error {
	profiles, err := c.sonarr.GetQualityProfiles(context.Background())
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to fetch quality profiles")
		return err
	}

	if len(profiles) == 0 {
		logging.Log.Warn().Msg("No quality profiles found")
		return nil
	}

	c.qualityProfiles = profiles

	return nil
}

// FetchRootFolders fetches all root folders
func (c *Client) FetchRootFolders() error {
	folders, err := c.sonarr.GetRootFolders(context.Background())
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to fetch root folders")
		return err
	}

	if len(folders) == 0 {
		logging.Log.Warn().Msg("No root folders found")
		return nil
	}

	c.rootFolders = folders

	return nil
}

// FetchLanguageProfiles fetches all language profiles
func (c *Client) FetchLanguageProfiles() error {
	profiles, err := c.sonarr.GetLanguageProfiles(context.Background())
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to fetch language profiles")
		return err
	}

	if len(profiles) == 0 {
		logging.Log.Warn().Msg("No language profiles found")
		return nil
	}

	c.languageProfiles = profiles

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

func (c *Client) GetQualityProfiles() []*sonarr.QualityProfileResource {
	return c.qualityProfiles
}

func (c *Client) GetRootFolders() []*sonarr.RootFolderResource {
	return c.rootFolders
}

func (c *Client) GetLanguageProfiles() []*sonarr.LanguageProfileResource {
	return c.languageProfiles
}

// GetQualityProfileByID returns a quality profile by id or an empty quality profile if not found
func (c *Client) GetLanguageProfileByID(id int32) *sonarr.LanguageProfileResource {
	if c.languageProfiles == nil {
		return nil
	}
	for _, profile := range c.languageProfiles {
		if profile.ID == id {
			return profile
		}
	}
	return new(sonarr.LanguageProfileResource)
}

func (c *Client) GetSerieQualityProfile() *sonarr.QualityProfileResource {
	return c.GetQualityProfileByID(c.serie.QualityProfileID)
}

// GetQualityProfileByID returns a quality profile by id or an empty quality profile if not found
func (c *Client) GetQualityProfileByID(id int32) *sonarr.QualityProfileResource {
	if c.qualityProfiles == nil {
		return nil
	}
	for _, profile := range c.qualityProfiles {
		if profile.ID == id {
			return profile
		}
	}
	return new(sonarr.QualityProfileResource)
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
