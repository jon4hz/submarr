package sonarr

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jon4hz/submarr/internal/config"
	"github.com/jon4hz/submarr/internal/httpclient"
)

// Client represents a sonarr client
type Client struct {
	http httpclient.Client
	cfg  *config.SonarrConfig
}

// New creates a new sonarr client
func New(httpClient httpclient.Client, cfg *config.SonarrConfig) *Client {
	return &Client{
		http: httpClient,
		cfg:  cfg,
	}
}

// Ping pings the sonarr server
func (c *Client) Ping(ctx context.Context) (*Ping, error) {
	var res Ping
	_, err := c.http.Get(ctx, c.cfg.Host, "/ping", &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetSeries returns a list of all series
func (c *Client) GetSeries(ctx context.Context) ([]*SeriesResource, error) {
	var res []*SeriesResource
	_, err := c.http.Get(ctx, c.cfg.Host, "/api/v3/series", &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

// GetSerie returns a serie by its TVDB ID
func (c *Client) GetSerie(ctx context.Context, tvdbID int32) (*SeriesResource, error) {
	var res []SeriesResource
	_, err := c.http.Get(ctx, c.cfg.Host, "/api/v3/series", &res,
		httpclient.WithParams(map[string]string{"tvdbId": strconv.FormatInt(int64(tvdbID), 10)}),
	)
	if err != nil {
		return nil, err
	}
	return &res[0], nil
}

// PutSerie updates a serie by its ID
func (c *Client) PutSerie(ctx context.Context, serie *SeriesResource, opts ...httpclient.RequestOpts) (*SeriesResource, error) {
	var res SeriesResource
	_, err := c.http.Put(ctx, c.cfg.Host, fmt.Sprintf("/api/v3/series/%d", serie.ID), &res, serie, opts...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// PostSerie adds a new serie
func (c *Client) PostSerie(ctx context.Context, serie *SeriesResource) (*SeriesResource, error) {
	var res SeriesResource
	_, err := c.http.Post(ctx, c.cfg.Host, "/api/v3/series", &res, serie)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) DeleteSerie(ctx context.Context, serieID int32, opts ...httpclient.RequestOpts) error {
	_, err := c.http.Delete(ctx, c.cfg.Host, fmt.Sprintf("/api/v3/series/%d", serieID), nil, nil, opts...)
	if err != nil {
		return err
	}
	return nil
}

// GetQueue returns the current download queue
func (c *Client) GetQueue(ctx context.Context, opts ...httpclient.RequestOpts) (*QueueResourcePagingResource, error) {
	var res QueueResourcePagingResource
	_, err := c.http.Get(ctx, c.cfg.Host, "/api/v3/queue", &res, opts...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetEpisodes returns a list of episodes for a given series and season
func (c *Client) GetEpisodes(ctx context.Context, seriesID, seasonNumber int32) ([]*EpisodeResource, error) {
	params := map[string]string{
		"seriesId":     strconv.Itoa(int(seriesID)),
		"seasonNumber": strconv.Itoa(int(seasonNumber)),
	}
	var res []*EpisodeResource
	_, err := c.http.Get(ctx, c.cfg.Host, "/api/v3/episode", &res,
		httpclient.WithParams(params),
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetEpisode returns an episode by its ID
func (c *Client) GetEpisode(ctx context.Context, episodeID int32) (*EpisodeResource, error) {
	var res EpisodeResource
	_, err := c.http.Get(ctx, c.cfg.Host, fmt.Sprintf("/api/v3/episode/%d", episodeID), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetAllEpisodes returns a list of all episodes for a given series
func (c *Client) GetAllEpisodes(ctx context.Context, seriesID int32) ([]*EpisodeResource, error) {
	params := map[string]string{
		"seriesId": strconv.Itoa(int(seriesID)),
	}
	var res []*EpisodeResource
	_, err := c.http.Get(ctx, c.cfg.Host, "/api/v3/episode", &res,
		httpclient.WithParams(params),
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetQualityProfiles returns a list of all quality profiles
func (c *Client) GetQualityProfiles(ctx context.Context) ([]*QualityProfileResource, error) {
	var res []*QualityProfileResource
	_, err := c.http.Get(ctx, c.cfg.Host, "/api/v3/qualityprofile", &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetQualityProfile returns a quality profile by its ID
func (c *Client) GetQualityProfile(ctx context.Context, id int) (*QualityProfileResource, error) {
	var res QualityProfileResource
	_, err := c.http.Get(ctx, c.cfg.Host, fmt.Sprintf("/api/v3/qualityprofile/%d", id), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// SetEpisodesMonitored sets the monitored status of a list of episodes
func (c *Client) SetEpisodesMonitored(ctx context.Context, params *EpisodesMonitoredResource) error {
	_, err := c.http.Put(ctx, c.cfg.Host, "/api/v3/episode/monitor", nil, params)
	return err
}

// PostCommand sends a command to sonarr
func (c *Client) PostCommand(ctx context.Context, params *CommandRequest) (*CommandResource, error) {
	var res CommandResource
	_, err := c.http.Post(ctx, c.cfg.Host, "/api/v3/command", &res, params)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetQueueDetails returns the queue for a certain series
func (c *Client) GetQueueDetails(ctx context.Context, seriesID int32) ([]*QueueResource, error) {
	var res []*QueueResource
	_, err := c.http.Get(ctx, c.cfg.Host, "/api/v3/queue/details", &res, httpclient.WithParams(map[string]string{"seriesId": fmt.Sprint(seriesID)}))
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetMissings returns all the missing episodes
func (c *Client) GetMissings(ctx context.Context) (*EpisodeResourcePagingResource, error) {
	var res EpisodeResourcePagingResource
	_, err := c.http.Get(ctx, c.cfg.Host, "/api/v3/wanted/missing", &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetSeriesLookup returns a list of series matching the given query
func (c *Client) GetSeriesLookup(ctx context.Context, query string) ([]*SeriesResource, error) {
	var res []*SeriesResource
	_, err := c.http.Get(ctx, c.cfg.Host, "/api/v3/series/lookup", &res, httpclient.WithParams(map[string]string{"term": query}))
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetRootFolders returns all root folders
func (c *Client) GetRootFolders(ctx context.Context) ([]*RootFolderResource, error) {
	var res []*RootFolderResource
	_, err := c.http.Get(ctx, c.cfg.Host, "/api/v3/rootfolder", &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetLanguageProfiles returns all language profiles
//
// Deprecated: Will be obsolete in Sonarr v4
func (c *Client) GetLanguageProfiles(ctx context.Context) ([]*LanguageProfileResource, error) {
	var res []*LanguageProfileResource
	_, err := c.http.Get(ctx, c.cfg.Host, "/api/v3/languageprofile", &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetHistory returns the history of an object
func (c *Client) GetHistory(ctx context.Context, opts ...httpclient.RequestOpts) (*HistoryResourcePagingResource, error) {
	var res *HistoryResourcePagingResource
	_, err := c.http.Get(ctx, c.cfg.Host, "/api/v3/history", &res, opts...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetEpisodeFiles returns all episode files for a given series
func (c *Client) GetEpisodeFiles(ctx context.Context, seriesID int32) ([]*EpisodeFileResource, error) {
	var res []*EpisodeFileResource
	_, err := c.http.Get(ctx, c.cfg.Host, "/api/v3/episodefile", &res, httpclient.WithParams(map[string]string{"seriesId": fmt.Sprint(seriesID)}))
	if err != nil {
		return nil, err
	}
	return res, nil
}
